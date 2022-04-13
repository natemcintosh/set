package bitset

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/natemcintosh/set"
	"golang.org/x/exp/slices"
)

func TestString(t *testing.T) {
	testCases := []struct {
		desc string
		set  []int
	}{
		{
			desc: "empty",
			set:  []int{},
		},
		{
			desc: "zero",
			set:  []int{0},
		},
		{
			desc: "one",
			set:  []int{1},
		},
		{
			desc: "a few",
			set:  []int{-1, 3, 10},
		},
		{
			desc: "comes out sorted",
			set:  []int{1, -1},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// Create the set
			set := NewSet(tC.set)

			// Create the string
			s := fmt.Sprintf("%v", set)

			// For each item in `tC.set`, see if it exists in the string
			for _, v := range tC.set {
				if !strings.Contains(s, fmt.Sprintf("%d", v)) {
					t.Errorf("Could not find %d in the set %v", v, s)
				}
			}
		})
	}
}

func TestSlots_from_uint64(t *testing.T) {
	testCases := []struct {
		desc string
		in   uint64
		want []int
	}{
		{
			desc: "0",
			in:   1,
			want: []int{0},
		},
		{
			desc: "1",
			in:   2,
			want: []int{1},
		},
		{
			desc: "0,1,2",
			in:   7,
			want: []int{0, 1, 2},
		},
		{
			desc: "0,2",
			in:   5,
			want: []int{0, 2},
		},
		{
			desc: "63",
			in:   1 << 63,
			want: []int{63},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := slots_from_uint64(tC.in)

			if !slices.Equal(got, tC.want) {
				t.Errorf("got %v; want %v", got, tC.want)
			}
		})
	}
}

func FuzzSlots_from_uint64(f *testing.F) {
	// This fuzz test is for checking that we don't hit any panics in getting the slots
	// where 1s are stored in a uint64
	f.Add(uint64(2))
	f.Add(uint64(10))

	f.Fuzz(func(t *testing.T, i uint64) {
		// When we run the function, does anything funny happen?
		slots_from_uint64(i)
	})
}

func TestConvertBackAndForth(t *testing.T) {
	testCases := []struct {
		desc string
		s    []int
	}{
		{
			desc: "empty",
			s:    []int{},
		},
		{
			desc: "zero",
			s:    []int{0},
		},
		{
			desc: "one",
			s:    []int{1},
		},
		{
			desc: "small",
			s:    []int{-1, 0, 1},
		},
		{
			desc: "medium",
			s:    []int{-2, 0, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			desc: "with duplicates",
			s:    []int{0, 0, 1, 1, 2},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// Create a bitset
			newset := NewSet(tC.s)

			// Get the numbers back out
			got := newset.Slice()
			slices.Sort(got)

			// Make sure to remove any duplicates
			want := make([]int, len(tC.s))
			copy(want, tC.s)
			slices.Sort(want)
			want = slices.Compact(want)

			// Compare the returned numbers to the originals
			if !equal(got, want) {
				t.Errorf("got %v; want %v", got, want)
			}
		})
	}
}

func FuzzConvertBackAndForth(f *testing.F) {
	f.Add(-2, 0, 3, 4, 5, 6, 7, 8, 9, 10)
	f.Add(-10, -4, -5, -11, -20, 12, 16, 13, 34, 35)

	f.Fuzz(func(
		t *testing.T,
		s1 int,
		s2 int,
		s3 int,
		s4 int,
		s5 int,
		s6 int,
		s7 int,
		s8 int,
		s9 int,
		s10 int,
	) {
		// When we get the numbers back, they will be sorted. Sort them here so we can
		// compare more easily later.
		// Remove any duplicates, becuase the set will not have duplicates
		want := []int{s1, s2, s3, s4, s5, s6, s7, s8, s9, s10}

		// Create a bitset
		newset := NewSet(want)

		// Get the numbers back out
		got := newset.Slice()
		slices.Sort(got)
		slices.Sort(want)
		want = slices.Compact(want)

		// Compare the returned numbers to the originals
		if !equal(want, got) {
			t.Errorf("got %v; want %v", got, want)
		}
	})
}

func BenchmarkConvertBackAndForth(b *testing.B) {
	benchCases := []struct {
		desc string
		s    []int
	}{
		{
			desc: "empty",
			s:    []int{},
		},
		{
			desc: "zero",
			s:    []int{0},
		},
		{
			desc: "one",
			s:    []int{1},
		},
		{
			desc: "small",
			s:    []int{-1, 0, 1},
		},
		{
			desc: "medium",
			s:    []int{-2, 0, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			desc: "with duplicates",
			s:    []int{0, 0, 1, 1, 2},
		},
	}
	for _, bC := range benchCases {
		b.Run(bC.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				set := NewSet(bC.s)
				set.Slice()
			}
		})
	}
}

func FuzzLen(f *testing.F) {
	// This fuzz test is for checking that Len() is always correct
	f.Add(2)
	f.Add(10)

	f.Fuzz(func(t *testing.T, _n int) {
		n := abs(_n)
		items := make([]int, n)
		// Create n random ints
		for i := 0; i < n; i++ {
			items[i] = rand.Int()
		}

		// Create the set
		set := NewSet(items)

		// Figure out how many unique items there
		slices.Sort(items)
		items = slices.Compact(items)
		want := len(items)

		// Check that Len is correct
		if set.Len() != want {
			t.Errorf("got %d, want %d. Set = %v", set.Len(), want, set)
		}
	})
}

// equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestContains(t *testing.T) {
	testCases := []struct {
		desc string
		set  Set
		item int
		want bool
	}{
		{
			desc: "empty",
			set:  NewSet([]int{}),
			item: 0,
			want: false,
		},
		{
			desc: "before",
			set:  NewSet([]int{1, 2, 3}),
			item: 0,
			want: false,
		},
		{
			desc: "after",
			set:  NewSet([]int{1, 2, 3}),
			item: 4,
			want: false,
		},
		{
			desc: "inside",
			set:  NewSet([]int{1, 2, 3}),
			item: 2,
			want: true,
		},
		{
			desc: "near outer edge",
			set:  NewSet([]int{0, 9, 10}),
			item: 9,
			want: true,
		},
		{
			desc: "near outer edge 2",
			set:  NewSet([]int{-2, -1, 5}),
			item: 5,
			want: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := tC.set.Contains(tC.item)

			if got != tC.want {
				t.Errorf("got %v, want %v", got, tC.want)
			}
		})
	}
}

