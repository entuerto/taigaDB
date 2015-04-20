// Copyright 2015 The av-vortex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

type Transaction interface {

	// Set sets the value for the given key. It overwrites any previous value
	// for that key; a DB is not a multi-map.
	Put(key, value interface{}) error

	// Delete deletes the value for the given key. It returns ErrNotFound if
	// the DB does not contain the key.
	Delete(key interface{}) error


	// Commit commits the transaction to the database. 
	Commit() error
}