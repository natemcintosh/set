// bitset is a set of sorted ints. Implemented as a slice of bools, and therefore
// designed for dense integer sets. If the set will be sparse (for example, holding a
// few very large integers), use `github.com/natemcintosh/set` instead.
package bitset

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// This error is returned when you try to remove an item from a set that doesn't exist
	ErrElementNotFound = errors.New("element not found")
)

type Set struct {
	// bits represents a set of numbers on a number line. Positive numbers are "to the
	// right" (larger indices in the slice), and negative numbers are "to the left"
	// (smaller indices in the slice). It is always as small as possible, i.e. the
	// minimum and maximum numbers in the set are at either end of the slice (the start
	// and end of the slice will always be `true`)
	//
	// Below is an example of how a set of integers would be mapped to the slice of
	// bools
	// { -2 ,          0 ,                 3 ,   4 ,   5 }
	// [true, false, true, false, false, true, true, true]
	bits []bool

	// smallest_item tells you the distance from 0 to the smallest number in
	// the set, which is the first item in the set. The smallest number can be negative.
	// In the case of the example shown for the `bits` field, `smallest_item`
	// would be -2
	smallest_item int

	// n_items is an easy way to keep track of how many elements without having to
	// iterate over the whole `bits` slice
	n_items int
}

// NewSet will return a bitset Set from an input slice of `ints`, or anything that has a
// slice of `ints` as the underlying data type.
func NewSet[S ~[]int](data S) Set {
	if len(data) == 0 {
		return Set{
			bits:          make([]bool, 0),
			smallest_item: 0,
			n_items:       0,
		}
	}

	var (
		min int = data[0]
		max int = data[0]
	)
	unique_nums := make(map[int]struct{}, len(data))
	// Get the min and max numbers, and keep track of how many unique numbers there are
	for _, v := range data {
		if v < min {
			min = v
		}

		if v > max {
			max = v
		}

		unique_nums[v] = struct{}{}
	}

	// Make the slice
	min_max_range := abs(max-min) + 1
	if max == min {
		min_max_range = 1
	}
	bits := make([]bool, min_max_range)

	// Iterate over the input numbers, and insert `true` at the correct indices
	for _, v := range data {
		bits[v-min] = true
	}

	return Set{bits: bits, smallest_item: min, n_items: len(unique_nums)}

}

// abs returns the absolute value of x.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (s Set) String() string {
	var b strings.Builder
	last_index := len(s.bits) - 1
	b.WriteString("{")
	for idx, v := range s.bits {

		if v {
			if idx < last_index {
				b.WriteString(fmt.Sprintf("%v, ", idx+s.smallest_item))
			} else {
				b.WriteString(fmt.Sprintf("%v", idx+s.smallest_item))
			}
		}
	}
	b.WriteString("}")

	return b.String()
}

// Slice returns a sorted slice representing the integers in the set
func (s *Set) Slice() []int {
	result := make([]int, 0, s.n_items)

	// Iterate over the slice and add to `result`
	for idx, v := range s.bits {
		if v {
			result = append(result, idx+s.smallest_item)
		}
	}

	return result
}

// Contains will return true if the set contains the item. If the set is empty, returns
// false
func (s *Set) Contains(item int) bool {
	if s.IsEmpty() {
		return false
	}

	// Check if the item is outside the bounds of the slice
	if s.under_lower_bound(item) {
		return false
	} else if s.over_upper_bound(item) {
		return false
	}

	// Check if the item at the correct offset is true
	return s.bits[s.calc_idx_of_item(item)]
}

// Len returns the length of the Set
func (s *Set) Len() int {
	return s.n_items
}

// IsEmpty returns true if the set is empty
func (s *Set) IsEmpty() bool {
	return s.Len() == 0
}

func (s *Set) under_lower_bound(item int) bool {
	return item < s.smallest_item
}

func (s *Set) over_upper_bound(item int) bool {
	return item > s.get_upper_value()
}

func (s *Set) get_upper_value() int {
	return s.smallest_item + len(s.bits) - 1
}