func BenchmarkContains(b *testing.B) {
	benchCases := []struct {
		desc string
		set  Set
		item int
	}{
		{
			desc: "empty",
			set:  NewSet([]int{}),
			item: 0,
		},
		{
			desc: "before",
			set:  NewSet([]int{1, 2, 3}),
			item: 0,
		},
		{
			desc: "after",
			set:  NewSet([]int{1, 2, 3}),
			item: 4,
		},
		{
			desc: "inside",
			set:  NewSet([]int{1, 2, 3}),
			item: 2,
		},
		{
			desc: "near outer edge",
			set:  NewSet([]int{0, 9, 10}),
			item: 9,
		},
		{
			desc: "near outer edge 2",
			set:  NewSet([]int{-2, -1, 5}),
			item: 5,
		},
	}
	for _, bC := range benchCases {
		b.Run(bC.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bC.set.Contains(bC.item)
			}
		})
	}
}

func FuzzContains(f *testing.F) {
	// We are hoping to find cases where bitset does not match set
	f.Add(2)
	f.Add(10)

	f.Fuzz(func(t *testing.T, _n int) {

		n := abs(_n)
		items := make([]int, n)
		// Create n random ints
		for i := 0; i < n; i++ {
			if _n > 0 {
				items[i] = rand.Int()
			} else {
				items[i] = -rand.Int()
			}
		}
		// Create the bitset
		set := set.NewSet(items)
		bitset := NewSet(items)

		// Run Contains, hope to find differences between the two sets
		to_find := rand.Int()
		bitset_found := bitset.Contains(to_find)
		set_found := set.Contains(to_find)

		if bitset_found != set_found {
			t.Errorf("bitset and set Contains() did not match")
		}

	})
}

func TestAdd(t *testing.T) {
	testCases := []struct {
		desc        string
		set         Set
		new_item    int
		should_grow bool
	}{
		{
			desc:        "before",
			set:         NewSet([]int{1, 2, 3}),
			new_item:    0,
			should_grow: true,
		},
		{
			desc:        "way before",
			set:         NewSet([]int{1, 2, 3}),
			new_item:    -100,
			should_grow: true,
		},
		{
			desc:        "after",
			set:         NewSet([]int{1, 2, 3}),
			new_item:    4,
			should_grow: true,
		},
		{
			desc:        "way after",
			set:         NewSet([]int{1, 2, 3}),
			new_item:    100,
			should_grow: true,
		},
		{
			desc:        "in range but not member yet",
			set:         NewSet([]int{1, 2, 4}),
			new_item:    3,
			should_grow: true,
		},
		{
			desc:        "in range and member",
			set:         NewSet([]int{1, 2, 3}),
			new_item:    2,
			should_grow: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// Get the current number of items
			n_items := tC.set.Len()

			// Add the new_item
			tC.set.Add(tC.new_item)

			// Check that the set contains the new item
			if !tC.set.Contains(tC.new_item) {
				t.Errorf("set %v should contain item %v", tC.set, tC.new_item)
			}

			// Check if it grew if it was supposed to, or not if not
			if tC.should_grow && (tC.set.Len() != (n_items + 1)) {
				t.Errorf("set %v adding item %v should have grown, but did not", tC.set, tC.new_item)
			} else if !tC.should_grow && (tC.set.Len() != n_items) {
				t.Errorf("set %v adding item %v should not have grown, but length changed", tC.set, tC.new_item)
			}
		})
	}
}

func BenchmarkAdd(b *testing.B) {
	benchCases := []struct {
		desc     string
		set      Set
		new_item int
	}{
		{
			desc:     "before",
			set:      NewSet([]int{1, 2, 3}),
			new_item: 0,
		},
		{
			desc:     "way before",
			set:      NewSet([]int{1, 2, 3}),
			new_item: -100,
		},
		{
			desc:     "after",
			set:      NewSet([]int{1, 2, 3}),
			new_item: 4,
		},
		{
			desc:     "way after",
			set:      NewSet([]int{1, 2, 3}),
			new_item: 100,
		},
		{
			desc:     "in range but not member yet",
			set:      NewSet([]int{1, 2, 4}),
			new_item: 3,
		},
		{
			desc:     "in range and member",
			set:      NewSet([]int{1, 2, 3}),
			new_item: 2,
		},
	}
	for _, bC := range benchCases {
		b.Run(bC.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bC.set.Add(bC.new_item)
			}
		})
	}
}

func FuzzAdd(f *testing.F) {
	// We are hoping to find panics with this fuzz test
	f.Add(2)
	f.Add(10)

	f.Fuzz(func(t *testing.T, _n int) {
		n := abs(_n)
		items := make([]int, n)
		// Create n random ints
		for i := 0; i < n; i++ {
			if _n > 0 {
				items[i] = rand.Int()
			} else {
				items[i] = -rand.Int()
			}
		}
		// Create the bitset
		set := set.NewSet(items)
		bitset := NewSet(items)

		// Add a new item
		to_add := rand.Int()
		set.Add(to_add)
		bitset.Add(to_add)

		// Check that both sets are the same
		setslice := set.Slice()
		bitsetslice := bitset.Slice()

		slices.Sort(setslice)
		slices.Sort(bitsetslice)

		if !equal(setslice, bitsetslice) {
			t.Errorf("set and bitset did not match after adding element %d", to_add)
		}

	})
}

