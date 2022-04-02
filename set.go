// set is a library designed to help you do set operations on comparable data types
package set

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// This error is returned when you try to remove an item from a set that doesn't exist
	ErrElementNotFound = errors.New("element not found")
)

type Set[T comparable] struct {
	data map[T]struct{}
}

// NewSet will return a Set object from an input slice, or anything that has a slice as
// the underlying data type
func NewSet[T comparable, S ~[]T](data S) Set[T] {
	// Create an empty map of the correct size
	result := make(map[T]struct{}, len(data))

	// Fill it up
	for _, v := range data {
		result[v] = struct{}{}
	}

	return Set[T]{data: result}
}

// NewSetWithCapacity will return a Set object with a specific capacity. Note that if
// len(data) >= size, size will simply be ignored. This function is most useful in cases
// where you know you will be adding more elements than go in initially, and you have an
// estimate of how many total will go in
func NewSetWithCapacity[T comparable, S ~[]T](data S, size int) Set[T] {
	// Choose whichever has the larger size
	s := size
	if len(data) > size {
		s = len(data)
	}

	// Create an empty map of the correct size
	result := make(map[T]struct{}, s)

	// Fill it up
	for _, v := range data {
		result[v] = struct{}{}
	}

	return Set[T]{data: result}
}

func (s Set[T]) String() string {
	var b strings.Builder
	last_index := s.Len() - 1
	index := -1
	b.WriteString("{")
	for v := range s.data {
		index += 1

		if index < last_index {
			b.WriteString(fmt.Sprintf("%v, ", v))
		} else {
			b.WriteString(fmt.Sprintf("%v}", v))
		}
	}

	return b.String()
}

// Slice will return all the items in the set as a slice. They are not guaranteed in any
// particular order.
func (s *Set[T]) Slice() []T {
	result := make([]T, 0, s.Len())

	for v := range s.data {
		result = append(result, v)
	}

	return result
}

// Contains will return true if the set contains the item. If the set is empty, returns
// false
func (s *Set[T]) Contains(item T) bool {
	_, ok := s.data[item]
	return ok
}

// Len returns the length of the Set
func (s *Set[T]) Len() int {
	return len(s.data)
}

