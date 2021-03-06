// bitset is a set of sorted ints. Implemented as a slice of bools, and therefore
// designed for dense integer sets. If the set will be sparse (for example, holding a
// few very large integers), use `github.com/natemcintosh/set` instead.
package bitset

import (
	"errors"
	"fmt"
	"math/bits"
	"strings"
)

var (
	// This error is returned when you try to remove an item from a set that doesn't exist
	ErrElementNotFound = errors.New("element not found")
)

// abs returns the absolute value of x.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type key struct {
	is_positive bool
	multiplier  uint64
}

type Set struct {
	data map[key]uint64
}

func NewSet[S ~[]int](data S) Set {
	// Create the underlying set
	uset := make(map[key]uint64)

	for _, v := range data {
		// Get the new data representation
		is_positive, multiplier, slot := number_to_bitset_representation(v)

		// Create the key for the map
		key := key{
			is_positive: is_positive,
			multiplier:  multiplier,
		}

		// Union if it already exists, else just add it
		if bits, ok := uset[key]; ok {
			uset[key] = bits | slot
		} else {
			uset[key] = slot
		}
	}

	return Set{data: uset}
}

// number_to_bitset_representation will take an int and return the following
//
// - `is_positive`: true if n >= 0
//
// - `multiplier`: how many times 64 goes into n: abs(n) / 64
//
// - `slot`: using an uint64 to represent the 64 bins for the remainder: abs(n) % 64
func number_to_bitset_representation(n int) (
	is_positive bool,
	multiplier uint64,
	slot uint64,
) {
	if n >= 0 {
		is_positive = true
		multiplier = uint64(n) / 64
		if n%64 == 0 {
			slot = uint64(1)
		} else {
			slot = two_to_power_n_minus_1(n % 64)
		}
	} else {
		is_positive = false
		multiplier = uint64(-n) / 64
		if -n%64 == 0 {
			slot = 1
		} else {
			slot = two_to_power_n_minus_1(-n % 64)
		}
	}
	return
}

func two_to_power_n_minus_1(n int) uint64 {
	return 1 << uint64(n)
}

func (u Set) String() string {
	var b strings.Builder
	b.WriteRune('{')
	upper_idx := -1
	for key, bits := range u.data {
		upper_idx += 1
		// For each bit in `bits`, want to extract the index of the bit if it is 1
		// and then add it to the string
		m := 64 * key.multiplier
		vals := slots_from_uint64(bits)
		for idx, v := range vals {
			val := m + uint64(v)
			if key.is_positive {
				if (idx+1)+(upper_idx+1) != len(vals)+len(u.data) {
					b.WriteString(fmt.Sprintf("%d, ", val))
				} else {
					b.WriteString(fmt.Sprintf("%d", val))
				}
			} else {
				if (idx+1)+(upper_idx+1) != len(vals)+len(u.data) {
					b.WriteString(fmt.Sprintf("-%d, ", val))
				} else {
					b.WriteString(fmt.Sprintf("-%d", val))
				}
			}
		}
	}
	b.WriteRune('}')
	return b.String()
}

func slots_from_uint64(u uint64) []int {
	if u == 0 {
		return []int{0}
	}
	var idx int
	result := make([]int, 0, bits.OnesCount64(u))
	for u != 0 {
		idx = bits.TrailingZeros64(u)
		result = append(result, idx)
		u &= ^(1 << uint64(idx))
	}

	return result
}

// Slice will return all the items in the set as a slice. They are not guaranteed in any
// particular order.
func (s *Set) Slice() []int {
	result := make([]int, 0)
	for key, bits := range s.data {
		// For each bit in `bits`, want to extract the index of the bit if it is 1
		// and then add it to the string
		m := 64 * int(key.multiplier)
		vals := slots_from_uint64(bits)
		for _, v := range vals {
			val := m + v
			if key.is_positive {
				result = append(result, val)
			} else {
				result = append(result, -val)
			}
		}
	}
	return result
}