func TestRemove(t *testing.T) {
	testCases := []struct {
		desc           string
		s              Set
		v              int
		want_set       Set
		want_err_value error
	}{
		{
			desc:           "valid remove",
			s:              NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			v:              5,
			want_set:       NewSet([]int{1, 2, 3, 4, 6, 7, 8, 9, 10}),
			want_err_value: nil,
		},
		{
			desc:           "invalid remove",
			s:              NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			v:              11,
			want_set:       NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			want_err_value: ErrElementNotFound,
		},
		{
			desc:           "smallest item",
			s:              NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			v:              1,
			want_set:       NewSet([]int{2, 3, 4, 5, 6, 7, 8, 9, 10}),
			want_err_value: nil,
		},
		{
			desc:           "largest item",
			s:              NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			v:              10,
			want_set:       NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			want_err_value: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := tC.s.Remove(tC.v)
			if err != tC.want_err_value {
				t.Errorf("got error %v, want %v", err, tC.want_err_value)
			}
			if !tC.s.Equals(tC.want_set) {
				t.Errorf("got %v, want %v", tC.s, tC.want_set)
			}

		})
	}
}

func FuzzRemoveDiscard(f *testing.F) {
	// We are hoping to find places where
	f.Add(2)
	f.Add(10)

	f.Fuzz(func(t *testing.T, _n int) {
		n := abs(_n)
		items := make([]int, n)
		// Create n random ints
		for i := 0; i < n; i++ {
			if _n > 0 {
				items[i] = rand.Int()
			} else {
				items[i] = -rand.Int()
			}
		}
		// Create the bitset
		set := set.NewSet(items)
		bitset := NewSet(items)

		// Remove a new item
		to_remove := rand.Int()
		set.Remove(to_remove)
		bitset.Remove(to_remove)

		// Check that both sets are the same
		setslice := set.Slice()
		bitsetslice := bitset.Slice()

		slices.Sort(setslice)
		slices.Sort(bitsetslice)

		if !equal(setslice, bitsetslice) {
			t.Errorf("set and bitset did not match after removing element %d", to_remove)
		}

		// Discard an item from both sets
		to_discard := rand.Int()
		set.Discard(to_discard)
		bitset.Discard(to_discard)

		// Check that both sets are the same
		setslice = set.Slice()
		bitsetslice = bitset.Slice()

		slices.Sort(setslice)
		slices.Sort(bitsetslice)

		if !equal(setslice, bitsetslice) {
			t.Errorf("set and bitset did not match after discarding element %d", to_discard)
		}

	})
}

func BenchmarkRemove(b *testing.B) {
	benchCases := []struct {
		desc string
		s    Set
		v    int
	}{
		{
			desc: "valid remove",
			s:    NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			v:    5,
		},
		{
			desc: "invalid remove",
			s:    NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			v:    11,
		},
		{
			desc: "smallest item",
			s:    NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			v:    1,
		},
		{
			desc: "largest item",
			s:    NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			v:    10,
		},
	}
	for _, bC := range benchCases {
		b.Run(bC.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bC.s.Remove(bC.v)
			}
		})
	}
}

func BenchmarkDiscard(b *testing.B) {
	benchCases := []struct {
		desc string
		s    Set
		v    int
	}{
		{
			desc: "valid remove",
			s:    NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			v:    5,
		},
		{
			desc: "invalid remove",
			s:    NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			v:    11,
		},
		{
			desc: "smallest item",
			s:    NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			v:    1,
		},
		{
			desc: "largest item",
			s:    NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			v:    10,
		},
	}
	for _, bC := range benchCases {
		b.Run(bC.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bC.s.Discard(bC.v)
			}
		})
	}
}

func TestPop(t *testing.T) {
	testCases := []struct {
		desc     string
		s        Set
		want_err error
	}{
		{
			desc:     "valid pop",
			s:        NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			want_err: nil,
		},
		{
			desc:     "invalid pop",
			s:        NewSet([]int{}),
			want_err: ErrElementNotFound,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			original_len := tC.s.Len()
			_, err := tC.s.Pop()
			if err != tC.want_err {
				t.Errorf("got error %v, want %v", err, tC.want_err)
			}
			// if the error is nil, check that the length is one less than before
			if err == nil && tC.s.Len() != original_len-1 {
				t.Errorf("got %v, want %v", tC.s.Len(), original_len-1)
			}
		})
	}
}

func TestClear(t *testing.T) {
	s := NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	s.Clear()
	if !s.IsEmpty() {
		t.Errorf("got %v, want empty", s)
	}
}

func BenchmarkMonteCarloRuns(b *testing.B) {
	// Create a set of numbers from 1 to 1,000
	mcslice := make([]int, 1000)
	// Fill it with numbers from 1 to 1,000
	for i := 0; i < 1000; i++ {
		mcslice[i] = i + 1
	}
	// Create a set from the slice
	mcs := NewSet(mcslice)

	// Create a set that is a subset of `mcs`
	mcs_subset := mcs.Copy()
	mcs_subset.Discard(1)
	mcs_subset.Discard(20)
	mcs_subset.Discard(50)
	mcs_subset.Discard(143)
	mcs_subset.Discard(999)

	// Reset the benchmark timer
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Discover which mcs are not in the subset
		mcs.Union(mcs_subset)
	}
}

func BenchmarkMonteCarloRunsInPlace(b *testing.B) {
	// Create a set of numbers from 1 to 1,000
	mcslice := make([]int, 1000)
	// Fill it with numbers from 1 to 1,000
	for i := 0; i < 1000; i++ {
		mcslice[i] = i + 1
	}
	// Create a set from the slice
	mcs := NewSet(mcslice)

	// Create a set that is a subset of `mcs`
	mcs_subset := mcs.Copy()
	mcs_subset.Discard(1)
	mcs_subset.Discard(20)
	mcs_subset.Discard(50)
	mcs_subset.Discard(143)
	mcs_subset.Discard(999)

	// Reset the benchmark timer
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Discover which mcs are not in the subset
		mcs.DifferenceInPlace(mcs_subset)
	}
}

func TestUnion(t *testing.T) {
	testCases := []struct {
		desc string
		in1  Set
		in2  Set
		want Set
	}{
		{
			desc: "entirely overlapping",
			in1:  NewSet([]int{1, 2, 3}),
			in2:  NewSet([]int{1, 2, 3}),
			want: NewSet([]int{1, 2, 3}),
		},
		{
			desc: "some overlap",
			in1:  NewSet([]int{1, 2, 3}),
			in2:  NewSet([]int{2, 3, 4, 5}),
			want: NewSet([]int{1, 2, 3, 4, 5}),
		},
		{
			desc: "no overlap",
			in1:  NewSet([]int{1, 2, 3}),
			in2:  NewSet([]int{4, 5, 6, 7}),
			want: NewSet([]int{1, 2, 3, 4, 5, 6, 7}),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := tC.in1.Union(tC.in2)

			if !got.Equals(tC.want) {
				t.Errorf("got %v; want %v", tC.want, got)
			}
		})
	}
}

