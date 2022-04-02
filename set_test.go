package set

import (
	"fmt"
	"strings"
	"testing"
)

func TestNewStringSet(t *testing.T) {
	in := []string{"a", "b", "c"}
	want := Set[string]{data: map[string]struct{}{"a": {}, "b": {}, "c": {}}}
	got := NewSet(in)

	if !want.Equals(got) {
		t.Errorf("got %v; want %v", want, got)
	}
}

func BenchmarkNewStringSet(b *testing.B) {
	in := []string{"a", "b", "c", "longer string"}
	for i := 0; i < b.N; i++ {
		NewSet(in)
	}
}

func TestNewFloatSet(t *testing.T) {
	in := []float64{1.0, 2.0, 3.0}
	want := Set[float64]{data: map[float64]struct{}{1.0: {}, 2.0: {}, 3.0: {}}}
	got := NewSet(in)

	if !want.Equals(got) {
		t.Errorf("got %v; want %v", want, got)
	}
}

func BenchmarkNewFloatSet(b *testing.B) {
	in := []float64{1.0, 2.0, 3.0}
	for i := 0; i < b.N; i++ {
		NewSet(in)
	}
}

func TestAdd(t *testing.T) {
	s1 := NewSet([]int{1, 2, 3})
	s1.Add(3)
	want1 := NewSet([]int{1, 2, 3})
	if !s1.Equals(want1) {
		t.Errorf("got %v; want %v", s1, want1)
	}

	s1.Add(4)
	want1 = NewSet([]int{1, 2, 3, 4})
	if !s1.Equals(want1) {
		t.Errorf("got %v; want %v", s1, want1)
	}

	s2 := NewSet([]string{"a", "b", "c", "longer string"})
	s2.Add("a")
	want2 := NewSet([]string{"a", "b", "c", "longer string"})
	if !s2.Equals(want2) {
		t.Errorf("got %v; want %v", s2, want2)
	}

	s2.Add("d")
	want2 = NewSet([]string{"a", "b", "c", "longer string", "d"})
	if !s2.Equals(want2) {
		t.Errorf("got %v; want %v", s2, want2)
	}
}

func TestUnionInt(t *testing.T) {
	testCases := []struct {
		desc string
		in1  Set[int]
		in2  Set[int]
		want Set[int]
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

			if !tC.want.Equals(got) {
				t.Errorf("got %v; want %v", tC.want, got)
			}
		})
	}
}

