// Copyright 2015 The av-vortex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import (
	"hash/crc32"
)

var table = crc32.MakeTable(crc32.Castagnoli)

// Checksum returns the CRC-32 checksum of data
func Checksum32(b []byte) uint32 {
	c := crc32.Update(0, table, b)
	return uint32(c >> 15 | c << 17) + 0xa282ead8
}