func BenchmarkUnionInt(b *testing.B) {
	benchCases := []struct {
		desc string
		in1  Set
		in2  Set
	}{
		{
			desc: "entirely overlapping",
			in1:  NewSet([]int{1, 2, 3}),
			in2:  NewSet([]int{1, 2, 3}),
		},
		{
			desc: "some overlap",
			in1:  NewSet([]int{1, 2, 3}),
			in2:  NewSet([]int{2, 3, 4, 5}),
		},
		{
			desc: "no overlap",
			in1:  NewSet([]int{1, 2, 3}),
			in2:  NewSet([]int{4, 5, 6, 7}),
		},
	}
	for _, bC := range benchCases {
		b.Run(bC.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bC.in1.Union(bC.in2)
			}
		})
	}
}

func FuzzUnion(f *testing.F) {
	// This fuzz test is for checking that Union always matches between the two set types
	f.Add(2)
	f.Add(10)

	f.Fuzz(func(t *testing.T, _n int) {
		n := abs(_n)
		items := make([]int, n)
		// Create n random ints
		for i := 0; i < n; i++ {
			items[i] = rand.Int()
		}

		// Create the sets
		half_way := len(items) / 2
		bitset1 := NewSet(items[:half_way])
		bitset2 := NewSet(items[half_way:])
		set1 := set.NewSet(items[:half_way])
		set2 := set.NewSet(items[half_way:])

		// Take the union
		bitunion := bitset1.Union(bitset2)
		union := set1.Union(set2)

		// Convert them to slices to compare
		bitslice := bitunion.Slice()
		slice := union.Slice()
		slices.Sort(bitslice)
		slices.Sort(slice)

		if !equal(bitslice, slice) {
			t.Errorf("bit set %v did not match set %v", bitslice, slice)
		}
	})
}

func FuzzUnionInPlace(f *testing.F) {
	// This fuzz test is for checking that UnionInPlace always matches between the two
	// set types
	f.Add(2)
	f.Add(10)

	f.Fuzz(func(t *testing.T, _n int) {
		n := abs(_n)
		items := make([]int, n)
		// Create n random ints
		for i := 0; i < n; i++ {
			items[i] = rand.Int()
		}

		// Create the sets
		var split_point int
		if n < 2 {
			split_point = 0
		} else {
			split_point = rand.Intn(len(items))
		}
		bitset1 := NewSet(items[:split_point])
		bitset2 := NewSet(items[split_point:])
		set1 := set.NewSet(items[:split_point])
		set2 := set.NewSet(items[split_point:])

		// Take the union
		bitset1.UnionInPlace(bitset2)
		set1.UnionInPlace(set2)

		// Convert them to slices to compare
		bitslice := bitset1.Slice()
		slice := set1.Slice()
		slices.Sort(bitslice)
		slices.Sort(slice)

		if !equal(bitslice, slice) {
			t.Errorf("bit set %v did not match set %v", bitslice, slice)
		}
	})
}

func TestNumber_to_bitset_representation(t *testing.T) {
	testCases := []struct {
		desc             string
		in               int
		want_is_positive bool
		want_multiplier  uint64
		want_rem_slot    uint64
	}{
		{
			desc:             "0",
			in:               0,
			want_is_positive: true,
			want_multiplier:  0,
			want_rem_slot:    1,
		},
		{
			desc:             "-1",
			in:               -1,
			want_is_positive: false,
			want_multiplier:  0,
			want_rem_slot:    2,
		},
		{
			desc:             "-10",
			in:               -10,
			want_is_positive: false,
			want_multiplier:  0,
			want_rem_slot:    1024,
		},
		{
			desc:             "63",
			in:               63,
			want_is_positive: true,
			want_multiplier:  0,
			want_rem_slot:    9223372036854775808,
		},
		{
			desc:             "64",
			in:               64,
			want_is_positive: true,
			want_multiplier:  1,
			want_rem_slot:    1,
		},
		{
			desc:             "-64",
			in:               -64,
			want_is_positive: false,
			want_multiplier:  1,
			want_rem_slot:    1,
		},
		{
			desc:             "155",
			in:               155,
			want_is_positive: true,
			want_multiplier:  2,
			want_rem_slot:    134217728,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			is_positive, multiplier, rem_slot := number_to_bitset_representation(tC.in)

			if is_positive != tC.want_is_positive {
				t.Errorf("Sign is wrong way: got %v; want %v", is_positive, tC.want_is_positive)
			}

			if multiplier != tC.want_multiplier {
				t.Errorf("Multiplier is incorrect: got %v; want %v", multiplier, tC.want_multiplier)
			}

			if rem_slot != tC.want_rem_slot {
				t.Errorf("Remainder slot is incorrect: got %v; want %v", rem_slot, tC.want_rem_slot)
			}
		})
	}
}

func BenchmarkNumber_to_bitset_representation(b *testing.B) {
	benchCases := []struct {
		desc string
		in   int
	}{
		{
			desc: "0",
			in:   0,
		},
		{
			desc: "-1",
			in:   -1,
		},
		{
			desc: "-10",
			in:   -10,
		},
		{
			desc: "63",
			in:   63,
		},
		{
			desc: "64",
			in:   64,
		},
		{
			desc: "-64",
			in:   -64,
		},
		{
			desc: "155",
			in:   155,
		},
		{
			desc: "4321",
			in:   4321,
		},
		{
			desc: "3849328234",
			in:   3849328234,
		},
	}
	for _, bC := range benchCases {
		b.Run(bC.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				number_to_bitset_representation(bC.in)
			}
		})
	}
}