func BenchmarkUnionInt(b *testing.B) {
	benchCases := []struct {
		desc string
		in1  Set[int]
		in2  Set[int]
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

func BenchmarkUnionString(b *testing.B) {
	benchCases := []struct {
		desc string
		in1  Set[string]
		in2  Set[string]
	}{
		{
			desc: "entirely overlapping",
			in1:  NewSet(strings.Fields("hello world, what is your name")),
			in2:  NewSet(strings.Fields("hello world, what is your name")),
		},
		{
			desc: "some overlap",
			in1:  NewSet(strings.Fields("hello world, what is your name")),
			in2:  NewSet(strings.Fields("hello world, here is something else entirely")),
		},
		{
			desc: "no overlap",
			in1:  NewSet(strings.Fields("hello world, what is your name")),
			in2:  NewSet(strings.Fields("something else entirely here")),
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

func BenchmarkUnionInPlaceInt(b *testing.B) {
	benchCases := []struct {
		desc string
		in1  Set[int]
		in2  Set[int]
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
				bC.in1.UnionInPlace(bC.in2)
			}
		})
	}
}

func BenchmarkUnionInPlaceString(b *testing.B) {
	benchCases := []struct {
		desc string
		in1  Set[string]
		in2  Set[string]
	}{
		{
			desc: "entirely overlapping",
			in1:  NewSet(strings.Fields("hello world, what is your name")),
			in2:  NewSet(strings.Fields("hello world, what is your name")),
		},
		{
			desc: "some overlap",
			in1:  NewSet(strings.Fields("hello world, what is your name")),
			in2:  NewSet(strings.Fields("hello world, here is something else entirely")),
		},
		{
			desc: "no overlap",
			in1:  NewSet(strings.Fields("hello world, what is your name")),
			in2:  NewSet(strings.Fields("something else entirely here")),
		},
	}
	for _, bC := range benchCases {
		b.Run(bC.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bC.in1.UnionInPlace(bC.in2)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	testCases := []struct {
		desc           string
		s              Set[float64]
		v              float64
		want_set       Set[float64]
		want_err_value error
	}{
		{
			desc:           "valid remove",
			s:              NewSet([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			v:              5,
			want_set:       NewSet([]float64{1, 2, 3, 4, 6, 7, 8, 9, 10}),
			want_err_value: nil,
		},
		{
			desc:           "invalid remove",
			s:              NewSet([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			v:              11,
			want_set:       NewSet([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			want_err_value: ErrElementNotFound,
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

func TestDiscard(t *testing.T) {
	testCases := []struct {
		desc     string
		s        Set[string]
		v        string
		want_set Set[string]
	}{
		{
			desc:     "value in set",
			s:        NewSet([]string{"hello", "world", "what", "is", "up", "dude"}),
			v:        "hello",
			want_set: NewSet([]string{"world", "what", "is", "up", "dude"}),
		},
		{
			desc:     "value not in set",
			s:        NewSet([]string{"hello", "world", "what", "is", "up", "dude"}),
			v:        "not",
			want_set: NewSet([]string{"hello", "world", "what", "is", "up", "dude"}),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tC.s.Discard(tC.v)
			if !tC.s.Equals(tC.want_set) {
				t.Errorf("got %v, want %v", tC.s, tC.want_set)
			}
		})
	}
}

func TestPop(t *testing.T) {
	testCases := []struct {
		desc     string
		s        Set[int]
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

func TestContains(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	testCases := []struct {
		desc string
		s    Set[Person]
		v    Person
		want bool
	}{
		{
			desc: "valid contains",
			s:    NewSet([]Person{{"Bob", 42}, {"Alice", 24}, {"Charlie", 12}}),
			v:    Person{"Bob", 42},
			want: true,
		},
		{
			desc: "invalid contains, partial match",
			s:    NewSet([]Person{{"Bob", 42}, {"Alice", 24}, {"Charlie", 12}}),
			v:    Person{"Bob", 43},
			want: false,
		},
		{
			desc: "invalid contains, no match at all",
			s:    NewSet([]Person{{"Bob", 42}, {"Alice", 24}, {"Charlie", 12}}),
			v:    Person{"Nate", 43},
			want: false,
		},
		{
			desc: "invalid contains (empty set)",
			s:    NewSet([]Person{}),
			v:    Person{"Bob", 42},
			want: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if got := tC.s.Contains(tC.v); got != tC.want {
				t.Errorf("got %v, want %v", got, tC.want)
			}
		})
	}
}

func BenchmarkContains(b *testing.B) {
	type Person struct {
		Name string
		Age  int
	}

	benchCases := []struct {
		desc string
		s    Set[Person]
		v    Person
		want bool
	}{
		{
			desc: "valid contains",
			s:    NewSet([]Person{{"Bob", 42}, {"Alice", 24}, {"Charlie", 12}}),
			v:    Person{"Bob", 42},
			want: true,
		},
		{
			desc: "invalid contains, partial match",
			s:    NewSet([]Person{{"Bob", 42}, {"Alice", 24}, {"Charlie", 12}}),
			v:    Person{"Bob", 43},
			want: false,
		},
		{
			desc: "invalid contains, no match at all",
			s:    NewSet([]Person{{"Bob", 42}, {"Alice", 24}, {"Charlie", 12}}),
			v:    Person{"Nate", 43},
			want: false,
		},
		{
			desc: "invalid contains (empty set)",
			s:    NewSet([]Person{}),
			v:    Person{"Bob", 42},
			want: false,
		},
	}
	for _, bC := range benchCases {
		b.Run(bC.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bC.s.Contains(bC.v)
			}
		})
	}
}

func TestIntersection(t *testing.T) {

	testCases := []struct {
		desc string
		s1   Set[int]
		s2   Set[int]
		want Set[int]
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
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if got := tC.s1.Intersection(tC.s2); !got.Equals(tC.want) {
				t.Errorf("got %v, want %v", got, tC.want)
			}
		})
	}
}

func BenchmarkIntersectionString(b *testing.B) {
	benchCases := []struct {
		desc string
		in1  Set[string]
		in2  Set[string]
	}{
		{
			desc: "entirely overlapping",
			in1:  NewSet(strings.Fields("hello world, what is your name")),
			in2:  NewSet(strings.Fields("hello world, what is your name")),
		},
		{
			desc: "some overlap",
			in1:  NewSet(strings.Fields("hello world, what is your name")),
			in2:  NewSet(strings.Fields("hello world, here is something else entirely")),
		},
		{
			desc: "no overlap",
			in1:  NewSet(strings.Fields("hello world, what is your name")),
			in2:  NewSet(strings.Fields("something else entirely here")),
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

func BenchmarkIntersectionInt(b *testing.B) {
	benchCases := []struct {
		desc string
		in1  Set[int]
		in2  Set[int]
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

func TestIsDisjoint(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	type Dog struct {
		Name  string
		Age   int
		Owner Person
	}

	testCases := []struct {
		desc string
		s1   Set[Dog]
		s2   Set[Dog]
		want bool
	}{
		{
			desc: "are disjoint",
			s1:   NewSet([]Dog{{"Fido", 3, Person{"Bob", 42}}, {"Rover", 4, Person{"Alice", 24}}}),
			s2:   NewSet([]Dog{{"Spot", 5, Person{"Bob", 42}}, {"Snoopy", 6, Person{"Bob", 42}}}),
			want: true,
		},
		{
			desc: "are not disjoint",
			s1:   NewSet([]Dog{{"Fido", 3, Person{"Bob", 42}}, {"Rover", 4, Person{"Alice", 24}}}),
			s2:   NewSet([]Dog{{"Fido", 3, Person{"Bob", 42}}, {"Snoopy", 6, Person{"Bob", 42}}}),
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
		in1  Set[int]
		in2  Set[int]
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

func BenchmarkIsDijointString(b *testing.B) {
	benchCases := []struct {
		desc string
		in1  Set[string]
		in2  Set[string]
	}{
		{
			desc: "entirely overlapping",
			in1:  NewSet(strings.Fields("hello world, what is your name")),
			in2:  NewSet(strings.Fields("hello world, what is your name")),
		},
		{
			desc: "some overlap",
			in1:  NewSet(strings.Fields("hello world, what is your name")),
			in2:  NewSet(strings.Fields("hello world, something else here")),
		},
		{
			desc: "no overlap",
			in1:  NewSet(strings.Fields("hello world, what is your name")),
			in2:  NewSet(strings.Fields("something else entirely here")),
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

func TestIsSubset(t *testing.T) {
	testCases := []struct {
		desc string
		s1   Set[int]
		s2   Set[int]
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
		s1   Set[int]
		s2   Set[int]
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

func TestIsProperSubset(t *testing.T) {
	testCases := []struct {
		desc string
		s1   Set[int]
		s2   Set[int]
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
		s1   Set[int]
		s2   Set[int]
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

func TestIsSuperSetOf(t *testing.T) {
	testCases := []struct {
		desc string
		s1   Set[int]
		s2   Set[int]
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
		s1   Set[int]
		s2   Set[int]
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

func TestIsProperSuperSetOf(t *testing.T) {
	testCases := []struct {
		desc string
		s1   Set[int]
		s2   Set[int]
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
		s1   Set[int]
		s2   Set[int]
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

func TestDifference(t *testing.T) {
	testCases := []struct {
		desc string
		s1   Set[int]
		s2   Set[int]
		want Set[int]
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
			if got := tC.s1.Difference(tC.s2); !got.Equals(tC.want) {
				t.Errorf("got %v, want %v", got, tC.want)
			}
		})
	}
}

func BenchmarkDifference(b *testing.B) {
	benchCases := []struct {
		desc string
		s1   Set[int]
		s2   Set[int]
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

func TestDifferenceInPlace(t *testing.T) {
	testCases := []struct {
		desc string
		s1   Set[int]
		s2   Set[int]
		want Set[int]
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

func TestSymmetricDifference(t *testing.T) {
	testCases := []struct {
		desc string
		s1   Set[int]
		s2   Set[int]
		want Set[int]
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
			if got := tC.s1.SymmetricDifference(tC.s2); !got.Equals(tC.want) {
				t.Errorf("got %v, want %v", got, tC.want)
			}
		})
	}
}

func BenchmarkSymmetricDifference(b *testing.B) {
	benchCases := []struct {
		desc string
		s1   Set[int]
		s2   Set[int]
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

func TestSymmetricDifferenceInPlace(t *testing.T) {
	testCases := []struct {
		desc string
		s1   Set[int]
		s2   Set[int]
		want Set[int]
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
		mcs.Difference(mcs_subset)
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

func TestString(t *testing.T) {
	s := NewSet([]int{1, 2, 3, 4})
	str_version := s.String()

	if !strings.HasPrefix(str_version, "{") {
		t.Errorf("%v doesn't start with a '{'", str_version)
	}

	if !strings.HasSuffix(str_version, "}") {
		t.Errorf("%v doesn't end with a '}'", str_version)
	}

	// This would obviously fail if testing a set of strings which contain commas
	counted_commas := strings.Count(str_version, ", ")
	expected_commas := s.Len() - 1
	if counted_commas != expected_commas {
		t.Errorf("saw %d ', '; wanted %d", counted_commas, expected_commas)
	}
}

func TestFormat(t *testing.T) {
	s := NewSet([]int{1, 2, 3, 4})
	str_version := fmt.Sprintf("%v", s)

	if !strings.HasPrefix(str_version, "{") {
		t.Errorf("%v doesn't start with a '{'", str_version)
	}

	if !strings.HasSuffix(str_version, "}") {
		t.Errorf("%v doesn't end with a '}'", str_version)
	}

	// This would obviously fail if testing a set of strings which contain commas
	counted_commas := strings.Count(str_version, ", ")
	expected_commas := s.Len() - 1
	if counted_commas != expected_commas {
		t.Errorf("saw %d ', '; wanted %d", counted_commas, expected_commas)
	}
}
