# bitset

This module is inspired by the julia language's BitSet type and it's methods. This is
currently implemented as
```go
type key struct {
	is_positive bool
	multiplier  uint64
}

type Set struct {
	data map[key]uint64
}
```
The `uint64` in the value of the map is used as a container; each bit indicates if a
number is stored in that index. The two fields in the `key` struct tell you which chunk
of 64 consecutive numbers is stored in that `uint64`.

This is a "hybrid" approach.

Cons:
1. We have to use a map as the underlying data container
1. We lose the possibility of straight iteration over all the bits in the set

Pros:
1. We can handle sparse sets relatively well
1. We still get access to bit operations for fast comparisons of `uint64` containers

### Questions for later
Instead of using `uint64` as the value type in the map, could also use a fixed size 
array of `uint64`. This might have some benefits if there are larger continuous runs of
numbers.
```go
type key struct {
	is_positive bool
	multiplier  uint64
}

type Set struct {
	data map[key][5]uint64
}
```