// Contains will return true if the set contains the item. If the set is empty, returns
// false
func (s *Set) Contains(item int) bool {
	if len(s.data) == 0 {
		return false
	}

	// Get the new data representation
	is_positive, multiplier, slot := number_to_bitset_representation(item)

	key := key{is_positive: is_positive, multiplier: multiplier}

	if bits, ok := s.data[key]; ok {
		if bits&slot != 0 {
			return true
		}
	}

	return false

}

// Len returns the length of the Set
func (s *Set) Len() int {
	res := 0
	for _, v := range s.data {
		res += bits.OnesCount64(v)
	}
	return res
}

// IsEmpty returns true if the set is empty
func (s *Set) IsEmpty() bool {
	return s.Len() == 0
}

// Add will add a new item to `s`. If it already exists, it is ignored
func (s *Set) Add(item int) {
	// Get the new data representation
	is_positive, multiplier, slot := number_to_bitset_representation(item)

	key := key{is_positive: is_positive, multiplier: multiplier}

	// Union if it already exists, else just add it
	if bits, ok := s.data[key]; ok {
		s.data[key] = bits | slot
	} else {
		s.data[key] = slot
	}
}

// Remove removes an item from the set. Returns an error if the item doesn't exist
func (s *Set) Remove(item int) error {
	if len(s.data) == 0 {
		return ErrElementNotFound
	}

	// Get the new data representation
	is_positive, multiplier, slot := number_to_bitset_representation(item)

	key := key{is_positive: is_positive, multiplier: multiplier}

	if bits, ok := s.data[key]; !ok {
		// This uint64 doesn't exist in the map
		return ErrElementNotFound
	} else {
		if bits&slot == 0 {
			// Was not found in this uint64
			return ErrElementNotFound
		}
		// Remove the element
		s.data[key] = bits ^ slot
	}
	return nil
}

// Discard removes an item from the set. If it doesn't exist, it is ignored
func (s *Set) Discard(item int) {
	if len(s.data) == 0 {
		return
	}

	// Get the new data representation
	is_positive, multiplier, slot := number_to_bitset_representation(item)

	key := key{is_positive: is_positive, multiplier: multiplier}

	if bits, ok := s.data[key]; !ok {
		// This uint64 doesn't exist in the map
		return
	} else {
		// Remove the element
		s.data[key] = bits ^ slot
	}
	return
}

// Pop will remove and return an arbitrary item from the set. If the set is empty,
// it will return an error
func (s *Set) Pop() (item int, err error) {
	if s.IsEmpty() {
		return item, ErrElementNotFound
	}

	var to_return int
	// Iterate to the first item
	for key, slots := range s.data {
		to_return = bits.TrailingZeros64(slots)
		// Erase that bit
		s.data[key] &= ^(1 << uint(to_return))
		break
	}

	return to_return, nil

}

// Clear will remove all items from the set
func (s *Set) Clear() {
	s.data = make(map[key]uint64)
}

// Copy makes a deep copy as quickly as possible
func (s *Set) Copy() Set {
	// Make sure to allocate the same size
	copy := make(map[key]uint64, len(s.data))

	// Fill it up
	for key, slots := range s.data {
		copy[key] = slots
	}

	return Set{data: copy}
}

// Equals will return true if `s` and `t` are
// - the same length
// - contain the same elements
func (s *Set) Equals(t Set) bool {
	if len(s.data) != len(t.data) {
		return false
	}

	if s.Len() != t.Len() {
		return false
	}

	for skey, sbits := range s.data {
		if tbits, ok := t.data[skey]; !ok {
			return false
		} else if sbits != tbits {
			return false
		}
	}

	// We've checked that all keys in `s` are in `t`, but not the other way around
	for tkey := range t.data {
		if _, ok := s.data[tkey]; !ok {
			return false
		}
	}

	return true
}

