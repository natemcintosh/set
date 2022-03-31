package bitset

import (
	"fmt"
	"testing"

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
