// bitset is a set of sorted ints. Implemented as a slice of bools, and therefore
// designed for dense integer sets. If the set will be sparse (for example, holding a
// few very large integers), use `github.com/natemcintosh/set` instead.
package bitset

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

	// min_item_offset_from_0 tells you the distance from 0 to the smallest number in
	// the set, which is the first item in the set. The smallest number can be negative.
	// In the case of the example shown for the `bits` field, `min_item_offset_from_0`
	// would be -2
	min_item_offset_from_0 int

	// length is an easy way to keep track of how many elements without having to
	// iterate over the whole `bits` slice
	length int
}

// NewSet will return a bitset Set from an input slice of `ints`, or anything that has a
// slice of `ints` as the underlying data type.
func NewSet[S ~[]int](data S) Set {
	if len(data) == 0 {
		return Set{}
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

	return Set{bits: bits, min_item_offset_from_0: min, length: len(unique_nums)}

}

// Slice returns a sorted slice representing the integers in the set
func (s *Set) Slice() []int {
	result := make([]int, 0, s.length)

	// Iterate over the slice and add to `result`
	for idx, v := range s.bits {
		if v {
			result = append(result, idx+s.min_item_offset_from_0)
		}
	}

	return result
}

// abs returns the absolute value of x.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
