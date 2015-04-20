// Copyright 2015 The av-vortex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package table

import (
	"bytes"
	"testing"
)


var (
	// Offset: 10, Size: 20
	testBlockData1 = []byte{10, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	// Offset: 5000, Size: 8000
    testBlockData2 = []byte{136, 39, 192, 62, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
)

func TestBlockHandleEncode(t *testing.T) {
	bh := NewHandle(10, 20)

	buf := make([]byte, MaxEncodedLength)

	if s, err := bh.Encode(buf); err != nil {
		t.Error(err)
	} else {
		if !bytes.Equal(buf, testBlockData1) {
			t.Error("Should be {10, 20}, got: ", buf)
		}
		if s != 2 {
			t.Error("Should be 2, got: ", s)
		}
	}

	bh.Offset = 5000
	bh.Size   = 8000

	if s, err := bh.Encode(buf); err != nil {
		t.Error(err)
	} else {
		if !bytes.Equal(buf, testBlockData2) {
			t.Error("Should be {136, 39, 192, 62}, got: ", buf)
		}
		if s != 4 {
			t.Error("Should be 4, got: ", 4)
		}
	}
}

func TestBlockHandleDecode(t *testing.T) {
	var bh BlockHandle

	if _, err := bh.Decode(testBlockData1); err == nil {
		if bh.Offset != 10 || bh.Size != 20 {
			t.Error("Should be {Offset: 10, Size: 20}, got: ", bh)
		}
	}

	if _, err := bh.Decode(testBlockData2); err == nil {
		if bh.Offset != 5000 || bh.Size != 8000 {
			t.Error("Should be {Offset: 5000, Size: 8000}, got: ", bh)
		}
	}

}

var testFooterData = []byte{10, 20, 30, 40, 0, 0, 0, 0, 0, 
		                     0,  0,  0,  0, 0, 0, 0, 0, 0, 
		                     0,  0,  0,  0, 0, 0, 0, 0, 0,
		                     0,  0,  0,  0, 0, 0, 0, 0, 0, 
		                     0,  0,  0,  0, 87, 251, 128, 139, 36, 117, 71, 219}

func TestFooterEncode(t *testing.T) {
	footer := NewFooter(NewHandle(10, 20), NewHandle(30, 40))

	buf := make([]byte, FooterEncodedLength)

	if s, err := footer.Encode(buf); err != nil {
		t.Error(err)
	} else {
		if !bytes.Equal(buf, testFooterData) {
			t.Error("Should be ", testFooterData, ", got: ", buf)
		}
		if s != FooterEncodedLength {
			t.Error("Should be 48, got: ", FooterEncodedLength)
		}
	}

}

func TestFooterDecode(t *testing.T) {
	var footer Footer

	if _, err := footer.Decode(testFooterData); err == nil {

		if footer.MetaIndexHandle.Offset != 10 || footer.MetaIndexHandle.Size != 20 {
			t.Error("Should be {Offset: 10, Size: 20}, got: ", footer.MetaIndexHandle)
		}

		if footer.BlockIndexHandle.Offset != 30 || footer.BlockIndexHandle.Size != 40 {
			t.Error("Should be {Offset: 30, Size: 40}, got: ", footer.BlockIndexHandle)
		}

		if footer.MagicNumber != TableMagicNumber {
			t.Errorf("Should be 0xdb4775248b80fb57, got: 0x%x", footer.MagicNumber)
		}
	}
}

func BenchmarkEncodeBlockHandle(b *testing.B) {
	b.StopTimer()

	var bh = NewHandle(5000, 8000)

	buf := make([]byte, MaxEncodedLength)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if _, err := bh.Encode(buf); err != nil {
			b.Error(err)
		} else if !bytes.Equal(buf, testBlockData2) {
			b.Error("Should be {136, 39, 192, 62}, got: ", buf)
		}
	}
}

func BenchmarkDecodeBlockHandle(b *testing.B) {
	b.StopTimer()

	var bh BlockHandle

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if _, err := bh.Decode(testBlockData2); err != nil {
			b.Error(err)
		} else if bh.Offset != 5000 || bh.Size != 8000 {
			b.Error("Should be {Offset: 5000, Size: 8000}, got: ", bh)
		}
	}
}

func BenchmarkEncodeFooter(b *testing.B) {
	b.StopTimer()

	footer := NewFooter(NewHandle(10, 20), NewHandle(30, 40))

	buf := make([]byte, FooterEncodedLength)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if _, err := footer.Encode(buf); err != nil {
			b.Error(err)
		} else if !bytes.Equal(buf, testFooterData) {
			b.Error("Should be ", testFooterData, ", got: ", buf)
		}
	}
}

func BenchmarkDecodeFooter(b *testing.B) {
	b.StopTimer()

	var footer Footer

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if _, err := footer.Decode(testFooterData); err == nil {
	
			if footer.MetaIndexHandle.Offset != 10 || footer.MetaIndexHandle.Size != 20 {
				b.Error("Should be {Offset: 10, Size: 20}, got: ", footer.MetaIndexHandle)
			}
	
			if footer.BlockIndexHandle.Offset != 30 || footer.BlockIndexHandle.Size != 40 {
				b.Error("Should be {Offset: 30, Size: 40}, got: ", footer.BlockIndexHandle)
			}
	
			if footer.MagicNumber != TableMagicNumber {
				b.Errorf("Should be 0xdb4775248b80fb57, got: 0x%x", footer.MagicNumber)
			}
		}
	}
}