func TestIntersection(t *testing.T) {

	testCases := []struct {
		desc string
		s1   Set
		s2   Set
		want Set
	}{
		{
			desc: "no intersection",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}),
			want: NewSet([]int{}),
		},
		{
			desc: "some intersection",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}),
			want: NewSet([]int{5, 6, 7, 8, 9, 10}),
		},
		{
			desc: "all intersection",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			want: NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		},
		{
			desc: "a fuzz case",
			s1:   NewSet([]int{6129484611666145821, 4037200794235010051, 5577006791947779410, 8674665223082153551}),
			s2:   NewSet([]int{3916589616287113937, 6334824724549167320, 605394647632969758, 1443635317331776148, 894385949183117216, 2775422040480279449}),
			want: NewSet([]int{}),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if got := tC.s1.Intersection(tC.s2); !got.Equals(tC.want) {
				t.Errorf("got %v, want %v", got, tC.want)
			}
		})
	}
}

func BenchmarkIntersection(b *testing.B) {
	benchCases := []struct {
		desc string
		in1  Set
		in2  Set
	}{
		{
			desc: "entirely overlapping",
			in1:  NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			in2:  NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		},
		{
			desc: "some overlap",
			in1:  NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			in2:  NewSet([]int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}),
		},
		{
			desc: "no overlap",
			in1:  NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			in2:  NewSet([]int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}),
		},
	}
	for _, bC := range benchCases {
		b.Run(bC.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bC.in1.Intersection(bC.in2)
			}
		})
	}
}

func FuzzIntersection(f *testing.F) {
	// This fuzz test is for checking that Intersection always matches between the two
	// set types
	f.Add(2)
	f.Add(10)

	f.Fuzz(func(t *testing.T, _n int) {
		n := abs(_n)
		items := make([]int, n)
		// Create n random ints
		for i := 0; i < n; i++ {
			items[i] = rand.Int()
		}

		// Create the sets
		var split_point int
		if n < 2 {
			split_point = 0
		} else {
			split_point = rand.Intn(len(items))
		}
		bitset1 := NewSet(items[:split_point])
		bitset2 := NewSet(items[split_point:])
		set1 := set.NewSet(items[:split_point])
		set2 := set.NewSet(items[split_point:])

		// Take the intersection
		bitintersection := bitset1.Intersection(bitset2)
		setintersection := set1.Intersection(set2)

		// Convert them to slices to compare
		bitslice := bitintersection.Slice()
		slice := setintersection.Slice()
		slices.Sort(bitslice)
		slices.Sort(slice)

		if !equal(bitslice, slice) {
			t.Errorf("bit set %v did not match set %v\nSet 1 = %v\nSet 2 = %v", bitslice, slice, bitset1, bitset2)
		}
	})
}

func FuzzIntersectionInPlace(f *testing.F) {
	// This fuzz test is for checking that IntersectionInPlace always matches between the two
	// set types
	f.Add(2)
	f.Add(10)

	f.Fuzz(func(t *testing.T, _n int) {
		n := abs(_n)
		items := make([]int, n)
		// Create n random ints
		for i := 0; i < n; i++ {
			items[i] = rand.Int()
		}

		// Create the sets
		var split_point int
		if n < 2 {
			split_point = 0
		} else {
			split_point = rand.Intn(len(items))
		}
		bitset1 := NewSet(items[:split_point])
		bitset2 := NewSet(items[split_point:])
		set1 := set.NewSet(items[:split_point])
		set2 := set.NewSet(items[split_point:])

		// Take the intersection
		bitset1.IntersectionInPlace(bitset2)
		set1.IntersectionInPlace(set2)

		// Convert them to slices to compare
		bitslice := bitset1.Slice()
		slice := set1.Slice()
		slices.Sort(bitslice)
		slices.Sort(slice)

		if !equal(bitslice, slice) {
			t.Errorf("bit set %v did not match set %v\nSet 1 = %v\nSet 2 = %v", bitslice, slice, bitset1, bitset2)
		}
	})
}

func TestIntersectionInPlace(t *testing.T) {

	testCases := []struct {
		desc string
		s1   Set
		s2   Set
		want Set
	}{
		{
			desc: "no intersection",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}),
			want: NewSet([]int{}),
		},
		{
			desc: "some intersection",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}),
			want: NewSet([]int{5, 6, 7, 8, 9, 10}),
		},
		{
			desc: "all intersection",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			want: NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		},
		{
			desc: "a fuzz case",
			s1:   NewSet([]int{6129484611666145821, 4037200794235010051, 5577006791947779410, 8674665223082153551}),
			s2:   NewSet([]int{3916589616287113937, 6334824724549167320, 605394647632969758, 1443635317331776148, 894385949183117216, 2775422040480279449}),
			want: NewSet([]int{}),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			s1 := tC.s1.Copy()
			s2 := tC.s2.Copy()
			s1.IntersectionInPlace(s2)
			if !s1.Equals(tC.want) {
				t.Errorf("got %v, want %v", s1, tC.want)
			}
		})
	}
}

func TestIsDisjoint(t *testing.T) {
	testCases := []struct {
		desc string
		s1   Set
		s2   Set
		want bool
	}{
		{
			desc: "are disjoint",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}),
			want: true,
		},
		{
			desc: "are not disjoint",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}),
			want: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if got := tC.s1.IsDisjoint(tC.s2); got != tC.want {
				t.Errorf("got %v, want %v", got, tC.want)
			}
		})
	}
}

func BenchmarkIsDijointInt(b *testing.B) {
	benchCases := []struct {
		desc string
		in1  Set
		in2  Set
	}{
		{
			desc: "entirely overlapping",
			in1:  NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			in2:  NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		},
		{
			desc: "some overlap",
			in1:  NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			in2:  NewSet([]int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}),
		},
		{
			desc: "no overlap",
			in1:  NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			in2:  NewSet([]int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}),
		},
	}
	for _, bC := range benchCases {
		b.Run(bC.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bC.in1.IsDisjoint(bC.in2)
			}
		})
	}
}

