// Copyright 2015 The av-vortex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"errors"
)

package db

var (
	ErrKeyNotFound = errors.New("db: key not found")
)
