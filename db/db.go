// Copyright 2015 The av-vortex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

// DB is a key/value store.
//
// It is safe to call Get and Find from concurrent goroutines. 
//
// Tu Put or Delete a key/value, you have to obtain a transaction.
type DB interface {
	// Get gets the value for the given key. It returns ErrNotFound if the DB
	// does not contain the key.
	//
	// The caller should not modify the contents of the returned slice, but
	// it is safe to modify the contents of the argument after Get returns.
	Get(key interface{}) (value interface{}, err error)

	// Find returns an iterator positioned before the first key/value pair
	// whose key is 'greater than or equal to' the given key. There may be no
	// such pair, in which case the iterator will return false on Next.
	//
	// Any error encountered will be implicitly returned via the iterator. An
	// error-iterator will yield no key/value pairs and closing that iterator
	// will return that error.
	//
	// It is safe to modify the contents of the argument after Find returns.
	Find(key interface{}) Iterator

	//
	Tx() Transaction

	// Close closes the DB. It may or may not close any underlying io.Reader
	// or io.Writer, depending on how the DB was created.
	//
	// It is not safe to close a DB until all outstanding iterators are closed.
	// It is valid to call Close multiple times. Other methods should not be
	// called after the DB has been closed.
	Close() error
}

//---------------------------------------------------------------------------------------
//
//---------------------------------------------------------------------------------------
