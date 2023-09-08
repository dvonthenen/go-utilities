// Copyright 2023 dvonthenen/go-utilities contributors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	initlib "github.com/dvonthenen/go-utilities/file-distribute"
	distribute "github.com/dvonthenen/go-utilities/file-distribute/pkg/distribute"
)

const (
	// DefaultDirAppend destination directory of new file structure/organization
	DefaultDirAppend string = "_NEW"

	// MaxNumOfFolders per directory
	MaxNumOfFolders int64 = 255
)

func main() {
	// flags
	maxFolders := *flag.Int64("maxFolders", 0, "The maximum folders")

	var srcDir string
	flag.StringVar(&srcDir, "src", "", "The source directory for all music files")

	var dstDir string
	flag.StringVar(&dstDir, "dst", "", "The destination directory for all music files")

	flag.Parse()
	// flags

	initlib.Init(initlib.FileDistributeInit{
		LogLevel: initlib.LogLevelStandard,
	})

	if len(srcDir) == 0 {
		fmt.Println("Provided src path is empty. Must provide a valid directory.")
		os.Exit(1)
	}

	// src
	absSrcPath, err := filepath.Abs(srcDir)
	if err != nil {
		fmt.Println("Source filepath.Abs failed. Err: %v\n", err)
		os.Exit(1)
	}

	stat, err := os.Stat(absSrcPath)
	if err != nil {
		fmt.Println("Invalid src=%s directory. Must provide a valid directory.", absSrcPath)
		os.Exit(1)
	}
	if !stat.IsDir() {
		fmt.Println("Invalid src=%s directory. Must provide a valid directory.", absSrcPath)
		os.Exit(1)
	}
	fmt.Printf("Src Path: %s\n", absSrcPath)

	//dst
	var absDstPath string
	if len(dstDir) == 0 {
		absDstPath = fmt.Sprintf("%s%s", absSrcPath, DefaultDirAppend)
	} else {
		absDstPath, err := filepath.Abs(srcDir)
		if err != nil {
			fmt.Println("Destination filepath.Abs failed. Err: %v\n", err)
			os.Exit(1)
		}

		err = os.MkdirAll(absDstPath, 0755)
		if err != nil {
			fmt.Println("MkdirAll(%s) failed. Err: %v\n", absDstPath, err)
			os.Exit(1)
		}
	}
	fmt.Printf("Dst Path: %s\n", absDstPath)

	if maxFolders == 0 {
		maxFolders = MaxNumOfFolders
	}

	dist := distribute.New(distribute.DistributeOpts{
		RootSrcPath: absSrcPath,
		RootDstPath: absDstPath,
		MaxFolders:  maxFolders,
	})

	err = dist.Process()
	if err == nil {
		fmt.Printf("Distribute Completed!\n")
	} else {
		fmt.Printf("Process failed. Err: %v\n", err)
	}
}
