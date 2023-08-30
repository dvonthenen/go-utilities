// Copyright 2023 dvonthenen/go-utilities contributors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"

	klog "k8s.io/klog/v2"
)

const (
	// DefaultDirAppend
	DefaultDirAppend string = "_NEW"
)

func main() {
	var init FileDistributeInit
	Init(&init)

	klog.V(2).Infof("FileDistribute Init")
	fmt.Printf("srcDir: %s\n", init.SrcDir)
}
