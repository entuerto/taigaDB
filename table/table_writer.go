// Copyright 2015 The taigaDB Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package table

import (
	_ "encoding/binary"
	_ "fmt"
	"os"

	_ "github.com/entuerto/taigaDB/util"
	_ "code.google.com/p/snappy-go/snappy"
)

func NewWriter(filename string, opt *Options) (TableWriter, error) {
	var table = &ssTableWriter{
		options: opt,
	}

	// Read only
	file, err := os.Create(filename) 
	if err != nil {
		return nil, err
	}
	table.file = file

	if table.options == nil {
		table.options = DefaultOptions()
	}

	return table, nil
}

//---------------------------------------------------------------------------------------
// Writer for Sorted String Table
//---------------------------------------------------------------------------------------

type ssTableWriter struct {
	file *os.File

	options *Options
}

func (self ssTableWriter) Write(key, value Slice) error {
	
	return nil
}

func (self ssTableWriter) Close() error {
	return self.file.Close()
}
