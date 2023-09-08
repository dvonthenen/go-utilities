// Copyright 2023 dvonthenen/go-utilities contributors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package distribute

type DistributeOpts struct {
	RootSrcPath string
	RootDstPath string
	MaxFolders  int64
}

type Distribute struct {
	options DistributeOpts
}
