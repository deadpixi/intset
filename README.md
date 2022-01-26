# intset

This package provides unsigned integer sets with various useful properties.
All of the data structures in this package support the following
operations with the associated time complexity:

- `Contains(n)` - Checks if *n* is a member of the set, in *O(1)* time
- `Pop()`       - Remove and return an arbitrary integer from the set, in *O(1)* time
- `Size()`      - Return the number of items in the set, in *O(1)* time
- `Values()`    - Returns a slice of integers of the members of the set

None of the data structures in this package allocate or deallocate memory
after construction.

The various data structures provide other operations which may be useful
in different situations.

# GrowSet

A `GrowSet` starts out empty and can have items added to it.
It supports the following additional operations with the associated time
complexity:

- `Add(n)`  - Add integer *n* to the set, in *O(1)* time.
- `Clear()` - Removes all elements from the set, in *O(1)* time.

`GrowSet` is based on *An Efficient Representation for Sparse Sets* by Briggs and Torczon.

# ShrinkSet

A `ShrinkSet` starts out containing a set of numbers, and
can have numbers removed. It supports the following additional
operations:

- `Remove(n)` - Remove *n* from the set in *O(1)* time.
- `Refill()` - Refills the set in *O(1)* time.

`ShrinkSet` is, as far as I know, a novel data structure.