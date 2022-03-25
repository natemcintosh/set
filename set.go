// set is a library designed to help you do set operations on comparable data types
package set

import "errors"

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

// Equals will return true if `s` and `t` are
// - the same length
// - contain the same elements
func (s *Set[T]) Equals(t Set[T]) bool {
	if s.Len() != t.Len() {
		return false
	}

	for v := range s.data {
		if _, ok := t.data[v]; !ok {
			return false
		}
	}

	return true
}

// IsEmpty returns true if the set is empty
func (s *Set[T]) IsEmpty() bool {
	return s.Len() == 0
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

// Len returns the length of the Set
func (s *Set[T]) Len() int {
	return len(s.data)
}

// Add will add a new item to `s`. If it already exists, it is ignored
func (s *Set[T]) Add(element T) {
	s.data[element] = struct{}{}
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

// Remove removes an item from the set. Returns an error if the item doesn't exist
func (s *Set[T]) Remove(element T) error {
	if _, ok := s.data[element]; !ok {
		return ErrElementNotFound
	}

	delete(s.data, element)
	return nil
}

// Discard removes an item from the set. If it doesn't exist, it is ignored
func (s *Set[T]) Discard(element T) {
	delete(s.data, element)
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

// Contains will return true if the set contains the item. If the set is empty, returns
// false
func (s *Set[T]) Contains(element T) bool {
	_, ok := s.data[element]
	return ok
}

// Intersection will create a new Set, and fill it with the intersection of `s` and `t`
func (s *Set[T]) Intersection(t Set[T]) Set[T] {
	// Create an empty set result
	result := NewSet([]T{})

	// Iterate over the smaller of the two sets, and add the item to `result` if it is
	// in the larger of the two sets
	if s.Len() < t.Len() {
		for v := range s.data {
			if _, ok := t.data[v]; ok {
				result.Add(v)
			}
		}
	} else {
		for v := range t.data {
			if _, ok := s.data[v]; ok {
				result.Add(v)
			}
		}
	}

	return result
}

// IntersectionInPlace will remove any items from `s` that are not in `t`
func (s *Set[T]) IntersectionInPlace(t Set[T]) {
	for v := range s.data {
		if _, ok := t.data[v]; !ok {
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
			if _, ok := t.data[v]; ok {
				return false
			}
		}
	} else {
		for v := range t.data {
			if _, ok := s.data[v]; ok {
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
		if _, ok := t.data[v]; !ok {
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
		if _, ok := t.data[v]; !ok {
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
