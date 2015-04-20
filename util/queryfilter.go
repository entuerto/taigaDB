// Copyright 2015 The taigaDB Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Taken from leveldb-go (Copyright 2013 The LevelDB-Go Authors.)

package util

type QueryFilter interface {
	// Name of the query filter
	Name() string

	// KeyMayMatch returns whether the filter may contain given key. False positives
	// are possible, where it returns true for keys not in the original set.
	KeyMayMatch(key []byte) bool
}

// A Bloom Filter is an encoded set of []byte keys.
type BloomFilter []byte

func (f BloomFilter) Name() string {
	return "Bloom Filter"
}

func (f BloomFilter) KeyMayMatch(key []byte) bool {
	if len(f) < 2 {
		return false
	}

	k := f[len(f) - 1]
	if k > 30 {
		// This is reserved for potentially new encodings for short Bloom filters.
		// Consider it a match.
		return true
	}

	numBits := uint32(8 * (len(f) - 1))
	h     := hash(key)
	delta := h >> 17 | h << 15

	for j := uint8(0); j < k; j++ {
		bitPos := h % numBits
		if f[bitPos / 8] & (1 << (bitPos % 8)) == 0 {
			return false
		}
		h += delta
	}
	return true
}

// NewBloomFilter returns a new Bloom filter that encodes a set of []byte keys with
// the given number of bits per key. The returned Filter may be a sub-slice of
// buf[:cap(buf)] if it is large enough, otherwise the Filter will be allocated
// separately.
//
// A good value for bits per key is 10, which yields a filter with ~ 1% false positive rate.
func NewBloomFilter(buf []byte, keys [][]byte, bitsPerKey int) QueryFilter {
	if bitsPerKey < 0 {
		bitsPerKey = 0
	}
	// 0.69 is approximately ln(2).
	k := uint32(float64(bitsPerKey) * 0.69)
	if k < 1 {
		k = 1
	}
	if k > 30 {
		k = 30
	}

	numBits := len(keys) * int(bitsPerKey)
	// For small n, we can see a very high false positive rate. Fix it
	// by enforcing a minimum bloom filter length.
	if numBits < 64 {
		numBits = 64
	}
	numBytes := (numBits + 7) / 8
	numBits = numBytes * 8

	if numBytes + 1 <= cap(buf) {
		buf = buf[:numBytes + 1]
		for i := range buf {
			buf[i] = 0
		}
	} else {
		buf = make([]byte, numBytes + 1)
	}

	for _, key := range keys {
		h := hash(key)
		delta := h >> 17 | h << 15
		for j := uint32(0); j < k; j++ {
			bitPos := h % uint32(numBits)
			buf[bitPos / 8] |= 1 << (bitPos % 8)
			h += delta
		}
	}
	buf[numBytes] = uint8(k)
	return BloomFilter(buf)	
}

// hash implements a hashing algorithm similar to the Murmur hash.
func hash(b []byte) uint32 {
	const (
		seed = 0xbc9f1d34
		m    = 0xc6a4a793
	)

	h := uint32(seed) ^ uint32(len(b) * m)
	for ; len(b) >= 4; b = b[4:] {
		h += uint32(b[0]) | uint32(b[1]) << 8 | uint32(b[2]) << 16 | uint32(b[3]) << 24
		h *= m
		h ^= h >> 16
	}
	switch len(b) {
	case 3:
		h += uint32(b[2]) << 16
		fallthrough
	case 2:
		h += uint32(b[1]) << 8
		fallthrough
	case 1:
		h += uint32(b[0])
		h *= m
		h ^= h >> 24
	}
	return h
}