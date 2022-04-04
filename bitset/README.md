# bitset

This module is inspired by the julia language's BitSet type and it's methods. This is
currently implemented as a slice of `bool`s, each index representing a spot on the
numberline. It is designed for dense integer sets. If the set will be sparse (for 
example, holding a few very large integers), use `set` instead. For dense integer sets,
`bitset` can be significantly faster than `set`.

### Future Improvements
Julia uses a `Vector{UInt64}` so that bitwise operations can be used. Ideally switch to
a `[]uint64`, the go equivalent of what julia uses to also take advantage of fast
bitwise operations.