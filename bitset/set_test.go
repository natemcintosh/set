package bitset

import (
	"fmt"
	"testing"

	"github.com/natemcintosh/set"
	"golang.org/x/exp/slices"
)

func TestString(t *testing.T) {
	testCases := []struct {
		desc string
		set  []int
		want string
	}{
		{
			desc: "empty",
			set:  []int{},
			want: "{}",
		},
		{
			desc: "one",
			set:  []int{1},
			want: "{1}",
		},
		{
			desc: "a few",
			set:  []int{-1, 3, 10},
			want: "{-1, 3, 10}",
		},
		{
			desc: "comes out sorted",
			set:  []int{1, -1},
			want: "{-1, 1}",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// Create the set
			set := NewSet(tC.set)

			// Create the string
			s := fmt.Sprintf("%v", set)

			// Compare
			if s != tC.want {
				t.Errorf("got %v, want %v", s, tC.want)
			}
		})
	}
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
		slices.Sort(want)
		want = slices.Compact(want)

		// Create a bitset
		newset := NewSet(want)

		// Get the numbers back out
		got := newset.Slice()

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
	// We are hoping to find out of bounds panics with this fuzz test
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
		// Create the set
		set := NewSet([]int{s1, s2, s3, s4, s5, s6, s7, s8, s9})

		// Run Contains, hoping to find panics
		set.Contains(s10)

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
	// We are hoping to find out of bounds panics with this fuzz test
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
		// Create the set
		set := NewSet([]int{s1, s2, s3, s4, s5, s6, s7, s8, s9})

		// Run Add, hoping to find panics
		set.Add(s10)

	})
}

func TestRemove(t *testing.T) {
	testCases := []struct {
		desc               string
		s                  Set
		v                  int
		want_set           Set
		want_err_value     error
		want_smallest_item int
	}{
		{
			desc:               "valid remove",
			s:                  NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			v:                  5,
			want_set:           NewSet([]int{1, 2, 3, 4, 6, 7, 8, 9, 10}),
			want_err_value:     nil,
			want_smallest_item: 1,
		},
		{
			desc:               "invalid remove",
			s:                  NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			v:                  11,
			want_set:           NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			want_err_value:     ErrElementNotFound,
			want_smallest_item: 1,
		},
		{
			desc:               "smallest item",
			s:                  NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			v:                  1,
			want_set:           NewSet([]int{2, 3, 4, 5, 6, 7, 8, 9, 10}),
			want_err_value:     nil,
			want_smallest_item: 2,
		},
		{
			desc:               "largest item",
			s:                  NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			v:                  10,
			want_set:           NewSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			want_err_value:     nil,
			want_smallest_item: 1,
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
			if tC.want_smallest_item != tC.s.smallest_item {
				t.Errorf("smallest item is %v; should be %v", tC.s.smallest_item, tC.want_smallest_item)
			}
		})
	}
}

func FuzzRemoveDiscard(f *testing.F) {
	// We are hoping to find out of bounds panics with this fuzz test
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
		// Create the set
		set := NewSet([]int{s1, s2, s3, s4, s5, s6, s7, s8, s9})

		prev_smallest := set.smallest_item
		prev_largest := set.get_upper_value()

		// Run Add, hoping to find panics
		set.Remove(s10)

		// If s10 == the smallest item in the set, then check that the smallest item
		// is no longer equal to that
		if (s10 == prev_smallest) && (s10 == set.smallest_item) {
			t.Errorf("The smallest item was removed, but the internal field has not been updated")
		}

		// If s10 == the largest item in the set, then check that the largest item is
		// no longer equal to that
		if (s10 == prev_largest) && (s10 == set.get_upper_value()) {
			t.Errorf("The largest item was removed, but the slice hasn't been updated")
		}

		// Recreate set and redo with Discard
		set = NewSet([]int{s1, s2, s3, s4, s5, s6, s7, s8, s9})

		prev_smallest = set.smallest_item
		prev_largest = set.get_upper_value()

		// Run Add, hoping to find panics
		set.Remove(s10)

		// If s10 == the smallest item in the set, then check that the smallest item
		// is no longer equal to that
		if (s10 == prev_smallest) && (s10 == set.smallest_item) {
			t.Errorf("The smallest item was discarded, but the internal field has not been updated")
		}

		// If s10 == the largest item in the set, then check that the largest item is
		// no longer equal to that
		if (s10 == prev_largest) && (s10 == set.get_upper_value()) {
			t.Errorf("The largest item was discarded, but the slice hasn't been updated")
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
	// We are hoping to find out of bounds panics with this fuzz test
	// We are also checking that both set and bitset produce the same result
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
		// Create the sets
		bitset1 := NewSet([]int{s1, s2, s3, s4, s5})
		bitset2 := NewSet([]int{s6, s7, s8, s9, s10})

		set1 := set.NewSet([]int{s1, s2, s3, s4, s5})
		set2 := set.NewSet([]int{s6, s7, s8, s9, s10})

		// Take the union
		bitunion := bitset1.Union(bitset2)
		union := set1.Union(set2)

		// Convert them to slices to compare
		bitslice := bitunion.Slice()
		slice := union.Slice()
		slices.Sort(slice)

		if !equal(bitslice, slice) {
			t.Errorf("bit set %v did not match set %v", bitslice, slice)
		}
	})
}
