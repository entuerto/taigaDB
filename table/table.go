// Copyright 2015 The taigaDB Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package table

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	"github.com/entuerto/taigaDB/util"
)

const (
	NoCompression     = 0
	SnappyCompression = 1
)

var ( 
	BlockReadCorruptionErr = errors.New("Table.Block: Block read corruption")
	BlockCRC32CorruptionErr = errors.New("Table.Block: Block checksum mismatch")
	TableMagicNumberErr = errors.New("Table: Wrong table format")
	TableBlockCompressionErr = errors.New("Table.Block: Wrong compression format")

	NotImplementedErr = errors.New("Table: Not implemented")
)

type Slice []byte

type Iterator interface {}

// A Table is a sorted map from strings to strings.
type Table interface {
	// Returns an iterator over the table contents.
	Iterator() Iterator

	// Given a key, return an approximate byte offset in the file where
	// the data for that key begins (or would begin if the key were
	// present in the file). 
	ApproximateOffsetOf(key interface{}) uint64 
}

func Open(filename string) (Table, error) {
	var table = new(ssTable)

	// Read only
	file, err := os.Open(filename) 
	if err != nil {
		return nil, err
	}
	table.file = file

	if err = table.readFooter(); err != nil {
		return nil, err
	}

	if err = table.readFilter(); err != nil {
		return nil, err
	}

	if err = table.readMeta(); err != nil {
		return nil, err
	}

	// Read the index block
	if err = table.readIndex(); err != nil {
		return nil, err
	}

	return table, nil
}

//---------------------------------------------------------------------------------------
//
//---------------------------------------------------------------------------------------

type ssTable struct {
	file *os.File

	// Handles to the specified file location
	MetaIndexHandle  *BlockHandle
	BlockIndexHandle *BlockHandle
}

func (self ssTable) String() string {
	return fmt.Sprintf("SSTable { MetaIndex: %v, BlockIndex: %v }", self.MetaIndexHandle, self.BlockIndexHandle)
}

func (self *ssTable) Iterator() Iterator {
	return nil
}

func (self *ssTable) ApproximateOffsetOf(key interface{}) uint64 {
	return uint64(0)
}

func (self *ssTable) readFooter() error {
	var offset int64

	if fi, err := self.file.Stat(); err != nil {
		return err
	} else {
		offset = fi.Size() - FooterEncodedLength
	}

	var buffer [FooterEncodedLength]byte

	if n, err := self.file.ReadAt(buffer[0:], offset); err != nil || n != FooterEncodedLength {
		return err
	}

	var footer Footer

	if _, err := footer.Decode(buffer[0:]); err != nil {
		return err
	}

	if footer.MagicNumber != TableMagicNumber {
		return TableMagicNumberErr
	}

	self.MetaIndexHandle  = footer.MetaIndexHandle
	self.BlockIndexHandle = footer.BlockIndexHandle

	return nil
}

func (self *ssTable) readIndex() error {

	indexBlock, err := self.readBlock(self.BlockIndexHandle)

	if  err != nil {
		return err
	} 

	fmt.Printf("Restarts: %#v \n", indexBlock.NumberOfRestarts())

	var be BlockIndexEntry
	var n0, n1, n2 int

	for i := 0; i < int(indexBlock.NumberOfRestarts()); i++ {
		be.Shared,   n0 = binary.Uvarint(indexBlock)
		be.Unshared, n1 = binary.Uvarint(indexBlock[n0:])
		be.ValueLen, n2 = binary.Uvarint(indexBlock[n0 + n1:])
	
		n := n0 + n1 + n2
		be.Key   = Slice(append(be.Key[:int(be.Shared)], indexBlock[n:n + int(be.Unshared)]...))
		be.Value = Slice(indexBlock[n + int(be.Unshared):n + int(be.Unshared + be.ValueLen)])
		indexBlock = indexBlock[n + int(be.Unshared + be.ValueLen):]
	
		be.Handle.Decode(be.Value)
	
		fmt.Println(be)
	}

	return nil
}

func (self *ssTable) readMeta() error {

	//readBlock(self.MetaIndexHandle)
	return nil
}

func (self *ssTable) readFilter() error {
	return nil
}

func (self *ssTable) readBlock(bh *BlockHandle) (Block, error) {
	var buffer = make([]byte, bh.Size + BlockTrailerSize)

	n, err := self.file.ReadAt(buffer[:], int64(bh.Offset))
	if err != nil {
		return nil, err
	}
	if n != int(bh.Size + BlockTrailerSize) {
		return nil, BlockReadCorruptionErr
	}

	// Checksum from block trailer
	checksum1 := binary.LittleEndian.Uint32(buffer[bh.Size + 1:])
	// Checksum calculated from block buffer
	checksum2 := util.Checksum32(buffer[:bh.Size + 1])

	if checksum1 != checksum2 {
		return nil, BlockCRC32CorruptionErr
	}

	switch buffer[bh.Size] {
	case NoCompression:
		return buffer[:bh.Size], nil
	case SnappyCompression:
		/*
		b, err := snappy.Decode(nil, b[:bh.length])
		if err != nil {
			return nil, err
		}
		return b, nil
		*/
		return nil, NotImplementedErr
	}

	return nil, TableBlockCompressionErr
}

