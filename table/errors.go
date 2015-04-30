// Copyright 2015 The taigaDB Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package table

import (
	"errors"
)

var (
	ErrBlockReadCorruption = errors.New("Table.Block: Block read corruption")
	ErrBlockCRC32Corruption = errors.New("Table.Block: Block checksum mismatch")

	ErrDecodeSmallBuffer = errors.New("Decode: Buffer to small")
	ErrDecodeNot64bits   = errors.New("Decode: Value is not 64bits")

	ErrEncodeBlockHandleBuffer = errors.New("BlockHandle.Encode: Buffer too small")
	ErrEncodeFooterBuffer      = errors.New("Footer.Encode: Buffer too small")
	
	ErrTableMagicNumber = errors.New("Table: Wrong table format")
	ErrTableBlockCompression = errors.New("Table.Block: Wrong compression format")

	ErrNotFound = errors.New("Table: Value was not found")
	ErrNotImplemented = errors.New("Table: Not implemented")
)
