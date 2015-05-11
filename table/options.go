// Copyright 2015 The taigaDB Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package table

import (
	"github.com/entuerto/taigaDB/util"
)

// Compression algorithm for block entries.
type Compression int

const (
	NoCompression      Compression = iota
	SnappyCompression
)

// Options holds the parameters for the table implementation.
type Options struct {
	// Number of keys between restart points for delta encoding of keys.
	//
	// The default value is 16.
	BlockRestartInterval int

	// Approximate size of user data packed per block.  Note that the
	// block size specified here corresponds to uncompressed data.  The
	// actual size of the unit read from disk may be smaller if
 	// compression is enabled.  
 	//
	// The default value is 4096 (4K).
	BlockSize int

	// Used to define the order of keys in the table. keys: a 'less
	// than' relationship. The same comparison algorithm must be used for reads
	// and writes.
	//
	// The default value uses the same ordering as bytes.Compare.
	Comparator util.Comparator

	// Compress blocks using the specified compression algorithm.
	//
	// The default value is no compression.
	Compression Compression

	// Whether to verify the per-block checksums in a table.
	//
	// The default value is false.
	VerifyChecksums bool
}

func DefaultOptions() *Options {
	return &Options{
		BlockRestartInterval: 16,
		BlockSize: 4096,
		Comparator: util.BytewiseComparator{},
		Compression: NoCompression,
		VerifyChecksums: false,
	}
}