# set
Educating myself about making sets in go. Attempting to match the python
[set](https://docs.python.org/3/library/stdtypes.html#set-types-set-frozenset)
methods. All methods are written in more or less the naive way.

In operations where we may have to allocate more space for the underlying map, is it
faster to just iterate over the items and let the runtime allocate for us? Or is it
faster to figure out how much space we'll need first, then allocate manually, and then
add the various items to the new map?