func FuzzIsDisjoint(f *testing.F) {
	// This fuzz test is for checking that IntersectionInPlace always matches between the two
	// set types
	f.Add(2)
	f.Add(10)

	f.Fuzz(func(t *testing.T, _n int) {
		n := abs(_n)
		items := make([]int, n)
		// Create n random ints
		for i := 0; i < n; i++ {
			items[i] = rand.Int()
		}

		// Create the sets
		var split_point int
		if n < 2 {
			split_point = 0
		} else {
			split_point = rand.Intn(len(items))
		}
		bitset1 := NewSet(items[:split_point])
		bitset2 := NewSet(items[split_point:])
		set1 := set.NewSet(items[:split_point])
		set2 := set.NewSet(items[split_point:])

		// Take the intersection
		bitresult := bitset1.IsDisjoint(bitset2)
		setresult := set1.IsDisjoint(set2)

		if bitresult != setresult {
			t.Errorf("got %v, want %v", bitresult, setresult)
		}
	})
}

func TestIsSubset(t *testing.T) {
	testCases := []struct {
		desc string
		s1   Set
		s2   Set
		want bool
	}{
		{
			desc: "exact match",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			want: true,
		},
		{
			desc: "some overlap, but not subset",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}),
			want: false,
		},
		{
			desc: "is small subset",
			s1:   NewSet([]int{1, 5, 8, 9}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			want: true,
		},
		{
			desc: "is not a subset",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			want: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if got := tC.s1.IsSubsetOf(tC.s2); got != tC.want {
				t.Errorf("got %v, want %v", got, tC.want)
			}
		})
	}
}

func BenchmarkIsSubsetInt(b *testing.B) {
	benchCases := []struct {
		desc string
		s1   Set
		s2   Set
	}{
		{
			desc: "exact match",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		},
		{
			desc: "some overlap, but not subset",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}),
		},
		{
			desc: "is small subset",
			s1:   NewSet([]int{1, 5, 8, 9}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		},
		{
			desc: "is not a subset",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		},
	}
	for _, bC := range benchCases {
		b.Run(bC.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bC.s1.IsSubsetOf(bC.s2)
			}
		})
	}
}

func FuzzIsSubsetOf(f *testing.F) {
	// This fuzz test is for checking that IsSubsetOf always matches between the two
	// set types
	f.Add(2)
	f.Add(10)

	f.Fuzz(func(t *testing.T, _n int) {
		n := abs(_n)
		items := make([]int, n)
		// Create n random ints
		for i := 0; i < n; i++ {
			items[i] = rand.Int()
		}

		// Create the sets
		var split_point int
		if n < 2 {
			split_point = 0
		} else {
			split_point = rand.Intn(len(items))
		}
		bitset1 := NewSet(items[:split_point])
		bitset2 := NewSet(items[split_point:])
		set1 := set.NewSet(items[:split_point])
		set2 := set.NewSet(items[split_point:])

		// Take the intersection
		bitresult := bitset1.IsSubsetOf(bitset2)
		setresult := set1.IsSubsetOf(set2)

		if bitresult != setresult {
			t.Errorf("got %v, want %v", bitresult, setresult)
		}
	})
}

func TestIsProperSubset(t *testing.T) {
	testCases := []struct {
		desc string
		s1   Set
		s2   Set
		want bool
	}{
		{
			desc: "exact match",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			want: false,
		},
		{
			desc: "some overlap, but not subset",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}),
			want: false,
		},
		{
			desc: "is small subset",
			s1:   NewSet([]int{1, 5, 8, 9}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			want: true,
		},
		{
			desc: "is not a subset",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			want: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if got := tC.s1.IsProperSubsetOf(tC.s2); got != tC.want {
				t.Errorf("got %v, want %v", got, tC.want)
			}
		})
	}
}

func BenchmarkIsProperSubsetInt(b *testing.B) {
	benchCases := []struct {
		desc string
		s1   Set
		s2   Set
	}{
		{
			desc: "exact match",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		},
		{
			desc: "some overlap, but not subset",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}),
		},
		{
			desc: "is small subset",
			s1:   NewSet([]int{1, 5, 8, 9}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		},
		{
			desc: "is not a subset",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		},
	}
	for _, bC := range benchCases {
		b.Run(bC.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bC.s1.IsProperSubsetOf(bC.s2)
			}
		})
	}
}

func FuzzIsProperSubsetOf(f *testing.F) {
	// This fuzz test is for checking that IsProperSubsetOf always matches between the two
	// set types
	f.Add(2)
	f.Add(10)

	f.Fuzz(func(t *testing.T, _n int) {
		n := abs(_n)
		items := make([]int, n)
		// Create n random ints
		for i := 0; i < n; i++ {
			items[i] = rand.Int()
		}

		// Create the sets
		var split_point int
		if n < 2 {
			split_point = 0
		} else {
			split_point = rand.Intn(len(items))
		}
		bitset1 := NewSet(items[:split_point])
		bitset2 := NewSet(items[split_point:])
		set1 := set.NewSet(items[:split_point])
		set2 := set.NewSet(items[split_point:])

		// Take the intersection
		bitresult := bitset1.IsProperSubsetOf(bitset2)
		setresult := set1.IsProperSubsetOf(set2)

		if bitresult != setresult {
			t.Errorf("got %v, want %v", bitresult, setresult)
		}
	})
}

func TestIsSuperSetOf(t *testing.T) {
	testCases := []struct {
		desc string
		s1   Set
		s2   Set
		want bool
	}{
		{
			desc: "exact match",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			want: true,
		},
		{
			desc: "some overlap, but not subset",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}),
			want: false,
		},
		{
			desc: "is small subset",
			s1:   NewSet([]int{1, 5, 8, 9}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			want: false,
		},
		{
			desc: "is not a subset",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			want: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if got := tC.s1.IsSuperSetOf(tC.s2); got != tC.want {
				t.Errorf("got %v, want %v", got, tC.want)
			}
		})
	}
}

func BenchmarkIsSuperSetOf(b *testing.B) {
	benchCases := []struct {
		desc string
		s1   Set
		s2   Set
	}{
		{
			desc: "exact match",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		},
		{
			desc: "some overlap, but not subset",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}),
		},
		{
			desc: "is small subset",
			s1:   NewSet([]int{1, 5, 8, 9}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		},
		{
			desc: "is not a subset",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		},
	}
	for _, bC := range benchCases {
		b.Run(bC.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bC.s1.IsSuperSetOf(bC.s2)
			}
		})
	}
}

