// Copyright 2015 The av-vortex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import (
	"bytes"
)

// Provides a total order across slices that are used as keys in an sstable or 
// a database.  A Comparator implementation must be thread-safe since it may 
// be invoked concurrently  from multiple threads.
type Comparator interface {
	// Three-way comparison.  Returns value:
	//   < 0 iff "a" < "b",
	//   == 0 iff "a" == "b",
	//   > 0 iff "a" > "b"
	Compare(a, b []byte) int

	// Name returns the name of the comparator.
	//
	// The Level-DB on-disk format stores the comparator name, and opening a
	// database with a different comparator from the one it was created with
	// will result in an error.
	Name() string
}


type BytewiseComparator struct{}

func (BytewiseComparator) Name() string {
	// This string is part of the C++ Level-DB implementation's default file format,
	// and should not be changed.
	return "leveldb.BytewiseComparator"
}

func (BytewiseComparator) Compare(a, b []byte) int {
	return bytes.Compare(a, b)
}

// Returns the largest pos such that a[:pos] equals b[:pos].
func SharedPrefix(a, b []byte) int {
	i, n := 0, len(a)
	if n > len(b) {
		n = len(b)
	}
	for i < n && a[i] == b[i] {
		i++
	}
	return i
}