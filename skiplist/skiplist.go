// Copyright 2015 The av-vortex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package skiplist

import (
	"math/rand"
)

const (
	// p is the fraction of nodes with level i pointers that also have
	// level i+1 pointers. p equal to 1/4 is a good value from the point
	// of view of speed and space requirements. If variability of running
	// times is a concern, 1/2 is a better value for p.
	defaultP = 0.25

	defaultMaxLevel = 32
)

type LessFunc func(l, r interface{}) bool

// Interface that you can use to implement an iterator that iterates 
// through a skip list
type Iterator interface {
	// Is positioned at a valid node
	Valid() bool

	// Key returns the current key.
	Key() interface{}
	
	// Value returns the current value.
	Value() interface{}

	// Advances to the next position and returns true if valid node
	Next() bool

	// Advances to the previous position and returns true if valid node
	Prev() bool

	// Advance to the first entry with a key >= target
	Seek(key interface{}) bool
}

type iter struct {
	current *node
	list    *SkipList
}

func (self iter) Valid() bool {
	return self.current != nil
}

func (self iter) Key() interface{} {
	if self.Valid() {
		return self.current.Key
	}
	return nil
}

func (self iter) Value() interface{} {
	if self.Valid() {
		return self.current.Value
	}
	return nil
}

func (self *iter) Next() bool {
	if next := self.current.next(); next != nil {
		self.current = next
		return true
	}
	return false
}

func (self *iter) Prev() bool {
	// Instead of using explicit "prev" links, we just search for the
	// last node that falls before key.
	if !self.Valid() {
		return false
	}

//	node := self.list.FindLessThan(self.current.key)
//	return node != self.list.head
	return false
}

func (self *iter) Seek(key interface{}) bool {
	current := self.current

	// If the existing iterator outside of the known key range, we should set the
	// position back to the beginning of the list.
	if current == nil {
		current = self.list.head
	}

	self.current = self.list.findNode(current, nil, key)

	return self.current != nil
}

//---------------------------------------------------------------------------------------
//
//---------------------------------------------------------------------------------------

type KV struct {
	Key     interface{}
	Value   interface{}
}

//---------------------------------------------------------------------------------------
//
//---------------------------------------------------------------------------------------

type node struct {
	*KV

	forward []*node
}

func (self *node) next() *node {
	if len(self.forward) == 0 {
		return nil
	}

	return self.forward[0]
}

//---------------------------------------------------------------------------------------
//
//---------------------------------------------------------------------------------------

type SkipList struct {
	head     *node
	length   uint

	P        float64
	MaxLevel int
	less     LessFunc
}

// Create a new SkipList object that will use "cmp" for comparing keys
func New(less LessFunc) *SkipList {
	return &SkipList{
		head: &node{
			KV: new(KV),
			forward: []*node{nil},
		},
		length: 0,
		P: defaultP,
		MaxLevel: defaultMaxLevel,
		less: less,
	}
}

func (self *SkipList) level() int {
	// Returns the level-1 of the skip list, used for slices indices.
	// The level of an empty skip list is 1.
	return len(self.head.forward) - 1
}

func (self *SkipList) randomLevel() (n int) {
	// Returns a random level in the range [0, s.level()+1] been at most
	// equal to s.maxLevel-1. Used for slices indices.
	for n = 0; rand.Float64() < self.P && n < self.MaxLevel - 1; n++ {
	}
	return
}

// The length of the skip list
func (self *SkipList) Len() uint {
	return self.length
}

// Returns the value associated with key. 
func (self *SkipList) Get(key interface{}) (interface{}, bool) {
	
	if candidate := self.findNode(self.head, nil, key); candidate != nil && candidate.Key == key {
		return candidate.Value, true
	}

	return nil, false
}

// Put key into the list, existing key is replaced
func (self *SkipList) Put(key, value interface{}) {
	update := make([]*node, self.level() + 1)
	candidate := self.findNode(self.head, update, key)

	if candidate != nil && candidate.Key == key {
		candidate.Value = value
		return
	}

	newLevel := self.randomLevel()

	if level := self.level(); newLevel > level {
		for i := level + 1; i <= newLevel; i++ {
			update = append(update, self.head)
			self.head.forward = append(self.head.forward, nil)
		}
	}

	node := &node{
		KV: &KV{key, value},
		forward: make([]*node, newLevel + 1), 
	}
	for i := 0; i <= newLevel; i++ {
		node.forward[i] = update[i].forward[i]
		update[i].forward[i] = node
	}

	self.length++
}

// True if an entry that compares equal to key is in the list
func (self *SkipList) Contains(key interface{}) bool {
	_, ok := self.Get(key)
	return ok
}

func (self *SkipList) GreaterOrEqual(key interface{}) *KV {

	if candidate := self.findNode(self.head, nil, key); candidate != nil {
		return candidate.KV
	}
	return nil
}

// Returns an iterator
func (self *SkipList) Iterator() Iterator {
	return &iter{
		current: self.head,
		list: self,
	}
}

// Removes the key from the list.
func (self *SkipList) Remove(key interface{}) (*KV, bool) {
	if key == nil {
		return nil, false
	}

	update := make([]*node, self.level() + 1)
	candidate := self.findNode(self.head, update, key)

	if candidate == nil || candidate.Key != key {
		return nil, false
	}

	for i := 0; i <= self.level() && update[i].forward[i] == candidate; i++ {
		update[i].forward[i] = candidate.forward[i]
	}

	for self.level() > 0 && self.head.forward[self.level()] == nil {
		self.head.forward = self.head.forward[:self.level()]
	}

	self.length--
	return candidate.KV, true
}

func (self *SkipList) Min() *KV {
	if min := self.head.next(); min != nil {
		return min.KV
	}
	return nil
}

func (self *SkipList) Max() *KV {
	current := self.head
	for i := self.level(); i >= 0; i-- {
		for current.forward[i] != nil {
			current = current.forward[i]
		}
	}

	if current == self.head {
		return nil
	}
	return current.KV
}

// findNode populates update with nodes that constitute the path to the
// node that may contain key. 
//
// The candidate node will be returned. If update is nil, it will be not used
// (the candidate node will still be returned). If update is not nil, but it 
// doesn't have enough height (levels) for all the nodes in the path, 
// findNode will panic.
func (self *SkipList) findNode(current *node, update []*node, key interface{}) *node {
	depth := len(current.forward) - 1

	for i := depth; i >= 0; i-- {
		for current.forward[i] != nil && self.less(current.forward[i].Key, key) {
			current = current.forward[i]
		}
		if update != nil {
			update[i] = current
		}
	}
	return current.next()
}