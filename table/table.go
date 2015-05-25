// Copyright 2015 The taigaDB Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package table

type Slice []byte

type Entry struct {
	Key   Slice
	Value Slice
}

// A Table is a sorted map from strings to strings.
type Table interface {
	Closer

	// Returns an iterator over the table contents.
	Iterator() Iterator

	// Given a key, return an approximate byte offset in the file where
	// the data for that key begins (or would begin if the key were
	// present in the file). 
	ApproximateOffsetOf(key Slice) uint64 
}

type TableReader interface {
	Table
	Reader
}

type TableWriter interface {
	Closer
	Writer
}

// Closer is the interface that wraps the basic Close method.
type Closer interface {
	Close() error
}

// Reader is the interface that wraps the basic Read method.
type Reader interface {
	// Look for a key/value in the table.
	Read(key Slice) (Slice, error)
}

// Writer is the interface that wraps the basic Write method.
type Writer interface {
	// Write a key/value to the table
	Write(key, value Slice) error
}

type Builder interface {
	Add(key, value Slice) error
	Flush() error
	Finish() error
}