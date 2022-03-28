# set
Educating myself about making generic sets in go. Attempting to match the python
[set](https://docs.python.org/3/library/stdtypes.html#set-types-set-frozenset)
API. For version 1.0, all methods are written in more or less the naive way. Elements of
a set must be `comparable`. `Comparable` is an interface that is implemented by all
comparable types (booleans, numbers, strings, pointers, channels, arrays of comparable
types, structs whose fields are all comparable types).

## API
Create a new `Set` object with `NewSet()` or `NewSetWithCapacity()`. They take a slice,
or any type that has a slice as the underlying data type, e.g. `type MyList = []int` is
an acceptable input type.

```go
package main

import (
	"fmt"

	"github.com/natemcintosh/set"
)

func main() {
	// The input slice (notice the duplicates)
	raw_data_1 := []int{1, 2, 1, 4, 5, 6, 6, 6, 7}
	// Create the Set
	s1 := set.NewSet(raw_data_1)
	// See that duplicate values no longer exist
	fmt.Println(s1)
	// {map[1:{} 2:{} 4:{} 5:{} 6:{} 7:{}]}

	// Check if an item exists in the set
	if s1.Contains(10) {
		fmt.Println("10 is in the set")
	} else {
		fmt.Println("10 is not in the set")
		// 10 is not in the set
	}

	// How many items are in the set?
	fmt.Printf("The set has %v items\n", s1.Len())
	// The set has 6 items

	// Is the set empty?
	fmt.Printf("The set is empty: %t\n", s1.IsEmpty())
	// The set is empty: false

	// Add new items to the set
	s1.Add(3)
	s1.Add(10)
	s1.Add(1) // Does nothing since 1 already exists in the set

	// Remove items, and return an error if it doesn't exist
	err := s1.Remove(100)
	if err != nil {
		fmt.Println("The set did not have the item 100. Failed to remove it from the set.")
		// The set did not have the item 100. Failed to remove it from the set.
	}

	// Discard items, and don't worry about if it exists in the set or not
	s1.Discard(1)
	fmt.Printf("The set contains 1: %t\n", s1.Contains(1))
	// The set contains 1: false

	// Pop a random item from the set. Get an error if the set is empty
	item, err := s1.Pop()
	if err != nil {
		fmt.Printf("The set was empty")
	}
	fmt.Printf("The item popped from the set was: %v\n", item)
	// The item popped from the set was: 10

	// Clear all items from the set
	s1.Clear()
	fmt.Printf("The set is empty: %t\n", s1.IsEmpty())
	// The set is empty: true

	// Create a new set with capacity for 50 numbers. Useful if we know that this will
	// be filled up with at least that many items in the future.
	// Also note that we can pass in a type that has a slice as the underlying type.
	type some_floats []float64
	var numbers some_floats = []float64{1.0, 2.0, 3.14, 4.95}
	s2 := set.NewSetWithCapacity(numbers, 50)

	// Create a deep copy of s2
	s3 := s2.Copy()

	// Make `s3` a little different than `s2`
	s3.Discard(1.0)
	s3.Discard(2.0)
	s3.Add(5.67)
	s3.Add(-100.21)

	// Are the two sets equal?
	fmt.Printf("s2 equals s3?: %t\n", s2.Equals(s3))
	// s2 equals s3?: false

	// Note that we cannot compare two sets of different types
	// s1.Equals(s2)
	// Compile error:
	// cannot use s2 (variable of type set.Set[float64]) as type set.Set[int] in argument to s1.Equals

	// What's the union of the two sets
	union := s2.Union(s3)
	fmt.Printf("Created a new set (deep copy), containing: %v\n", union)
	// Created a new set (deep copy), containing: {map[-100.21:{} 1:{} 2:{} 3.14:{} 4.95:{} 5.67:{}]}

	// Union, but update the set in place. This will be faster as it doesn't create a copy
	s2.UnionInPlace(s3)
	fmt.Printf("s2 has become: %v\n", s2)
	// s2 has become: {map[-100.21:{} 1:{} 2:{} 3.14:{} 4.95:{} 5.67:{}]}

	// Recreate original `s2`
	s2 = set.NewSet(numbers)

	// What is the intersection of the two sets?
	intersection := s2.Intersection(s3)
	fmt.Printf("s2 intersection with s3 = %v\n", intersection)
	// s2 intersection with s3 = {map[3.14:{} 4.95:{}]}

	// IntersectionInPlace will update a set in place. Faster than Intersection
	s2.IntersectionInPlace(s3)
	fmt.Printf("s2 has become: %v\n", s2)
	// s2 has become: {map[3.14:{} 4.95:{}]}

	// Recreate original `s2`
	s2 = set.NewSet(numbers)

	s4 := set.NewSet([]float64{-10, -9})
	// Are the two sets disjoint? As in, are there no elements common to both?
	fmt.Printf("s2 has no common elements with s4? %t\n", s2.IsDisjoint(s4))
	// s2 has no common elements with s4? true

	s3 = s2.Copy()
	// Is `s3` a subset of `s2`?
	fmt.Printf("s3 is a subset of s2? %t\n", s3.IsSubsetOf(s2))
	// s3 is a subset of s2? true

	// Proper subset
	fmt.Printf("s3 is a proper subset of s2? %t\n", s3.IsProperSubsetOf(s2))
	// s3 is a proper subset of s2? false

	s3.Add(23.45)
	// Is `s3` a super set of `s2`
	fmt.Printf("s3 is a super set of s2? %t\n", s3.IsSuperSetOf(s2))
	// s3 is a super set of s2? true

	// Proper superset
	fmt.Printf("s3 is a proper super set of s2? %t\n", s3.IsProperSuperSetOf(s2))
	// s3 is a proper super set of s2? true

	// What is the difference of two sets
	taking_science := set.NewSet([]string{"Larry", "Curly", "Moe", "Shemp"})
	taking_math := set.NewSet([]string{"Curly", "Shemp", "Albert"})
	not_taking_math := taking_science.Difference(taking_math)
	fmt.Printf("The following students take science but not math: %v\n", not_taking_math)
	// The following students take science but not math: {map[Larry:{} Moe:{}]}

	// DifferenceInPlace is faster, but alters the original
	taking_science.DifferenceInPlace(taking_math)
	fmt.Printf("taking_science is now: %v\n", taking_science)
	// taking_science is now: {map[Larry:{} Moe:{}]}

	// Reset `taking_science`
	taking_science = set.NewSet([]string{"Larry", "Curly", "Moe", "Shemp"})

	// What is the symmetric difference? I.e. not in the intersection of the two
	sym_diff := taking_science.SymmetricDifference(taking_math)
	fmt.Printf("The following students are taking science or math, but not both: %v\n", sym_diff)
	// The following students are taking science or math, but not both: {map[Albert:{} Larry:{} Moe:{}]}

	// Once again, the in place version is faster, but alters the original
	taking_science.SymmetricDifferenceInPlace(taking_math)
	fmt.Printf("taking_science has become %v\n", taking_science)
	// taking_science has become {map[Albert:{} Larry:{} Moe:{}]}
}
```


### Questions for later
In operations where we may have to allocate more space for the underlying map, is it
faster to just iterate over the items and let the runtime allocate for us? Or is it
faster to figure out how much space we'll need first, then allocate manually, and then
add the various items to the new map? For now, I'm guessing that it probably depends
heavily on the use case.