func FuzzIsSuperSetOf(f *testing.F) {
	// This fuzz test is for checking that IsSuperSetOf always matches between the two
	// set types
	f.Add(2)
	f.Add(10)

	f.Fuzz(func(t *testing.T, _n int) {
		n := abs(_n)
		items := make([]int, n)
		// Create n random ints
		for i := 0; i < n; i++ {
			items[i] = rand.Int()
		}

		// Create the sets
		var split_point int
		if n < 2 {
			split_point = 0
		} else {
			split_point = rand.Intn(len(items))
		}
		bitset1 := NewSet(items[:split_point])
		bitset2 := NewSet(items[split_point:])
		set1 := set.NewSet(items[:split_point])
		set2 := set.NewSet(items[split_point:])

		// Take the intersection
		bitresult := bitset1.IsSuperSetOf(bitset2)
		setresult := set1.IsSuperSetOf(set2)

		if bitresult != setresult {
			t.Errorf("got %v, want %v", bitresult, setresult)
		}
	})
}

func TestIsProperSuperSetOf(t *testing.T) {
	testCases := []struct {
		desc string
		s1   Set
		s2   Set
		want bool
	}{
		{
			desc: "exact match",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			want: false,
		},
		{
			desc: "some overlap, but not subset",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}),
			want: false,
		},
		{
			desc: "is small subset",
			s1:   NewSet([]int{1, 5, 8, 9}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			want: false,
		},
		{
			desc: "is not a subset",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			want: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if got := tC.s1.IsProperSuperSetOf(tC.s2); got != tC.want {
				t.Errorf("got %v, want %v", got, tC.want)
			}
		})
	}
}

func BenchmarkIsProperSuperSetOf(b *testing.B) {
	benchCases := []struct {
		desc string
		s1   Set
		s2   Set
	}{
		{
			desc: "exact match",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		},
		{
			desc: "some overlap, but not subset",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}),
		},
		{
			desc: "is small subset",
			s1:   NewSet([]int{1, 5, 8, 9}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		},
		{
			desc: "is not a subset",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		},
	}
	for _, bC := range benchCases {
		b.Run(bC.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bC.s1.IsProperSuperSetOf(bC.s2)
			}
		})
	}
}

func FuzzIsProperSuperSetOf(f *testing.F) {
	// This fuzz test is for checking that IsProperSuperSetOf always matches between the two
	// set types
	f.Add(2)
	f.Add(10)

	f.Fuzz(func(t *testing.T, _n int) {
		n := abs(_n)
		items := make([]int, n)
		// Create n random ints
		for i := 0; i < n; i++ {
			items[i] = rand.Int()
		}

		// Create the sets
		var split_point int
		if n < 2 {
			split_point = 0
		} else {
			split_point = rand.Intn(len(items))
		}
		bitset1 := NewSet(items[:split_point])
		bitset2 := NewSet(items[split_point:])
		set1 := set.NewSet(items[:split_point])
		set2 := set.NewSet(items[split_point:])

		// Take the intersection
		bitresult := bitset1.IsProperSuperSetOf(bitset2)
		setresult := set1.IsProperSuperSetOf(set2)

		if bitresult != setresult {
			t.Errorf("got %v, want %v", bitresult, setresult)
		}
	})
}

func TestDifference(t *testing.T) {
	testCases := []struct {
		desc string
		s1   Set
		s2   Set
		want Set
	}{
		{
			desc: "exact match",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			want: NewSet([]int{}),
		},
		{
			desc: "some overlap",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}),
			want: NewSet([]int{1, 2, 3, 4}),
		},
		{
			desc: "tiny overlap",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{1, 2}),
			want: NewSet([]int{3, 4, 5, 6, 7, 8, 9, 10}),
		},
		{
			desc: "no overlap",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}),
			want: NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := tC.s1.Difference(tC.s2)
			if !got.Equals(tC.want) {
				t.Errorf("got %v, want %v", got, tC.want)
			}
		})
	}
}

func BenchmarkDifference(b *testing.B) {
	benchCases := []struct {
		desc string
		s1   Set
		s2   Set
	}{
		{
			desc: "exact match",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		},
		{
			desc: "some overlap",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}),
		},
		{
			desc: "tiny overlap",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{1, 2}),
		},
		{
			desc: "no overlap",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}),
		},
	}
	for _, bC := range benchCases {
		b.Run(bC.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bC.s1.Difference(bC.s2)
			}
		})
	}
}

func FuzzDifference(f *testing.F) {
	// This fuzz test is for checking that Difference always matches between the two
	// set types
	f.Add(2)
	f.Add(10)

	f.Fuzz(func(t *testing.T, _n int) {
		n := abs(_n)
		items := make([]int, n)
		// Create n random ints
		for i := 0; i < n; i++ {
			items[i] = rand.Int()
		}

		// Create the sets
		var split_point int
		if n < 2 {
			split_point = 0
		} else {
			split_point = rand.Intn(len(items))
		}
		bitset1 := NewSet(items[:split_point])
		bitset2 := NewSet(items[split_point:])
		set1 := set.NewSet(items[:split_point])
		set2 := set.NewSet(items[split_point:])

		// Take the Difference
		bitset := bitset1.Difference(bitset2)
		set := set1.Difference(set2)

		// Check that both sets are the same
		setslice := set.Slice()
		bitsetslice := bitset.Slice()

		slices.Sort(setslice)
		slices.Sort(bitsetslice)

		if !equal(setslice, bitsetslice) {
			t.Errorf("got %v, want %v", bitsetslice, setslice)
		}
	})
}

func TestDifferenceInPlace(t *testing.T) {
	testCases := []struct {
		desc string
		s1   Set
		s2   Set
		want Set
	}{
		{
			desc: "exact match",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			want: NewSet([]int{}),
		},
		{
			desc: "some overlap",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}),
			want: NewSet([]int{1, 2, 3, 4}),
		},
		{
			desc: "tiny overlap",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{1, 2}),
			want: NewSet([]int{3, 4, 5, 6, 7, 8, 9, 10}),
		},
		{
			desc: "no overlap",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}),
			want: NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if tC.s1.DifferenceInPlace(tC.s2); !tC.s1.Equals(tC.want) {
				t.Errorf("got %v, want %v", tC.s1, tC.want)
			}
		})
	}
}

