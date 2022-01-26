package intset

import (
	"testing"
)

func assert(t *testing.T, result bool, message string, vals ...interface{}) {
	if !result {
		t.Fatalf(message, vals...)
	}
}

func TestGrowSetAddAndContains(t *testing.T) {
	set := NewGrowSet(6)

	set.Add(1)
	set.Add(3)
	set.Add(4)

	for _, v := range []int{0, 2, 5, 6} {
		assert(t, !set.Contains(v), "set should not contain %v", v)
	}

	for _, v := range []int{1, 3, 4} {
		assert(t, set.Contains(v), "set should contain %v", v)
	}
}

func TestGrowSetRepeatedAdd(t *testing.T) {
	set := NewGrowSet(6)

	for i := 0; i < 100; i++ {
		set.Add(1)
	}

	assert(t, set.Contains(1), "set doesn't contain 1")
	assert(t, !set.Contains(0), "set contains 0")
}

func TestGrowSetClear(t *testing.T) {
	set := NewGrowSet(6)

	set.Add(3)
	set.Add(4)

	for _, v := range []int{0, 1, 2, 5} {
		assert(t, !set.Contains(v), "set should not contain %v", v)
	}

	for _, v := range []int{3, 4} {
		assert(t, set.Contains(v), "set should contain %v", v)
	}

	set.Clear()

	for i := 0; i < 10; i++ {
		assert(t, !set.Contains(i), "set contains %v", i)
	}
}

func TestGrowSetValues(t *testing.T) {
	set := NewGrowSet(6)
	set.Add(0)
	set.Add(2)
	set.Add(4)

	for _, v := range set.Values() {
		assert(t, v == 0 || v == 2 || v == 4, "unknown value in set %v", v)
	}

	set.Clear()

	for range set.Values() {
		assert(t, false, "values present after clear")
	}
}

func TestGrowSetSize(t *testing.T) {
	set := NewGrowSet(6)

	set.Add(0)
	set.Add(2)
	set.Add(4)

	assert(t, set.Size() == 3, "set size isn't 3")

	set.Add(0)
	set.Add(2)
	set.Add(4)

	assert(t, set.Size() == 3, "set size should still be 3")

	set.Clear()

	assert(t, set.Size() == 0, "set size should be 0")
}

func TestGrowSetPop(t *testing.T) {
	set := NewGrowSet(6)
	set.Add(0)
	set.Add(2)
	set.Add(4)

	oldSize := set.Size()
	assert(t, oldSize == 3, "size isn't 3")

	popped, err := set.Pop()
	assert(t, err == nil, "error is not nil: %v", err)
	assert(t, popped == 0 || popped == 2 || popped == 4, "popped value isn't correct")
	assert(t, !set.Contains(popped), "set should not contain popped value")
	assert(t, set.Size() == oldSize-1, "set size is incorrect")

	popped1, err := set.Pop()
	assert(t, err == nil, "error is not nil: %v", err)
	assert(t, popped != popped1, "duplicate popped value")
	assert(t, popped1 == 0 || popped1 == 2 || popped1 == 4, "popped value isn't correct")

	popped2, err := set.Pop()
	assert(t, err == nil, "error is not nil: %v", err)
	assert(t, popped != popped2 && popped1 != popped2, "duplicate popped value")
	assert(t, popped2 == 0 || popped2 == 2 || popped2 == 4, "popped value isn't correct")

	assert(t, set.Size() == 0, "set size should be zero")

	_, err = set.Pop()
	assert(t, err != nil, "error should not be nil")
	assert(t, err == EmptySetError, "error should be EmptySetError")
}

func TestShrinkSetContainsAndSize(t *testing.T) {
	set := NewShrinkSet(5)
	assert(t, set.Size() == 5, "set size should be 5")

	for i := 0; i < 5; i++ {
		assert(t, set.Contains(i), "set should contain %v", i)
	}

	assert(t, !set.Contains(6), "set should not contain 6")
}

func TestShrinkSetRemove(t *testing.T) {
	set := NewShrinkSet(6)

	set.Remove(1)
	set.Remove(3)
	set.Remove(5)

	for _, v := range []int{0, 2, 4} {
		assert(t, set.Contains(v), "set should contain %v", v)
	}

	for _, v := range []int{1, 3, 5} {
		assert(t, !set.Contains(v), "set should not contain %v", v)
	}
}

func TestShrinkSetRefill(t *testing.T) {
	set := NewShrinkSet(6)

	set.Remove(1)
	set.Remove(3)
	set.Remove(5)

	assert(t, set.Size() == 3, "set size should be 3")

	for _, i := range []int{0, 2, 4} {
		assert(t, set.Contains(i), "set should contain %v", i)
	}

	for _, i := range []int{1, 3, 5} {
		assert(t, !set.Contains(i), "set should not contain %v", i)
	}

	set.Refill()

	for i := 0; i < 6; i++ {
		assert(t, set.Contains(i), "set should contain %v", i)
	}

	assert(t, set.Size() == 6, "set size should be 6")
}

func TestShrinkSetValues(t *testing.T) {
	set := NewShrinkSet(6)

	set.Remove(1)
	set.Remove(3)
	set.Remove(5)

	for _, i := range set.Values() {
		assert(t, i != 1 && i != 3 && i != 5, "value should not be %v", i)
	}

	set.Refill()

	count := 0
	for _, i := range set.Values() {
		count++
		assert(t, i >= 0 && i < 6, "value out of range: %v", i)
	}

	assert(t, count == 6, "count should be 6")
}

func TestShrinkSetPop(t *testing.T) {
	set := NewShrinkSet(4)
	set.Remove(3)

	last := -1
	for i := 0; i < 3; i++ {
		popped, err := set.Pop()
		assert(t, err == nil, "error should be nil")
		assert(t, popped == 0 || popped == 1 || popped == 2, "invalid popped value %v", popped)
		assert(t, popped != last, "duplicate popped value")
		last = popped
	}

	set.Refill()

	last = -1
	for i := 0; i < 4; i++ {
		popped, err := set.Pop()
		assert(t, err == nil, "error should be nil")
		assert(t, popped == 0 || popped == 1 || popped == 2 || popped == 3, "invalid popped value %v", popped)
		assert(t, popped != last, "duplicate popped value")
		last = popped
	}

	set.Refill()
	for i := 0; i < 4; i++ {
		_, err := set.Pop()
		assert(t, err == nil, "error should be nil")
	}

	_, err := set.Pop()
	assert(t, err == EmptySetError, "error should be EmptySetError")
}