// Union will create a new Set, and fill it with the union of `s` and `t`
func (s *Set) Union(t Set) Set {
	// Figure out which is has more key->value pairs
	s_is_larger := len(s.data) > len(t.data)

	// First create a copy of either `s` or `t`. Pick whichever is largest to reduce
	// allocations.
	var result Set
	if s_is_larger {
		result = s.Copy()
	} else {
		result = t.Copy()
	}

	// Iterate over the smaller set, and add all it's items to `result`
	if s_is_larger {
		for tkey, tslots := range t.data {
			// Get the key from s (if it exists)
			if sslots, ok := result.data[tkey]; ok {
				result.data[tkey] = sslots | tslots
			} else {
				result.data[tkey] = tslots
			}
		}
	} else {
		for skey, sslots := range s.data {
			// Get the key from t (if it exists)
			if tslots, ok := result.data[skey]; ok {
				result.data[skey] = sslots | tslots
			} else {
				result.data[skey] = sslots
			}
		}
	}

	return result
}

// UnionInPlace will add all the items in set `t` to set `s`
func (s *Set) UnionInPlace(t Set) {
	for tkey, tslots := range t.data {
		// Get the key from s (if it exists)
		if sslots, ok := s.data[tkey]; ok {
			s.data[tkey] = sslots | tslots
		} else {
			s.data[tkey] = tslots
		}
	}
}

// Intersection will create a new Set, and fill it with the intersection of `s` and `t`
func (s *Set) Intersection(t Set) Set {
	// Create an empty set result
	data := make(map[key]uint64)

	// Iterate over the smaller of the two sets, and add the item to `result` if it is
	// in the larger of the two sets
	if len(s.data) < len(t.data) {
		for skey, sslots := range s.data {
			// Get the key from t (if it exists)
			if tslots, ok := t.data[skey]; ok {
				if (sslots & tslots) != 0 {
					data[skey] = sslots & tslots
				}
			}
		}
	} else {
		for tkey, tslots := range t.data {
			// Get the key from s (if it exists)
			if sslots, ok := s.data[tkey]; ok {
				if (sslots & tslots) != 0 {
					data[tkey] = sslots & tslots
				}
			}
		}
	}

	return Set{data: data}
}

// IntersectionInPlace will remove any items from `s` that are not in `t`
func (s *Set) IntersectionInPlace(t Set) {
	// For each key in `s`, check if it is in `t`
	// If it is, perform the intersection, else remove it
	// If the intersection is 0, remove the key
	for skey, sslots := range s.data {
		// Get the key from t (if it exists)
		if tslots, ok := t.data[skey]; ok {
			if (sslots & tslots) == 0 {
				delete(s.data, skey)
			} else {
				s.data[skey] = sslots & tslots
			}
		} else {
			// The key does not exist in `t`, so remove it from `s`
			delete(s.data, skey)
		}
	}
}

// IsDisjoint will return true if the set has no elements in common with `t`. Sets are
// disjoint if and only if their intersection is the empty set
func (s *Set) IsDisjoint(t Set) bool {
	// Iterate over the smaller of the two sets. If we find an item in one that is in
	// the other, return false
	if len(s.data) < len(t.data) {
		for skey, sslots := range s.data {
			// Get the key from t (if it exists)
			if tslots, ok := t.data[skey]; ok {
				if (sslots & tslots) != 0 {
					return false
				}
			}
		}
	} else {
		for tkey, tslots := range t.data {
			// Get the key from s (if it exists)
			if sslots, ok := s.data[tkey]; ok {
				if (sslots & tslots) != 0 {
					return false
				}
			}

		}
	}
	return true
}

// IsSubsetOf tests whether every element in `s` is in `t`
func (s *Set) IsSubsetOf(t Set) bool {
	// Iterate over `s`. If we find an item in `s` that is not in `t`, return false
	for skey, sslots := range s.data {
		// Get the key from t (if it exists)
		if tslots, ok := t.data[skey]; ok {
			if (sslots & tslots) != sslots {
				return false
			}
		} else {
			// The key does not exist in `t`, so return false
			return false
		}
	}
	return true
}

