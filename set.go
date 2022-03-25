// set is a library designed to help you do set operations on comparable data types
package set

type Set[T comparable] struct {
	data map[T]struct{}
}

// NewSet will return a Set object from an input slice
func NewSet[T comparable](data []T) Set[T] {
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
func NewSetWithCapacity[T comparable](data []T, size int) Set[T] {
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
