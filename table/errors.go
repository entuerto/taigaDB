// Copyright 2015 The taigaDB Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package table

import (
	"errors"
)

var (
	BlockReadCorruptionErr = errors.New("Table.Block: Block read corruption")
	BlockCRC32CorruptionErr = errors.New("Table.Block: Block checksum mismatch")

	DecodeSmallBufferErr = errors.New("Decode: Buffer to small")
	DecodeNot64bitsErr   = errors.New("Decode: Value is not 64bits")

	EncodeBlockHandleBufferErr = errors.New("BlockHandle.Encode: Buffer too small")
	EncodeFooterBufferErr      = errors.New("Footer.Encode: Buffer too small")
	
	TableMagicNumberErr = errors.New("Table: Wrong table format")
	TableBlockCompressionErr = errors.New("Table.Block: Wrong compression format")

	NotFoundErr = errors.New("Table: Value was not found")
	NotImplementedErr = errors.New("Table: Not implemented")
)
