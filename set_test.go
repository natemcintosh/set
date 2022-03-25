package set

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewStringSet(t *testing.T) {
	in := []string{"a", "b", "c"}
	want := Set[string]{data: map[string]struct{}{"a": {}, "b": {}, "c": {}}}
	got := NewSet(in)

	if !reflect.DeepEqual(want, got) {
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

	if !reflect.DeepEqual(want, got) {
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
	if !reflect.DeepEqual(s1, want1) {
		t.Errorf("got %v; want %v", s1, want1)
	}

	s1.Add(4)
	want1 = NewSet([]int{1, 2, 3, 4})
	if !reflect.DeepEqual(s1, want1) {
		t.Errorf("got %v; want %v", s1, want1)
	}

	s2 := NewSet([]string{"a", "b", "c", "longer string"})
	s2.Add("a")
	want2 := NewSet([]string{"a", "b", "c", "longer string"})
	if !reflect.DeepEqual(s2, want2) {
		t.Errorf("got %v; want %v", s2, want2)
	}

	s2.Add("d")
	want2 = NewSet([]string{"a", "b", "c", "longer string", "d"})
	if !reflect.DeepEqual(s2, want2) {
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

			if !reflect.DeepEqual(tC.want, got) {
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

// func TestRemove(t *testing.T) {
// 	s := NewSet([]float64{1.0, 2.0, 4.5, 3.7})

// }
