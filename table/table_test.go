// Copyright 2015 The taigaDB Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package table

import (
	"fmt"
	"testing"
)



func TestTableOpen(t *testing.T) {
	if table, err := Open("../data/h.no-compression.sst"); err != nil {
		t.Error(err)
	} else {
		fmt.Println(table)

		sst := table.(*ssTable)
		sst.file.Close()
	}
}
/*
func TestTableReadBlock(t *testing.T) {
	if table, err := Open("../data/h.no-compression.sst"); err != nil {
		t.Error(err)
	} else {
		sst := table.(*ssTable)

		block, _ := sst.readBlock(sst.BlockIndexHandle)
		fmt.Println(string(block))

		sst.file.Close()
	}
}
*/