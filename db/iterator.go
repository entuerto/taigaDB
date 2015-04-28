// Copyright 2015 The av-vortex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

// Iterator iterates over a DB's key/value pairs in key order.
type Iterator interface {
	// Is positioned at a valid node
	Valid() bool
	
	// Next moves the iterator to the next key/value pair.
	// It returns whether the iterator is exhausted.
	Next() bool

	// Key returns the key of the current key/value pair, or nil if done.
	// The caller should not modify the returned contents.
	Key() interface{}

	// Value returns the value of the current key/value pair, or nil if done.
	// The caller should not modify the returned contents.
	Value() interface{}
}
