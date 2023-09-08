// Copyright 2023 dvonthenen/go-utilities contributors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package distribute

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	klog "k8s.io/klog/v2"
)

func New(opts DistributeOpts) *Distribute {
	dist := &Distribute{
		options: opts,
	}
	return dist
}

func (d *Distribute) Process() error {
	cnt := int64(0)
	err := filepath.Walk(d.options.RootSrcPath, func(path string, info os.FileInfo, err error) error {
		klog.V(5).Infof("Path: %s\n", path)

		if err != nil {
			klog.Errorf("filepath.Walk(%s) failed. Err: %v\n", path, err)
			return err
		}

		// get path helpers
		if !strings.HasPrefix(path, d.options.RootSrcPath) {
			klog.Errorf("Path %s does not have prefix %s", path, d.options.RootSrcPath)
			return fmt.Errorf("path %s does not have prefix %s", path, d.options.RootSrcPath)
		}
		srcFile := strings.TrimPrefix(path, d.options.RootSrcPath)

		if len(srcFile) == 0 {
			klog.V(5).Infof("Skip this directory (aka self ./)")
			return nil
		}

		// handle dir
		if info.IsDir() {
			klog.V(4).Infof("%s is folder", srcFile)
			dstDir := filepath.Join(d.options.RootDstPath, srcFile)

			klog.V(4).Infof("Creating MkdirAll(%s)\n", dstDir)
			// fmt.Printf("Creating folder %s\n", dstDir)
			err := os.MkdirAll(dstDir, os.ModePerm)
			if err != nil {
				klog.Errorf("os.MkdirAll(%s) failed. Err: %v\n", dstDir, err)
				return err
			}
			return nil
		}

		// does dir exist? create it!
		dstFolder := filepath.Join(d.options.RootDstPath, strconv.FormatInt(int64(cnt)%(d.options.MaxFolders+1), 10))

		_, err = os.Stat(dstFolder)
		if err != nil {
			klog.V(4).Infof("Creating MkdirAll(%s)\n", dstFolder)
			// fmt.Printf("Creating folder %s\n", dstFolder)
			err = os.MkdirAll(dstFolder, os.ModePerm)
			if err != nil {
				klog.Errorf("os.MkdirAll(%s) failed. Err: %v\n", dstFolder, err)
				return err
			}
		}

		// handle file
		dstFile := filepath.Join(dstFolder, srcFile)

		klog.V(4).Infof("Copying %s -> %s\n", path, dstFile)
		fmt.Printf("Creating MP3: %s\n", dstFile)
		_, err = d.copy(path, dstFile)
		if err != nil {
			klog.Errorf("copy file %s failed. Err: %v\n", path, err)
			return err
		}

		cnt++
		return nil
	})

	if err != nil {
		klog.Errorf("Process failed. Err: %v\n", err)
	} else {
		klog.V(2).Infof("Distribute.Process succeeded")
	}
	return err
}

func (d *Distribute) copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		klog.Errorf("os.Stat(%s) failed. Err: %v\n", src, err)
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		klog.Errorf("os.Open(%s) failed. Err: %v\n", src, err)
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		klog.Errorf("os.Create(%s) failed. Err: %v\n", dst, err)
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	if sourceFileStat.Size() != nBytes {
		klog.Errorf("copy byte size mismatch. src: %d != dst: %d\n", sourceFileStat.Size(), nBytes)
		return nBytes, fmt.Errorf("copy byte size mismatch. src: %d != dst: %d\n", sourceFileStat.Size(), nBytes)
	}
	return nBytes, err
}
