// Copyright 2023 dvonthenen/go-utilities contributors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package diff

import "io/fs"

type DIRECTION int

const (
	UNKNOWN_DIRECTION DIRECTION = iota
	DIRECTION_SRC_TO_DST
	DIRECTION_DST_TO_SRC
)

type DiffOpts struct {
	RootSrcPath   string
	RootDstPath   string
	SkipSrcUpdate bool
	DryRun        bool
}

type Diff struct {
	options DiffOpts
}

type DiffFile struct {
	Path    string
	RelPath string
	Attr    *fs.FileInfo
	Hash    string // crypto/sha256 calculated only if attr mod is different
}

type DiffCompare struct {
	SrcFile   *DiffFile
	DstFile   *DiffFile
	Direction DIRECTION
}
