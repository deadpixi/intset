// This package provides unsigned integer sets with various useful properties.
// All of the data structures in this package support the following
// operations with the associated time complexity:
//   Contains(n) - Checks if n is a member of the set, in O(1) time
//   Pop()       - Remove and return an arbitrary integer from the set, in O(1) time
//   Size()      - Return the number of items in the set, in O(1) time
//   Values()    - Returns a slice of integers of the members of the set
// None of the data structures in this package allocate or deallocate memory
// after construction.
//
// The various data structures provide other operations which may be useful
// in different situations.
package intset

import (
	"errors"
)

// Returned when an operation (e.g. Pop) that returns a value from the set
// is requested on an empty set.
var EmptySetError = errors.New("empty set")

// Returned when a value is too large or small to fit in a constructed set.
var ValueOutOfRangeError = errors.New("value out of range")

type set struct {
	n      int
	sparse []int
	dense  []int
}

// A GrowSet starts out empty and can have items added to it.
// It supports the following additional operations with the associated time
// complexity:
//
//   Add(n)  - Add integer n to the set, in O(1) time.
//   Clear() - Removes all elements from the set, in O(1) time.
type GrowSet set

// Allocate a new GrowSet.
// The resulting set will be able to store the integers less than
// capacity.
// Construction takes O(1) time.
func NewGrowSet(capacity int) *GrowSet {
	return &GrowSet{
		n:      0,
		sparse: make([]int, capacity, capacity),
		dense:  make([]int, capacity, capacity),
	}
}

// Returns true if value is a member of the set.
func (g *GrowSet) Contains(value int) bool {
	return value < len(g.sparse) && g.sparse[value] < g.n && g.dense[g.sparse[value]] == value
}

// Removes all elements from the set.
func (g *GrowSet) Clear() {
	g.n = 0
}

// Returns the size of the set.
func (g *GrowSet) Size() int {
	return g.n
}

// Adds value to the set. Adding the same value multiple times is not an error.
// If a value is less than zero or too large to be stored in the set, ValueOutOfRangeError
// is returned, otherwise nil.
func (g *GrowSet) Add(value int) error {
	if value >= len(g.sparse) || value < 0 {
		return ValueOutOfRangeError
	}

	if !g.Contains(value) {
		g.dense[g.n] = value
		g.sparse[value] = g.n
		g.n++
	}

	return nil
}

// Remove and return a random value from the set.
// If the set is empty, the result will be 0 and error will be EmptySetError.
func (g *GrowSet) Pop() (int, error) {
	if g.n == 0 {
		return 0, EmptySetError
	}

	value := g.dense[g.n-1]
	g.n--
	return value, nil
}

// Returns a slice of ints, which are the members of the set.
// This slice should not be modified.
func (g *GrowSet) Values() []int {
	return g.dense[:g.n]
}

// A ShrinkSet starts out containing a set of numbers, and
// can have numbers removed. It supports the following additional
// operations:
//
//   Remove(n) - Remove n from the set in O(1) time.
//   Refill() - Refills the set in O(1) time.
type ShrinkSet set

// Create a new ShrinkSet storing the numbers up to,
// but not including, capacity. This takes O(n) time,
// where n == capacity.
func NewShrinkSet(capacity int) *ShrinkSet {
	result := &ShrinkSet{
		n:      capacity,
		sparse: make([]int, capacity, capacity),
		dense:  make([]int, capacity, capacity),
	}

	for i := 0; i < capacity; i++ {
		result.sparse[i] = i
		result.dense[i] = i
	}

	return result
}

// Returns true if value is in the set.
func (s *ShrinkSet) Contains(value int) bool {
	return value < len(s.sparse) && s.sparse[value] < s.n
}

// Resets the set to its original state in O(1) time.
func (s *ShrinkSet) Refill() {
	s.n = len(s.dense)
}

// Returns the number of elements in the set.
func (s *ShrinkSet) Size() int {
	return s.n
}

// Returns a slice containing the members of the set.
// This slice should not be modified.
func (g *ShrinkSet) Values() []int {
	return g.dense[:g.n]
}

// Remove the item from the set. It is not an error to
// remove an item that does not exist.
func (s *ShrinkSet) Remove(item int) {
	if s.Contains(item) {
		itemIndex := s.sparse[item]
		lastItem := s.dense[s.n-1]
		lastItemIndex := s.sparse[lastItem]

		s.dense[lastItemIndex] = item
		s.dense[itemIndex] = lastItem
		s.sparse[lastItem] = itemIndex
		s.sparse[item] = lastItemIndex
		s.n--
	}
}

// Remove and return a random member from the set.
// If the set is empty, the result will be zero and
// error will be EmptySetError.
func (s *ShrinkSet) Pop() (int, error) {
	if s.n == 0 {
		return 0, EmptySetError
	}

	removed := s.dense[0]
	s.Remove(removed)
	return removed, nil
}