// IsProperSubsetOf tests whether every element in `s` is in `t`, but that
// `s.Equals(t) == false`
func (s *Set) IsProperSubsetOf(t Set) bool {

	// Iterate over `s`. If we find an item in `s` that is not in `t`, return false
	for skey, sslots := range s.data {
		// Get the key from t (if it exists)
		if tslots, ok := t.data[skey]; ok {
			if (sslots & tslots) != sslots {
				return false
			}
		} else {
			// The key does not exist in `t`, so return false
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
func (s *Set) IsSuperSetOf(t Set) bool {
	// Iterate over `t`. If we find an item in `t` that is not in `s`, return false
	for tkey, tslots := range t.data {
		// Get the key from s (if it exists)
		if sslots, ok := s.data[tkey]; ok {
			if (sslots & tslots) != tslots {
				return false
			}
		} else {
			// The key does not exist in `s`, so return false
			return false
		}
	}
	return true
}

// IsProperSuperSetOf tests whether every element in `t` is in `s`, but that
// `s.Equals(t) == false`
func (s *Set) IsProperSuperSetOf(t Set) bool {

	// Iterate over `t`. If we find an item in `t` that is not in `s`, return false
	for tkey, tslots := range t.data {
		// Get the key from s (if it exists)
		if sslots, ok := s.data[tkey]; ok {
			if (sslots & tslots) != tslots {
				return false
			}
		} else {
			// The key does not exist in `s`, so return false
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
func (s *Set) Difference(t Set) Set {
	// Copy `s`
	result := s.Copy()

	// Iterate over `t`. If we find an item in `result`, remove it from `result`
	for tkey, tslots := range t.data {
		// Get the key from result (if it exists)
		if sslots, ok := result.data[tkey]; ok {
			// Make sslots the intersection of sslots and tslots
			if sslots^tslots == 0 {
				delete(result.data, tkey)
			} else {
				result.data[tkey] = sslots &^ tslots
			}
		}
	}

	return result
}

// DifferenceInPlace removes any elements in `s` that are in `t`
func (s *Set) DifferenceInPlace(t Set) {
	// Iterate over `t`. If we find an item in `s`, remove it from `s`
	for tkey, tslots := range t.data {
		// Get the key from s (if it exists)
		if sslots, ok := s.data[tkey]; ok {
			// Make sslots the intersection of sslots and tslots
			if sslots^tslots == 0 {
				delete(s.data, tkey)
			} else {
				s.data[tkey] = sslots &^ tslots
			}
		}
	}
}

// SymmetricDifference returns a new set with elements in either `s` or `t`, but not both
func (s *Set) SymmetricDifference(t Set) Set {
	// Make an empty set to populate
	data := make(map[key]uint64)

	// Iterate over `s`, and add the item if it does not exist in `t`
	for skey, sslots := range s.data {
		// Get the key from t (if it exists)
		if tslots, ok := t.data[skey]; ok {
			// Make sslots the intersection of sslots and tslots
			if sslots^tslots == 0 {
				continue
			} else {
				data[skey] = sslots ^ tslots
			}
		} else {
			// The key does not exist in `t`, so add it to the result
			data[skey] = sslots
		}
	}

	// Iterate over `t`, and add the item if it does not exist in `s`
	for tkey, tslots := range t.data {
		// Get the key from s (if it exists)
		if sslots, ok := s.data[tkey]; ok {
			// Make sslots the intersection of sslots and tslots
			if sslots^tslots == 0 {
				continue
			} else {
				data[tkey] = sslots ^ tslots
			}
		} else {
			// The key does not exist in `s`, so add it to the result
			data[tkey] = tslots
		}
	}

	return Set{data: data}
}

// SymmerticDifferenceInPlace removes any elements in `s` that are in `t`, and adds any
// elements in `t` that are not in `s`
func (s *Set) SymmetricDifferenceInPlace(t Set) {
	// Iterate over `t`. If we find an item in `s`, remove it from `s`, otherwise add it
	for tkey, tslots := range t.data {
		// Get the key from s (if it exists)
		if sslots, ok := s.data[tkey]; ok {
			// Make sslots the intersection of sslots and tslots
			if sslots^tslots == 0 {
				delete(s.data, tkey)
			} else {
				s.data[tkey] = sslots ^ tslots
			}
		} else {
			// The key does not exist in `s`, so add it to the result
			s.data[tkey] = tslots
		}
	}

}
