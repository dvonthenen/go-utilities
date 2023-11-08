// Copyright 2023 dvonthenen/go-utilities contributors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	initlib "github.com/dvonthenen/go-utilities/how-alike"
	howalike "github.com/dvonthenen/go-utilities/how-alike/pkg/howalike"
)

func main() {
	var checkFile string
	flag.StringVar(&checkFile, "check", "", "The file to check against the actual file")

	var actualFile string
	flag.StringVar(&actualFile, "actual", "", "The actual transcription file")

	flag.Parse()

	// init lib
	initlib.Init(initlib.FileDistributeInit{
		LogLevel: initlib.LogLevelStandard,
	})

	if len(checkFile) == 0 {
		fmt.Println("The file to check is empty. Must provide a valid file.")
		os.Exit(1)
	}
	if len(actualFile) == 0 {
		fmt.Println("The transcription file is empty. Must provide a valid file.")
		os.Exit(1)
	}

	// src
	absCheckFile, err := filepath.Abs(checkFile)
	if err != nil {
		fmt.Println("Source filepath.Abs failed. Err: %v\n", err)
		os.Exit(1)
	}

	stat, err := os.Stat(absCheckFile)
	if err != nil {
		fmt.Println("Invalid src=%s directory. Must provide a valid directory.", absCheckFile)
		os.Exit(1)
	}
	if stat.IsDir() {
		fmt.Println("Invalid src=%s is directory. Must provide a valid file.", absCheckFile)
		os.Exit(1)
	}
	fmt.Printf("File to Check: %s\n", absCheckFile)

	//dst
	absActualFile, err := filepath.Abs(actualFile)
	if err != nil {
		fmt.Println("Source filepath.Abs failed. Err: %v\n", err)
		os.Exit(1)
	}

	stat, err = os.Stat(absActualFile)
	if err != nil {
		fmt.Println("Invalid src=%s directory. Must provide a valid directory.", absActualFile)
		os.Exit(1)
	}
	if stat.IsDir() {
		fmt.Println("Invalid src=%s is directory. Must provide a valid file.", absActualFile)
		os.Exit(1)
	}
	fmt.Printf("File to Check: %s\n", absActualFile)

	chk := howalike.New(howalike.HowAlikeOptions{
		CheckFile:  absCheckFile,
		ActualFile: absActualFile,
	})

	err = chk.Process()
	if err == nil {
		fmt.Printf("How Alike Completed!\n")
	} else {
		fmt.Printf("Process failed. Err: %v\n", err)
	}
}
