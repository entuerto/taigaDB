// Copyright 2015 The taigaDB Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package table

import (
	_ "fmt"
	"testing"
)



func TestTableOpen(t *testing.T) {
	if table, err := Open("../data/h.no-compression.sst"); err != nil {
		t.Error(err)
	} else {
		sst := table.(*ssTable)
		sst.Dump()
		sst.file.Close()
	}
}

func BenchmarkEntryIterator(b *testing.B) { 
	if table, err := Open("../data/h.no-compression.sst"); err != nil {
		b.Error(err)
	} else {
		b.StartTimer()

		sst := table.(*ssTable)
		idx := sst.BlockIndex.next

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
