// Copyright 2015 The taigaDB Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package table

import (
	_ "fmt"
	"testing"
)



func TestTableOpen(t *testing.T) {
	if _, err := Open("../data/h.no-compression.sst"); err != nil {
		t.Error(err)
	} 
}

func TestTableLookup(t *testing.T) {
	table, err := Open("../data/h.no-compression.sst")

	if err != nil {
		t.Error(err)
	} 

	value, err := table.Lookup(Slice("school")) 
	if err != nil {
		t.Error(err)
	}
	if string(value) != "1" {
		t.Errorf("Looking for school, got %s", value)
	}

}

func BenchmarkEntryIterator(b *testing.B) { 
	if table, err := Open("../data/h.no-compression.sst"); err != nil {
		b.Error(err)
	} else {
		b.StartTimer()

		sst := table.(*ssTable)
		idx := sst.BlockIndex[0]

		for i := 0; i < b.N; i++ {
			if block, err := sst.readBlock(&idx.Handle); err == nil {
				iter := NewEntryIterator(block) 

				for _, ok := iter.Next(); ok; { 
					
					_, ok = iter.Next()
				}
			}
		}
	}
}