func FuzzDifferenceInPlace(f *testing.F) {
	// This fuzz test is for checking that DifferenceInPlace always matches between the two
	// set types
	f.Add(2)
	f.Add(10)

	f.Fuzz(func(t *testing.T, _n int) {
		n := abs(_n)
		items := make([]int, n)
		// Create n random ints
		for i := 0; i < n; i++ {
			items[i] = rand.Int()
		}

		// Create the sets
		var split_point int
		if n < 2 {
			split_point = 0
		} else {
			split_point = rand.Intn(len(items))
		}
		bitset1 := NewSet(items[:split_point])
		bitset2 := NewSet(items[split_point:])
		set1 := set.NewSet(items[:split_point])
		set2 := set.NewSet(items[split_point:])

		// Take the Difference
		bitset1.DifferenceInPlace(bitset2)
		set1.DifferenceInPlace(set2)

		// Check that both sets are the same
		setslice := set1.Slice()
		bitsetslice := bitset1.Slice()

		slices.Sort(setslice)
		slices.Sort(bitsetslice)

		if !equal(setslice, bitsetslice) {
			t.Errorf("got %v, want %v", bitsetslice, setslice)
		}
	})
}

func TestSymmetricDifference(t *testing.T) {
	testCases := []struct {
		desc string
		s1   Set
		s2   Set
		want Set
	}{
		{
			desc: "exact match",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			want: NewSet([]int{}),
		},
		{
			desc: "some overlap",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}),
			want: NewSet([]int{1, 2, 3, 4, 11, 12, 13, 14}),
		},
		{
			desc: "tiny overlap",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{1, 2}),
			want: NewSet([]int{3, 4, 5, 6, 7, 8, 9, 10, 3, 4}),
		},
		{
			desc: "no overlap",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{-11, -12, -13, -14, -15, -16, -17, -18, -19, -20}),
			want: NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, -11, -12, -13, -14, -15, -16, -17, -18, -19, -20}),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if got := tC.s1.SymmetricDifference(tC.s2); !got.Equals(tC.want) {
				t.Errorf("got %v, want %v", got, tC.want)
			}
		})
	}
}

func BenchmarkSymmetricDifference(b *testing.B) {
	benchCases := []struct {
		desc string
		s1   Set
		s2   Set
	}{
		{
			desc: "exact match",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		},
		{
			desc: "some overlap",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}),
		},
		{
			desc: "tiny overlap",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{1, 2}),
		},
		{
			desc: "no overlap",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}),
		},
	}
	for _, bC := range benchCases {
		b.Run(bC.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bC.s1.SymmetricDifference(bC.s2)
			}
		})
	}
}

func FuzzSymmetricDifference(f *testing.F) {
	// This fuzz test is for checking that SymmetricDifference always matches between the two
	// set types
	f.Add(2)
	f.Add(10)

	f.Fuzz(func(t *testing.T, _n int) {
		n := abs(_n)
		items := make([]int, n)
		// Create n random ints
		for i := 0; i < n; i++ {
			items[i] = rand.Int()
		}

		// Create the sets
		var split_point int
		if n < 2 {
			split_point = 0
		} else {
			split_point = rand.Intn(len(items))
		}
		bitset1 := NewSet(items[:split_point])
		bitset2 := NewSet(items[split_point:])
		set1 := set.NewSet(items[:split_point])
		set2 := set.NewSet(items[split_point:])

		// Take the SymmetricDifference
		bitset := bitset1.SymmetricDifference(bitset2)
		set := set1.SymmetricDifference(set2)

		// Check that both sets are the same
		setslice := set.Slice()
		bitsetslice := bitset.Slice()

		slices.Sort(setslice)
		slices.Sort(bitsetslice)

		if !equal(setslice, bitsetslice) {
			t.Errorf("got %v, want %v", bitsetslice, setslice)
		}
	})
}

func TestSymmetricDifferenceInPlace(t *testing.T) {
	testCases := []struct {
		desc string
		s1   Set
		s2   Set
		want Set
	}{
		{
			desc: "exact match",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			want: NewSet([]int{}),
		},
		{
			desc: "some overlap",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}),
			want: NewSet([]int{1, 2, 3, 4, 11, 12, 13, 14}),
		},
		{
			desc: "tiny overlap",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{1, 2}),
			want: NewSet([]int{3, 4, 5, 6, 7, 8, 9, 10, 3, 4}),
		},
		{
			desc: "no overlap",
			s1:   NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			s2:   NewSet([]int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}),
			want: NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if tC.s1.SymmetricDifferenceInPlace(tC.s2); !tC.s1.Equals(tC.want) {
				t.Errorf("got %v, want %v", tC.s1, tC.want)
			}
		})
	}
}

func FuzzSymmetricDifferenceInPlace(f *testing.F) {
	// This fuzz test is for checking that SymmetricDifferenceInPlace always matches between the two
	// set types
	f.Add(2)
	f.Add(10)

	f.Fuzz(func(t *testing.T, _n int) {
		n := abs(_n)
		items := make([]int, n)
		// Create n random ints
		for i := 0; i < n; i++ {
			items[i] = rand.Int()
		}

		// Create the sets
		var split_point int
		if n < 2 {
			split_point = 0
		} else {
			split_point = rand.Intn(len(items))
		}
		bitset1 := NewSet(items[:split_point])
		bitset2 := NewSet(items[split_point:])
		set1 := set.NewSet(items[:split_point])
		set2 := set.NewSet(items[split_point:])

		// Take the SymmetricDifferenceInPlace
		bitset1.SymmetricDifferenceInPlace(bitset2)
		set1.SymmetricDifferenceInPlace(set2)

		// Check that both sets are the same
		setslice := set1.Slice()
		bitsetslice := bitset1.Slice()

		slices.Sort(setslice)
		slices.Sort(bitsetslice)

		if !equal(setslice, bitsetslice) {
			t.Errorf("got %v, want %v", bitsetslice, setslice)
		}
	})
}
