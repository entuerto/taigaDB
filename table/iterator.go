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
	idx *IndexEntry
}

func (self ssTableIterator) Valid() bool {
	return self.idx != nil
}

func (self *ssTableIterator) Next() bool {
	if self.Valid() {
		self.idx = self.idx.next
	}
	
	return self.Valid()
}

func (self ssTableIterator) Key() Slice {
	if self.Valid() {
		return self.idx.Key
	}
	return nil
}

func (self ssTableIterator) Value() Slice {
	if self.Valid() {
		return self.sst.get(self.idx)
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

	var n0, n1, n2 int
	self.entry.Shared,   n0 = binary.Uvarint(self.data)
	self.entry.Unshared, n1 = binary.Uvarint(self.data[n0:])
	self.entry.ValueLen, n2 = binary.Uvarint(self.data[n0 + n1:])
	
	n := n0 + n1 + n2
	self.entry.Key   = Slice(append(self.entry.Key[:int(self.entry.Shared)], self.data[n:n + int(self.entry.Unshared)]...))
	self.entry.Value = Slice(self.data[n + int(self.entry.Unshared):n + int(self.entry.Unshared + self.entry.ValueLen)])

	self.data = self.data[n + int(self.entry.Unshared + self.entry.ValueLen):]
	
	return &self.entry, true
}

func NewEntryIterator(b Block) EntryIterator {
	return &blockEntryIterator{
		data: b[:b.RestartStartOffset()],
	}
}
