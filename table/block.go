// Copyright 2015 The taigaDB Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package table

import (
	"encoding/binary"
	"errors"
	"fmt"
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

var (
	EncodeBlockHandleBufferErr = errors.New("BlockHandle.Encode: Buffer too small")
	EncodeFooterBufferErr      = errors.New("Footer.Encode: Buffer too small")
	DecodeSmallBufferErr = errors.New("Decode: Buffer to small")
	DecodeNot64bitsErr   = errors.New("Decode: Value is not 64bits")
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

type Block []byte

func (self Block) NumberOfRestarts() uint32 {
	return binary.LittleEndian.Uint32(self[len(self) - 4:])
}

func (self Block) RestartStartOffset() int {
	return len(self) - (1 + int(self.NumberOfRestarts())) * 4
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

//---------------------------------------------------------------------------------------
// Block Index Entry
//---------------------------------------------------------------------------------------

type BlockIndexEntry struct {
	BlockEntry

	Handle BlockHandle
}

func (self BlockIndexEntry) String() string {
	return fmt.Sprintf("BlockIndexEntry { Shared: %v, Unshared: %v, ValueLen: %v, Key: %8s, Handle: %v}", 
						self.Shared, 
						self.Unshared,
						self.ValueLen,
						self.Key,
						self.Handle)
}