// IsEmpty returns true if the set is empty
func (s *Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

// Add will add a new item to `s`. If it already exists, it is ignored
func (s *Set[T]) Add(item T) {
	s.data[item] = struct{}{}
}

// Remove removes an item from the set. Returns an error if the item doesn't exist
func (s *Set[T]) Remove(item T) error {
	if !s.Contains(item) {
		return ErrElementNotFound
	}

	delete(s.data, item)
	return nil
}

// Discard removes an item from the set. If it doesn't exist, it is ignored
func (s *Set[T]) Discard(item T) {
	delete(s.data, item)
}

// Pop will remove and return an arbitrary item from the set. If the set is empty,
// it will return an error
func (s *Set[T]) Pop() (item T, err error) {
	if s.IsEmpty() {
		return item, ErrElementNotFound
	}

	// Get the first item
	for item = range s.data {
		break
	}

	// Discard it
	s.Discard(item)

	return item, nil
}

// Clear will remove all items from the set
func (s *Set[T]) Clear() {
	s.data = make(map[T]struct{})
}

// Copy makes a deep copy as quickly as possible
func (s *Set[T]) Copy() Set[T] {
	// Make sure to allocate the same size
	copy := make(map[T]struct{}, len(s.data))

	// Fill it up
	for v := range s.data {
		copy[v] = struct{}{}
	}

	return Set[T]{data: copy}
}

// Equals will return true if `s` and `t` are
// - the same length
// - contain the same elements
func (s *Set[T]) Equals(t Set[T]) bool {
	if s.Len() != t.Len() {
		return false
	}

	for v := range s.data {
		if !t.Contains(v) {
			return false
		}
	}

	return true
}

// Union will create a new Set, and fill it with the union of `s` and `t`
func (s *Set[T]) Union(t Set[T]) Set[T] {
	// Figure out which is larger
	s_is_larger := s.Len() > t.Len()

	// First create a copy of either `s` or `t`. Pick whichever is largest to reduce
	// allocations.
	var result Set[T]
	if s_is_larger {
		result = s.Copy()
	} else {
		result = t.Copy()
	}

	// Iterate over the smaller set, and add all it's items to `result`
	if s_is_larger {
		for v := range t.data {
			result.Add(v)
		}
	} else {
		for v := range s.data {
			result.Add(v)
		}
	}

	return result
}

// UnionInPlace will add all the items in set `t` to set `s`
func (s *Set[T]) UnionInPlace(t Set[T]) {
	for v := range t.data {
		s.Add(v)
	}
}

// Intersection will create a new Set, and fill it with the intersection of `s` and `t`
func (s *Set[T]) Intersection(t Set[T]) Set[T] {
	// Create an empty set result
	result := NewSet([]T{})

	// Iterate over the smaller of the two sets, and add the item to `result` if it is
	// in the larger of the two sets
	if s.Len() < t.Len() {
		for v := range s.data {
			if t.Contains(v) {
				result.Add(v)
			}
		}
	} else {
		for v := range t.data {
			if s.Contains(v) {
				result.Add(v)
			}
		}
	}

	return result
}

// IntersectionInPlace will remove any items from `s` that are not in `t`
func (s *Set[T]) IntersectionInPlace(t Set[T]) {
	for v := range s.data {
		if !t.Contains(v) {
			s.Discard(v)
		}
	}
}

// IsDisjoint will return true if the set has no elements in common with `t`. Sets are
// disjoint if and only if their intersection is the empty set
func (s *Set[T]) IsDisjoint(t Set[T]) bool {
	// Iterate over the smaller of the two sets. If we find an item in one that is in
	// the other, return false
	if s.Len() < t.Len() {
		for v := range s.data {
			if t.Contains(v) {
				return false
			}
		}
	} else {
		for v := range t.data {
			if s.Contains(v) {
				return false
			}
		}
	}
	return true
}

// IsSubsetOf tests whether every element in `s` is in `t`
func (s *Set[T]) IsSubsetOf(t Set[T]) bool {
	// Iterate over `s`. If we find an item in `s` that is not in `t`, return false
	for v := range s.data {
		if !t.Contains(v) {
			return false
		}
	}
	return true
}

// IsProperSubsetOf tests whether every element in `s` is in `t`, but that
// `s.Equals(t) == false`
func (s *Set[T]) IsProperSubsetOf(t Set[T]) bool {

	// Iterate over `s`. If we find an item in `s` that is not in `t`, return false
	for v := range s.data {
		if !t.Contains(v) {
			return false
		}
	}

	// If the lengths are equal, we have just verified that the two sets are equal.
	if s.Len() == t.Len() {
		return false
	} else {
		return true
	}

}

// IsSuperSetOf tests whether every element in `t` is in `s`
func (s *Set[T]) IsSuperSetOf(t Set[T]) bool {
	// Iterate over `t`. If we find an item in `t` that is not in `s`, return false
	for v := range t.data {
		if !s.Contains(v) {
			return false
		}
	}
	return true
}

// IsProperSuperSetOf tests whether every element in `t` is in `s`, but that
// `s.Equals(t) == false`
func (s *Set[T]) IsProperSuperSetOf(t Set[T]) bool {

	// Iterate over `t`. If we find an item in `t` that is not in `s`, return false
	for v := range t.data {
		if !s.Contains(v) {
			return false
		}
	}

	// If the lengths are equal, we have just verified that the two sets are equal.
	if s.Len() == t.Len() {
		return false
	} else {
		return true
	}

}

// Difference returns a new set with elements in `s` that are not in `t`
func (s *Set[T]) Difference(t Set[T]) Set[T] {
	// Copy `s`
	result := s.Copy()

	// Iterate over `t`. If we find an item in `result`, remove it from `result`
	for v := range t.data {
		result.Discard(v)
	}

	return result
}

// DifferenceInPlace removes any elements in `s` that are in `t`
func (s *Set[T]) DifferenceInPlace(t Set[T]) {
	// Iterate over `t`. If we find an item in `s`, remove it from `s`
	for v := range t.data {
		s.Discard(v)
	}
}

// SymmetricDifference returns a new set with elements in either `s` or `t`, but not both
func (s *Set[T]) SymmetricDifference(t Set[T]) Set[T] {
	// Make an empty set to populate
	result := NewSet([]T{})

	// The big question here is whether it's worth allocating a little to save a few checks
	// For now, assume that it's best to just check everything, and store as little as
	// possible.

	// Iterate over `s`, and add the item if it does not exist in `t`
	for v := range s.data {
		if !t.Contains(v) {
			result.Add(v)
		}
	}

	// Iterate over `t`, and add the item if it does not exist in `s`
	for v := range t.data {
		if !s.Contains(v) {
			result.Add(v)
		}
	}

	return result
}

// SymmerticDifferenceInPlace removes any elements in `s` that are in `t`, and adds any
// elements in `t` that are not in `s`
func (s *Set[T]) SymmetricDifferenceInPlace(t Set[T]) {
	// Iterate over `t`. If we find an item in `s`, remove it from `s`, otherwise add it
	for v := range t.data {
		if s.Contains(v) {
			s.Discard(v)
		} else {
			s.Add(v)
		}
	}

}