// calc_idx_of_item does not check if your item is out of bounds. Use one of
// `under_lower_bound` or `over_upper_bound` to make sure you are in bounds. Also, if
// the returned index is negative, then it's definitely out of bounds.
func (s *Set) calc_idx_of_item(item int) int {
	return item - s.smallest_item
}

// Add will add a new item to `s`. If it already exists, it is ignored
func (s *Set) Add(item int) {
	// Check if the new element is outside the bounds
	if item < s.smallest_item {
		// Calculate how many new elements we need to have on the front
		new_items := s.smallest_item - item

		// Create a new slice that goes from the new element up to (but not including)
		// the start of the old slice
		to_add_to_front := make([]bool, new_items)
		to_add_to_front[0] = true

		// Append the old slice to the new one, and make it the `bits` field
		s.bits = append(to_add_to_front, s.bits...)

		// Update the `smallest_item` field
		s.smallest_item = item

		// Increment the length field
		s.n_items += 1

		// Return
		return
	} else if s.over_upper_bound(item) {
		// Calculate how many new elements we need to add to the end
		new_items := item - s.get_upper_value()

		// Create a new slice that goes from one after the end of the current slice to
		// the new element
		to_append := make([]bool, new_items)
		to_append[len(to_append)-1] = true

		// Append the new slice to the old
		s.bits = append(s.bits, to_append...)

		// Increment the length field
		s.n_items += 1

		// Return
		return
	}

	// Update `s.n_items` if necessary
	if !s.Contains(item) {
		s.n_items += 1
	}

	// Set the item at the correct index to true
	s.bits[s.calc_idx_of_item(item)] = true
}

// Remove removes an item from the set. Returns an error if the item doesn't exist
func (s *Set) Remove(item int) error {
	if !s.Contains(item) {
		return ErrElementNotFound
	}

	// Remove the item
	s.bits[s.calc_idx_of_item(item)] = false

	// Decrement the number of items field
	s.n_items -= 1

	// Was the value removed the smallest value?
	if item == s.smallest_item {
		// Make the slice smaller

		// Find the next true index in the slice, and make that the start
		for idx, v := range s.bits {
			if v {
				s.bits = s.bits[idx:]
				s.smallest_item = item + idx
				break
			}
		}

	} else if item == s.get_upper_value() {
		// Was the value removed the largest value?
		// Make the slice smaller

		// Find the index of the next true index in the slice from the rear, and make
		// that the end
		for idx := len(s.bits) - 1; idx >= 0; idx-- {
			if s.bits[idx] {
				// Keep everything up to and including this index
				s.bits = s.bits[:(idx + 1)]
				break
			}
		}

	}

	return nil
}

// Discard removes an item from the set. If it doesn't exist, it is ignored
func (s *Set) Discard(item int) {
	if !s.Contains(item) {
		return
	}

	// Remove the item
	s.bits[s.calc_idx_of_item(item)] = false

	// Decrement the number of items field
	s.n_items -= 1

	// Was the value removed the smallest value?
	if item == s.smallest_item {
		// Make the slice smaller

		// Find the next true index in the slice, and make that the start
		for idx, v := range s.bits {
			if v {
				s.bits = s.bits[idx:]
				s.smallest_item = item + idx
				break
			}
		}

	} else if item == s.get_upper_value() {
		// Was the value removed the largest value?
		// Make the slice smaller

		// Find the index of the next true index in the slice from the rear, and make
		// that the end
		for idx := len(s.bits) - 1; idx >= 0; idx-- {
			if s.bits[idx] {
				// Keep everything up to and including this index
				s.bits = s.bits[:(idx + 1)]
				break
			}
		}

	}
}

// Equals will return true if `s` and `t` are
// - the same length
// - contain the same elements
func (s *Set) Equals(t Set) bool {
	if s.Len() != t.Len() {
		return false
	}

	// If they don't have the same start, they are not equal
	if s.smallest_item != t.smallest_item {
		return false
	}

	for idx, v := range s.bits {
		if t.bits[idx] != v {
			return false
		}
	}

	return true
}
