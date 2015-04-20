// Copyright 2015 The av-vortex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package skiplist

import (
	_ "math/rand"
	"testing"
)

func less(l, r interface{}) bool {
	return l.(int) < r.(int)
}

func TestLen(t *testing.T) {
	s := New(less)
	if s.Len() != 0 {
		t.Error("length should be 0")
	}

	s.Put(1, 2)
	s.Put(1, 123)
	s.Put(2, 33)
	if v := s.Len(); v != 2 {
		t.Errorf("length should be 2 returned %d", v)
	}

	if kv, ok := s.Remove(1); !ok || kv.Value != 123 {
		t.Error("value should be 123")
	}
	if kv, ok := s.Remove(1); ok || kv != nil {
		t.Error("value should be nil")
	}
}

func TestLevel(t *testing.T) {
	s := New(less)
	if s.level() != 0 {
		t.Error("level should be 0")
	}
}

/*
func TestRandomLevel(t *testing.T) {
	s := New(less)
	s.MaxLevel = 32
	if v := s.randomLevel(); v != 2 {
		t.Errorf("random level should be 2 returned %d", v)
	}
}
*/

func TestEmptyNodeNext(t *testing.T) {
	n := new(node)
	if next := n.next(); next != nil {
		t.Errorf("Next() should be nil for an empty node.")
	}
}

func TestGet(t *testing.T) {
	s := New(less)
	s.Put(0, 0)

	if value, ok := s.Get(0); !(value == 0 && ok) {
		t.Errorf("%v, %v instead of %v, %v", value, ok, 0, true)
	}

	if value, ok := s.Get(100); value != nil || ok {
		t.Errorf("%v, %v instead of %v, %v", value, ok, nil, false)
	}
}

func TestGreaterOrEqual(t *testing.T) {
	s := New(less)

	if kv := s.GreaterOrEqual(5); kv != nil  {
		t.Errorf("s.GreaterOrEqual(5) should have returned nil and nil for an empty map, not %v and %v.", kv.Key, kv.Value)
	}

	s.Put(0, 0)

	if kv := s.GreaterOrEqual(5); kv != nil  {
		t.Errorf("s.GreaterOrEqual(5) should have returned nil and nil for an empty map, not %v and %v.", kv.Key, kv.Value)
	}

	s.Put(10, 100)

	if kv := s.GreaterOrEqual(5); !(kv.Key == 10 && kv.Value == 100) {
		t.Errorf("s.GreaterOrEqual(5) should have returned 10 and 100, not %v and %v.", kv.Key, kv.Value)
	}
}

func (s *SkipList) equals(t *testing.T, key, wanted int) {
	if got, _ := s.Get(key); got != wanted {
		t.Errorf("For key %v wanted value %v, got %v.", key, wanted, got)
	}
}

func TestPut(t *testing.T) {
	s := New(less)
	if l := s.Len(); l != 0 {
		t.Errorf("Len is not 0, it is %v", l)
	}

	s.Put(0, 0)
	s.Put(1, 1)
	if l := s.Len(); l != 2 {
		t.Errorf("Len is not 2, it is %v", l)
	}
	s.equals(t, 0, 0)
	if t.Failed() {
		t.Errorf("header.Next() after s.Set(0, 0) and s.Set(1, 1): %v.", s.head.next())
	}
	s.equals(t, 1, 1)

}

func TestChange(t *testing.T) {
	s := New(less)
	s.Put(0, 0)
	s.Put(1, 1)
	s.Put(2, 2)

	s.Put(0, 7)
	if value, _ := s.Get(0); value != 7 {
		t.Errorf("Value should be 7, not %d", value)
	}
	s.Put(1, 8)
	if value, _ := s.Get(1); value != 8 {
		t.Errorf("Value should be 8, not %d", value)
	}

}

func TestRemove(t *testing.T) {
	s := New(less)
	for i := 0; i < 10; i++ {
		s.Put(i, i)
	}
	for i := 0; i < 10; i += 2 {
		s.Remove(i)
	}

	for i := 0; i < 10; i += 2 {
		if _, ok := s.Get(i); ok {
			t.Errorf("%d should not be present in s", i)
		}
	}

	if kv, ok := s.Remove(10000); kv != nil || ok {
		t.Errorf("Deleting a non-existent key should return nil, false, and not %v, %v.", kv.Value, ok)
	}

}

func TestIteration(t *testing.T) {
	s := New(less)
	for i := 0; i < 20; i++ {
		s.Put(i, i)
	}

	seen := uint(0)
	var lastKey int

	i := s.Iterator()

	for i.Next() {
		seen++
		lastKey = i.Key().(int)
		if i.Key() != i.Value() {
			t.Errorf("Wrong value for key %v: %v.", i.Key(), i.Value())
		}
	}

	if seen != s.Len() {
		t.Errorf("Not all the items in s where iterated through (seen %d, should have seen %d). Last one seen was %d.", seen, s.Len(), lastKey)
	}
/*
	for i.Prev() {
		if i.Key() != i.Value() {
			t.Errorf("Wrong value for key %v: %v.", i.Key(), i.Value())
		}

		if i.Key().(int) >= lastKey {
			t.Errorf("Expected key to descend but ascended from %v to %v.", lastKey, i.Key())
		}

		lastKey = i.Key().(int)
	}

	if lastKey != 0 {
		t.Errorf("Expected to count back to zero, but stopped at key %v.", lastKey)
	}
*/
}

func TestIterationSeek(t *testing.T) {
	s := New(less)
	for i := 0; i < 20; i++ {
		s.Put(i, i)
	}

	i := s.Iterator()

	if !i.Seek(5) {
		t.Errorf("Could not seek to key of value 5 got %v.", i.Key)
	}
}
