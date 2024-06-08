// Copyright 2023 dvonthenen/go-utilities contributors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	initlib "github.com/dvonthenen/go-utilities/diff-directory"
	diffdirectory "github.com/dvonthenen/go-utilities/diff-directory/pkg/diff-directory"
)

func printHelp() {
	fmt.Println("Usage: diff-directory -src <src> -dst <dst> [-skipsrc] [-dryrun] [-logging <level>]")
	fmt.Println("Options:")
	fmt.Println("  -src string")
	fmt.Println("    	The source directory for all music files")
	fmt.Println("  -dst string")
	fmt.Println("    	The destination directory for all music files")
	fmt.Println("  -skipsrc")
	fmt.Println("    	Skip updating the source of the diff")
	fmt.Println("  -dryrun")
	fmt.Println("    	Do a run run only... don't update/copy any files")
	fmt.Println("  -logging int")
	fmt.Println("    	Set logging level: 2 - standard (default), 7 - very verbose")

}

func main() {
	// flags
	var skipSrc bool
	flag.BoolVar(&skipSrc, "skipsrc", false, "Skip updating the source of the diff")

	var dryrun bool
	flag.BoolVar(&dryrun, "dryrun", false, "Do a run run only... don't update/copy any files")

	var srcDir string
	flag.StringVar(&srcDir, "src", "", "The source directory for all music files")

	var dstDir string
	flag.StringVar(&dstDir, "dst", "", "The destination directory for all music files")

	var logging int
	flag.IntVar(&logging, "logging", 2, "Set logging level: 2 - standard (default), 7 - very verbose")

	flag.Parse()
	// flags

	initlib.Init(initlib.DiffDirectoryInit{
		LogLevel: initlib.LogLevel(logging),
	})

	// src
	if len(srcDir) == 0 {
		fmt.Println("Provided src path is empty. Must provide a valid directory.")
		fmt.Println()
		printHelp()
		os.Exit(1)
	}

	absSrcPath, err := filepath.Abs(srcDir)
	if err != nil {
		fmt.Printf("Source filepath.Abs failed. Err: %v\n", err)
		fmt.Println()
		printHelp()
		os.Exit(1)
	}

	stat, err := os.Stat(absSrcPath)
	if err != nil {
		fmt.Printf("Invalid src=%s directory. Must provide a valid directory.\n", absSrcPath)
		fmt.Println()
		printHelp()
		os.Exit(1)
	}
	if !stat.IsDir() {
		fmt.Printf("Invalid src=%s directory. Must provide a valid directory.\n", absSrcPath)
		fmt.Println()
		printHelp()
		os.Exit(1)
	}

	//dst
	if len(dstDir) == 0 {
		fmt.Println("Provided dst path is empty. Must provide a valid directory.")
		fmt.Println()
		printHelp()
		os.Exit(1)
	}

	absDstPath, err := filepath.Abs(dstDir)
	if err != nil {
		fmt.Printf("Destination filepath.Abs failed. Err: %v\n", err)
		fmt.Println()
		printHelp()
		os.Exit(1)
	}

	stat, err = os.Stat(absDstPath)
	if err != nil {
		fmt.Printf("Invalid dst=%s directory. Must provide a valid directory.\n", absDstPath)
		fmt.Println()
		printHelp()
		os.Exit(1)
	}
	if !stat.IsDir() {
		fmt.Printf("Invalid dst=%s directory. Must provide a valid directory.\n", absDstPath)
		fmt.Println()
		printHelp()
		os.Exit(1)
	}

	// output
	fmt.Printf("logging: %d\n", logging)
	fmt.Printf("Src Path: %s\n", absSrcPath)
	fmt.Printf("Dst Path: %s\n", absDstPath)
	fmt.Printf("Skip Src: %t\n", skipSrc)
	fmt.Printf("Dry Run: %t\n", dryrun)
	fmt.Printf("\n\n")

	dist := diffdirectory.New(diffdirectory.DiffOpts{
		RootSrcPath:   absSrcPath,
		RootDstPath:   absDstPath,
		SkipSrcUpdate: skipSrc,
		DryRun:        dryrun,
	})

	err = dist.Process()
	if err == nil {
		fmt.Printf("Diff Completed!\n")
	} else {
		fmt.Printf("Process failed. Err: %v\n", err)
	}
}
