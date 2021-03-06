// Copyright 2015 The taigaDB Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Table

A sorted string table stores a sequence of entries sorted by key.

The sequence of key/value pairs in the file are stored in sorted
order and partitioned into a sequence of data blocks.  These blocks
come one after another at the beginning of the file.  

Table Structure:
                                                         
    +------------------+
    | Data block 1     |    Each block followed by a 5-bytes trailer contains 
    +------------------+    compression type and checksum.
    | ...              |
    +------------------+
    | Data block n     | 
    +------------------+
    | Filter block     | 
    +------------------+    
    | Meta index block | 
    +------------------+    
    | Index block      |
    +------------------+        +-------------------------+
    | Footer           |   ->   | Meta index block handle |  Table footer: 48 bytes long
    +------------------+        +-------------------------+
                                | Index block handle      |
                                +-------------------------+
                                | [Padding up to 40 bytes]|
                                +-------------------------+
                                |  magic (8-bytes)        |  The magic are first 64-bit of SHA-1 sum 
                                +-------------------------+  of "http://code.google.com/p/leveldb/"

Blocks

Blocks have one or many key/value entries followed by a block trailer structure.

When keys are stored, we drop the prefix shared with the previous string. This helps 
reduce the space requirement. Furthermore, once every K keys, we do not apply the prefix
compression and store the entire key. This is called a "restart point". The tail end 
of the block stores the offsets of all of the restart points, and can be used to do 
a binary search when looking for a particular key. Values are stored without compression
immediately following the corresponding key.

Index Block:

The index block contains one entry per data block, where the key is a string >= last key 
in that data block and before the first key in the successive data block.  The value is the
BlockHandle for the data block.

Data Block Structure:

    +---------------+  -> Restart point
    | Block entry 1 |
    +---------------+
    | Block entry 2 |
    +---------------+  -> Restart point (depends on restart interval)
    | ...           | 
    +---------------+
    | Block entry n | 
    +---------------+
    | Trailer       |  -> Block trailer:
    +---------------+       +---------------------------------+   
                            | Restarts (uint32[num_restarts]) |  restarts[i] contains the offset within the block 
                            +---------------------------------+  of the nth restart point.  
                            | Number of restarts (uint32)     |
                            +---------------------------------+

Block Entry Structure:

    +----------------------+
    | Shared (varint32)    |
    +----------------------+  -> Key length is the shared and not shared values.
    | Unshared (varint32)  |
    +----------------------+
    | Value len (varint32) |
    +----------------------+
    | Key delta ([]byte)   |  -> Slice of unshared bytes
    +----------------------+
    | Value ([]byte)       |  -> Slice of value length
    +----------------------+

    Shared == 0 for restart point.

Block in SSTable file: 

An extra block "footer" is saved for each data block that indicates if the block is saved 
compressed or uncompressed in the file and a CRC checksum for error protection.

The checksum is a CRC-32 computed using Castagnoli's polynomial. Compression type 
also included in the checksum.

    +---------------------------+
    | Data Block Structure      |
    +---------------------------+
    | Compression type (1-byte) |   
    +---------------------------+   
    | Checksum (4-byte)         |
    +---------------------------+

*/
package table                           
