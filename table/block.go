// Copyright 2015 The taigaDB Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package table

import (
	"encoding/binary"
	"fmt"
	"sort"
)

const (
	// LevelDB format magic number. leading 64bits of
	//  echo http://code.google.com/p/leveldb/ | sha1sum
	TableMagicNumber = uint64(0xdb4775248b80fb57)

	// 1-byte type + 32-bit crc
	BlockTrailerSize = 5

	// Maximum encoding length of a BlockHandle
	MaxEncodedLength = binary.MaxVarintLen64 * 2

	// Encoded length of a Footer.
	FooterEncodedLength = 2 * MaxEncodedLength + 8
)

//---------------------------------------------------------------------------------------
// BlockHandle
//---------------------------------------------------------------------------------------

type BlockHandle struct {
	Offset uint64
	Size   uint64
}

func NewHandle(offset, size uint64) *BlockHandle {
	return &BlockHandle{
		Offset: offset,
		Size: size,
	}
}

func (self BlockHandle) String() string {
	return fmt.Sprintf("BlockHandle { Offset: %d, Size: %d}", self.Offset, self.Size)
}

// Encodes to data and returns the number of bytes written. If the buffer 
// is too small, it will panic.
func (self BlockHandle) Encode(data Slice) (int, error) {
	if len(data) < MaxEncodedLength {
		return 0, EncodeBlockHandleBufferErr
	}

	n := binary.PutUvarint(data, self.Offset)
	m := binary.PutUvarint(data[n:], self.Size)

	return n + m, nil
}

// Decodes from data and returns the number of bytes read
func (self *BlockHandle) Decode(data Slice) (int, error) {
	var n, m int

	self.Offset, n = binary.Uvarint(data)
	switch {
	case  n == 0 :
		return 0, DecodeSmallBufferErr
	case  n < 0 :
		return 0, DecodeNot64bitsErr
	}

	self.Size, m = binary.Uvarint(data[n:])
	switch {
	case  m == 0 :
		return 0, DecodeSmallBufferErr
	case  m < 0 :
		return 0, DecodeNot64bitsErr
	}

	return n + m, nil
}

//---------------------------------------------------------------------------------------
// Footer
//---------------------------------------------------------------------------------------

type Footer struct {
	MetaIndexHandle  *BlockHandle
	BlockIndexHandle *BlockHandle
	MagicNumber uint64
}

func NewFooter(metaIndexHandle, blockIndexHandle *BlockHandle) *Footer {
	return &Footer{
		MetaIndexHandle: metaIndexHandle,
		BlockIndexHandle: blockIndexHandle,
		MagicNumber: TableMagicNumber,
	}
}

func (self Footer) String() string {
	return fmt.Sprintf("Footer { MetaIndexHandle: %v, BlockIndexHandle: %v, MagicNumber: %d}", 
		              self.MetaIndexHandle, 
		              self.BlockIndexHandle,
		              self.MagicNumber)
}

// Encodes to data and returns the number of bytes written. If the buffer 
// is too small, it will panic.
func (self Footer) Encode(data Slice) (int, error) {
	if len(data) < FooterEncodedLength {
		return 0, EncodeFooterBufferErr
	}

	var pos int

	if metaSize,  _ := self.MetaIndexHandle.Encode(data); metaSize != 0 {
		pos += metaSize
	}

	if blockSize, _ := self.BlockIndexHandle.Encode(data[pos:]); blockSize != 0 {
		pos += blockSize
	}

	binary.LittleEndian.PutUint64(data[40:], TableMagicNumber)
	return FooterEncodedLength, nil
}

// Decodes from data and returns the number of bytes read
func (self *Footer) Decode(data Slice) (int, error) {
	var n, m int
	var err error

	if self.MetaIndexHandle == nil {
		self.MetaIndexHandle = NewHandle(0, 0)
	}

	if n, err = self.MetaIndexHandle.Decode(data); err != nil {
		return n, err
	}

	if self.BlockIndexHandle == nil {
		self.BlockIndexHandle = NewHandle(0, 0)
	}

	if m, err = self.BlockIndexHandle.Decode(data[n:]); err != nil {
		return m, err
	}

	self.MagicNumber = binary.LittleEndian.Uint64(data[40:])
	return FooterEncodedLength, nil
}

//---------------------------------------------------------------------------------------
// Block
//---------------------------------------------------------------------------------------

/* 
Blocks have one or many key/value entries followed by a block trailer structure.
*/
type Block []byte

func (self Block) NumberOfRestarts() int {
	return int(binary.LittleEndian.Uint32(self[len(self) - 4:]))
}

func (self Block) RestartStartOffset() int {
	return len(self) - (1 + int(self.NumberOfRestarts())) * 4
}

