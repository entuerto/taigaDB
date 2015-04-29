// Copyright 2015 The taigaDB Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package table

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/entuerto/taigaDB/util"
)

const (
	NoCompression     = 0
	SnappyCompression = 1
)

type Slice []byte

// A Table is a sorted map from strings to strings.
type Table interface {
	// Look for a key/value in the table.
	Lookup(key Slice) (Slice, error)

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

	// Read the meta block
	if metaBlock, err := table.readBlock(table.MetaIndexHandle); err != nil {
		return nil, err
	} else {
		table.MetaIndex = decodeIndexEntries(metaBlock)
	}

	// Read the index block
	if indexBlock, err := table.readBlock(table.BlockIndexHandle); err != nil {
		return nil, err
	} else {
		table.BlockIndex = decodeIndexEntries(indexBlock)
	}

	return table, nil
}

//---------------------------------------------------------------------------------------
// Sorted String Table
//---------------------------------------------------------------------------------------

type ssTable struct {
	file *os.File

	// Handles to the specified file location
	MetaIndexHandle  *BlockHandle
	BlockIndexHandle *BlockHandle

	MetaIndex  IndexSlice
	BlockIndex IndexSlice
}

func (self ssTable) String() string {
	return fmt.Sprintf("SSTable { MetaIndex: %v, BlockIndex: %v }", self.MetaIndexHandle, self.BlockIndexHandle)
}

func (self *ssTable) Iterator() Iterator {
	return &ssTableIterator{
		sst: self,
		idx: self.BlockIndex,
	}
}

func (self *ssTable) ApproximateOffsetOf(key interface{}) uint64 {
	return uint64(0)
}

func (self ssTable) Lookup(key Slice) (Slice, error) {
	//
	i := self.BlockIndex.Search(key)
	if i == self.BlockIndex.Len() {
		return nil, NotFoundErr
	}
	
	blockIdx := self.BlockIndex[i]
	block, err := self.readBlock(&blockIdx.Handle)
	if err != nil {
		return nil, err
	}

	entry := block.Search(key)
	if entry != nil {
		return entry.Value, nil
	}	
	return nil, NotFoundErr
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

func (self ssTable) Dump() {
	fmt.Println()
	fmt.Println(self)
	fmt.Println()
	fmt.Println("** Block Index **")
	fmt.Println()
	for _, i := range self.BlockIndex {
		fmt.Println(i)
	}
	fmt.Println()
	fmt.Println("** Metadata Index **")
	fmt.Println()
	for _, i := range self.MetaIndex {
		fmt.Println(i)
	}
/*
	fmt.Println()
	fmt.Println("** First Data Block **")
	fmt.Println()
	idx := self.BlockIndex[0]
	if block, err := self.readBlock(&idx.Handle); err == nil {
		block.Dump()
	}
*/
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

