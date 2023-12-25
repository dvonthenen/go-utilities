// Copyright 2023 dvonthenen/go-utilities contributors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package diff

import (
	"errors"
)

var (
	// ErrDiffSizeCopied the input and output copy size doesnt match
	ErrDiffSizeCopied = errors.New("the input and output copy size doesnt match")

	// ErrUnknownDirection unknown direction to copy file (src -> dst OR dst -> src)
	ErrUnknownDirection = errors.New("unknown direction to copy file (src -> dst OR dst -> src)")
)
