// Copyright 2015 The taigaDB Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package table

import (
	"encoding/binary"
)

// Iterator iterates over a Table's key/value pairs in key order.
type Iterator interface {
	// Is positioned at a valid node
	Valid() bool

	// Next moves the iterator to the next key/value pair.
	// It returns whether the iterator is exhausted.
	Next() bool

	// Key returns the key of the current key/value pair, or nil if done.
	// The caller should not modify the returned contents.
	Key() Slice

	// Value returns the value of the current key/value pair, or nil if done.
	// The caller should not modify the returned contents.
	Value() Slice
}

//---------------------------------------------------------------------------------------
// SSTable Iterator
//---------------------------------------------------------------------------------------

type ssTableIterator struct {
	sst *ssTable
	idx IndexSlice

	pos int
}

func (self ssTableIterator) Valid() bool {
	return self.pos != self.idx.Len()
}

func (self *ssTableIterator) Next() bool {
	if self.Valid() {
		self.pos++
	}
	
	return self.Valid()
}

func (self ssTableIterator) Key() Slice {
	if self.Valid() {
		return self.idx[self.pos].Key
	}
	return nil
}

func (self ssTableIterator) Value() Slice {
	if self.Valid() {
		return nil
	}
	return nil
}

//---------------------------------------------------------------------------------------
// Block Entry Iterator
//---------------------------------------------------------------------------------------

type EntryIterator interface {
	Next() (*BlockEntry, bool)
}

type blockEntryIterator struct {
	data Block

	entry BlockEntry
}

func (self *blockEntryIterator) Next() (*BlockEntry, bool) {
	if len(self.data) == 0 {
		return nil, false
	}

	b := self.data
	entry := &self.entry

	var n0, n1, n2 int
	entry.Shared,   n0 = binary.Uvarint(b)
	entry.Unshared, n1 = binary.Uvarint(b[n0:])
	entry.ValueLen, n2 = binary.Uvarint(b[n0 + n1:])
	
	n := n0 + n1 + n2
	entry.Key   = Slice(append(entry.Key[:int(entry.Shared)], b[n:n + int(entry.Unshared)]...))
	entry.Value = Slice(b[n + int(entry.Unshared):n + int(entry.Unshared + entry.ValueLen)])

	self.data = b[n + int(entry.Unshared + entry.ValueLen):]
	
	return entry, true
}

func NewEntryIterator(b Block) EntryIterator {
	return &blockEntryIterator{
		data: b[:b.RestartStartOffset()],
	}
}
