// Copyright 2015 The taigaDB Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package table

import (
	_ "fmt"
	"testing"
)



func TestTableOpen(t *testing.T) {
	table, err := NewReader("../data/h.no-compression.sst", DefaultOptions())
	defer table.Close()

	if err != nil {
		t.Error(err)
	} 

	//sst := table.(*ssTable)
	//sst.Dump()
}

func TestTableLookup(t *testing.T) {
	table, err := NewReader("../data/h.no-compression.sst", DefaultOptions())
	defer table.Close()

	if err != nil {
		t.Error(err)
	} 

	value, err := table.Read(Slice("school")) 
	if err != nil {
		t.Error(err)
	}
	if string(value) != "1" {
		t.Errorf("Looking for school, got %s", value)
	}

}

func BenchmarkEntryIterator(b *testing.B) { 
	table, err := NewReader("../data/h.no-compression.sst", DefaultOptions())
	defer table.Close()

	if err != nil {
		b.Error(err)
		return
	} 

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