func (self Block) RestartSlice() []uint32 {
	numRestarts := self.NumberOfRestarts()
	if numRestarts == 0  {
		return nil
	}

	restarts := make([]uint32, numRestarts)

	data := self[self.RestartStartOffset():]
	for i := 0; i < numRestarts; i++ {
		restarts[i] = binary.LittleEndian.Uint32(data[4 * i:])
	}
	return restarts
}

func (self Block) Search(key Slice) *BlockEntry {
	numRestarts := self.NumberOfRestarts()
	if numRestarts == 0  {
		return nil
	}

	restarts := self.RestartSlice()

	var entry BlockEntry
	// Search uses binary search to find and return the smallest index
	pos := sort.Search(numRestarts, func(i int) bool {
		readBlockEntry(self[restarts[i]:], &entry)

		return string(entry.Key) > string(key)	
	})

	if pos < numRestarts {
		// Goto previous restart because we found Key > restart(key)
		pos--
		
		iter := NewEntryIterator(self[restarts[pos]:]) 
		for e, ok := iter.Next(); ok; { 
			if string(e.Key) == string(key) {
				return e
			}
			e, ok = iter.Next()
		}

	}

	return nil
}

// Prints block entries for debuging porposes
func (self Block) Dump() {
	iter := NewEntryIterator(self) 
	for entry, ok := iter.Next(); ok; { 
		fmt.Println(entry)

		entry, ok = iter.Next()
	}
}

//---------------------------------------------------------------------------------------
// Block Entry
//---------------------------------------------------------------------------------------

type BlockEntry struct {
	Shared    uint64
	Unshared  uint64
	ValueLen  uint64

	Key    Slice
	Value  Slice
}

func (self BlockEntry) String() string {
	return fmt.Sprintf("BlockEntry { Shared: %v, Unshared: %v, ValueLen: %v, Key: %s, Value: %v}", 
		              self.Shared, 
		              self.Unshared,
		              self.ValueLen,
		              self.Key,
		              self.Value)
}

// Helper function to decode the next block entry. 
// It returns the truncated block slice that was read.
func readBlockEntry(b Block, entry *BlockEntry) Block {
	var n0, n1, n2 int
	entry.Shared,   n0 = binary.Uvarint(b)
	entry.Unshared, n1 = binary.Uvarint(b[n0:])
	entry.ValueLen, n2 = binary.Uvarint(b[n0 + n1:])
	
	n := n0 + n1 + n2
	entry.Key   = Slice(append(entry.Key[:int(entry.Shared)], b[n:n + int(entry.Unshared)]...))
	entry.Value = Slice(b[n + int(entry.Unshared):n + int(entry.Unshared + entry.ValueLen)])

	return b[n + int(entry.Unshared + entry.ValueLen):]
}

//---------------------------------------------------------------------------------------
// Index Entry
//---------------------------------------------------------------------------------------

/*
The index entry contains the last key in that data block and is less then the first key in 
the successive data block.  The value is the BlockHandle for the data block.
*/
type IndexEntry struct {
	Key Slice
	Handle BlockHandle
}

func (self IndexEntry) String() string {
	return fmt.Sprintf("IndexEntry { Key: %8s, Handle: %v}", self.Key, self.Handle)
}

// Helper function to decode the index entries. 
// It returns an index slice.
func decodeIndexEntries(b Block) IndexSlice {
	numRestarts := b.NumberOfRestarts()
	idxSlice := make(IndexSlice, numRestarts)

	for i := 0; i < numRestarts; i++ {
		shared,   n0 := binary.Uvarint(b)
		unshared, n1 := binary.Uvarint(b[n0:])
		valueLen, n2 := binary.Uvarint(b[n0 + n1:])
	
		var ie = new(IndexEntry)

		n := n0 + n1 + n2
		ie.Key =  Slice(append(ie.Key[:int(shared)], b[n:n + int(unshared)]...))
		value  := Slice(b[n + int(unshared):n + int(unshared + valueLen)])
		b = b[n + int(unshared + valueLen):]
	
		ie.Handle.Decode(value)
	
		idxSlice[i] = ie
	}

	return idxSlice
}

type IndexSlice []*IndexEntry

func (self IndexSlice) Len() int { 
	return len(self) 
}

func (self IndexSlice) Less(i, j int) bool { 
	return string(self[i].Key) < string(self[j].Key) 
}

func (self IndexSlice) Swap(i, j int) { 
	self[i], self[j] = self[j], self[i] 
}

// Sort is a convenience method.
func (self IndexSlice) Sort() { 
	sort.Sort(self) 
}

// Search returns the result of applying Search to the receiver and x. 
func (self IndexSlice) Search(key Slice) int { 
	return sort.Search(len(self), func(i int) bool { return string(self[i].Key) >= string(key) })
}