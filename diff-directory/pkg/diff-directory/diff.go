// Copyright 2023 dvonthenen/go-utilities contributors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package diff

import (
	sha256 "crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	klog "k8s.io/klog/v2"
)

func New(opts DiffOpts) *Diff {
	dist := &Diff{
		options: opts,
	}
	return dist
}

func (d *Diff) Process() error {
	diff := make([]*DiffCompare, 0)

	err := d.fileComparison(&diff)
	if err != nil {
		klog.Errorf("fileComparison failed. Err: %v\n", err)
		return err
	}

	err = d.resolveDifferences(&diff)
	if err != nil {
		klog.Errorf("resolveDifferences failed. Err: %v\n", err)
		return err
	}

	return nil
}

func (d *Diff) fileComparison(diff *[]*DiffCompare) error {
	srcPath := d.options.RootSrcPath
	lenSrc := len(srcPath)
	dstPath := d.options.RootDstPath
	lenDst := len(dstPath)
	klog.V(4).Infof("srcPath: %s\n", srcPath)
	klog.V(4).Infof("dstPath: %s\n", dstPath)

	srcMap := make(map[string]*DiffFile, 0)
	dstMap := make(map[string]*DiffFile, 0)

	err := filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		filename := filepath.Base(path)
		klog.V(6).Infof("[SRC] path: %s\n", path)
		klog.V(6).Infof("[SRC] filename: %s\n", filename)
		if err != nil {
			klog.Errorf("filepath.Walk Init. Err: %v\n", err)
			return err
		}
		if strings.EqualFold(srcPath, path) || strings.EqualFold(filename, ".") || strings.EqualFold(filename, "..") {
			klog.V(6).Infof("[SRC] filepath.Walk(%s) skip . and ..\n", path)
			return nil
		}

		if info.IsDir() {
			klog.V(4).Infof("IsDir\n")
			return nil
		}
		newRel := path[lenSrc+1:]
		klog.V(6).Infof("newRel: %s\n", newRel)
		srcMap[newRel] = &DiffFile{
			Path:    path,
			RelPath: newRel,
			Attr:    &info,
		}
		return nil
	})
	if err != nil {
		klog.Errorf("filepath.Walk(%s) Err: %v\n", srcPath, err)
		return err
	}

	err = filepath.Walk(dstPath, func(path string, info os.FileInfo, err error) error {
		filename := filepath.Base(path)
		klog.V(6).Infof("[DST] path: %s\n", path)
		klog.V(6).Infof("[DST] filename: %s\n", filename)
		if err != nil {
			klog.Errorf("filepath.Walk Init. Err: %v\n", err)
			return err
		}
		if strings.EqualFold(dstPath, path) || strings.EqualFold(filename, ".") || strings.EqualFold(filename, "..") {
			klog.V(6).Infof("filepath.Walk(%s) skip . and ..\n", path)
			return nil
		}

		if info.IsDir() {
			klog.V(4).Infof("IsDir\n")
			return nil
		}
		newRel := path[lenDst+1:]
		klog.V(6).Infof("newRel: %s\n", newRel)
		dstMap[newRel] = &DiffFile{
			Path:    path,
			RelPath: newRel,
			Attr:    &info,
		}
		return nil
	})
	if err != nil {
		klog.Errorf("filepath.Walk(%s) Err: %v\n", dstPath, err)
		return err
	}

	klog.V(6).Infof("File comparison...\n")

	for key, val := range srcMap {
		dst := dstMap[key]
		if dst == nil {
			klog.V(3).Infof("[ADDING] %s because dst is missing file.", val.Path)
			*diff = append(*diff, &DiffCompare{
				SrcFile:   val,
				DstFile:   nil,
				Direction: DIRECTION_SRC_TO_DST,
			})
			continue
		}
		if (*dst.Attr).ModTime().Before((*val.Attr).ModTime()) {
			srcHash, err := d.getHash(val.Path)
			if err != nil {
				klog.Errorf("Error calculating sha256(%s)\n", val.Path)
				continue
			}
			dstHash, err := d.getHash(dst.Path)
			if err != nil {
				klog.Errorf("Error calculating sha256(%s)\n", dst.Path)
				continue
			}

			if srcHash != dstHash {
				klog.V(3).Infof("[ADDING] %s hash: %s -> %s hash: %s\n", val.Path, srcHash, dst.Path, dstHash)
				*diff = append(*diff, &DiffCompare{
					SrcFile:   val,
					DstFile:   dst,
					Direction: DIRECTION_SRC_TO_DST,
				})
			}
		} else if (*dst.Attr).ModTime().After((*val.Attr).ModTime()) {
			srcHash, err := d.getHash(val.Path)
			if err != nil {
				klog.Errorf("Error calculating sha256(%s)\n", val.Path)
				continue
			}
			dstHash, err := d.getHash(dst.Path)
			if err != nil {
				klog.Errorf("Error calculating sha256(%s)\n", dst.Path)
				continue
			}

			if srcHash != dstHash {
				klog.V(3).Infof("[ADDING] %s hash: %s <- %s hash: %s\n", val.Path, srcHash, dst.Path, dstHash)
				*diff = append(*diff, &DiffCompare{
					SrcFile:   val,
					DstFile:   dst,
					Direction: DIRECTION_DST_TO_SRC,
				})
			}
		}
	}

	for key, val := range dstMap {
		src := srcMap[key]
		if src == nil {
			klog.V(3).Infof("[ADDING] %s because src is missing file.", val.Path)
			*diff = append(*diff, &DiffCompare{
				SrcFile:   nil,
				DstFile:   val,
				Direction: DIRECTION_DST_TO_SRC,
			})
		}
	}

	return nil
}

func (d *Diff) resolveDifferences(diffs *[]*DiffCompare) error {
	for _, diff := range *diffs {
		switch diff.Direction {
		case DIRECTION_SRC_TO_DST:
			newDst := filepath.Join(d.options.RootDstPath, diff.SrcFile.RelPath)
			err := d.buildDir(newDst)
			if err != nil {
				klog.Errorf("buildDir(%s) failed. Err: %v\n", newDst, err)
				return err
			}
			_, err = d.copy(diff.SrcFile.Path, newDst)
			if err != nil {
				klog.Errorf("copy(%s, %s) failed. Err: %v\n", diff.SrcFile.Path, newDst, err)
				return err
			}

			klog.V(4).Infof("[SRC -> DST] Paths: %s to %s\n", diff.SrcFile.Path, newDst)
			if d.options.DryRun {
				klog.Infof("[SRC -> DST] Diff: %s\n", diff.SrcFile.RelPath)
			} else {
				klog.Infof("[SRC -> DST] Copying... %s\n", diff.SrcFile.RelPath)
			}
			if diff.DstFile == nil {
				klog.Infof("\tDestination file does not exist\n")
			} else if diff.SrcFile.Hash != "" && diff.SrcFile.Hash != diff.DstFile.Hash {
				klog.Infof("\tHash mismatch: %s -> %s\n", diff.SrcFile.Hash, diff.DstFile.Hash)
			} else {
				srcTime := (*diff.SrcFile.Attr).ModTime()
				dstTime := (*diff.DstFile.Attr).ModTime()
				klog.Infof("\tSrc Mod Time: %d-%02d-%02dT%02d:%02d:%02d != Dst Mod Time: %d-%02d-%02dT%02d:%02d:%02d\n",
					srcTime.Year(), srcTime.Month(), srcTime.Day(),
					srcTime.Hour(), srcTime.Minute(), srcTime.Second(),
					dstTime.Year(), dstTime.Month(), dstTime.Day(),
					dstTime.Hour(), dstTime.Minute(), dstTime.Second(),
				)
			}
			klog.Infof("\n")
		case DIRECTION_DST_TO_SRC:
			if d.options.SkipSrcUpdate {
				klog.V(3).Infof("Skipping src update because SkipSrcUpdate is true\n")
				continue
			}

			newSrc := filepath.Join(d.options.RootSrcPath, diff.DstFile.RelPath)
			err := d.buildDir(newSrc)
			if err != nil {
				klog.Errorf("buildDir(%s) failed. Err: %v\n", newSrc, err)
				return err
			}
			_, err = d.copy(diff.DstFile.Path, newSrc)
			if err != nil {
				klog.Errorf("copy(%s, %s) failed. Err: %v\n", diff.DstFile.Path, newSrc, err)
				return err
			}

			klog.V(4).Infof("[DST -> SRC] Paths: %s to %s\n", diff.DstFile.Path, newSrc)
			if d.options.DryRun {
				klog.Infof("[DST -> SRC] Diff: %s\n", diff.DstFile.RelPath)
			} else {
				klog.Infof("[DST -> SRC] Copying... %s\n", diff.DstFile.RelPath)
			}
			if diff.SrcFile == nil {
				klog.Infof("\tSource file does not exist\n")
			} else if diff.SrcFile.Hash != "" && diff.SrcFile.Hash != diff.DstFile.Hash {
				klog.Infof("\tHash mismatch: %s -> %s\n", diff.SrcFile.Hash, diff.DstFile.Hash)
			} else {
				srcTime := (*diff.SrcFile.Attr).ModTime()
				dstTime := (*diff.DstFile.Attr).ModTime()
				klog.Infof("\tSrc Mod Time: %d-%02d-%02dT%02d:%02d:%02d != Dst Mod Time: %d-%02d-%02dT%02d:%02d:%02d\n",
					srcTime.Year(), srcTime.Month(), srcTime.Day(),
					srcTime.Hour(), srcTime.Minute(), srcTime.Second(),
					dstTime.Year(), dstTime.Month(), dstTime.Day(),
					dstTime.Hour(), dstTime.Minute(), dstTime.Second(),
				)
			}
			klog.Infof("\n")
		default:
			klog.Errorf("Unknown direction: %d\n", diff.Direction)
			return ErrUnknownDirection
		}
	}

	if d.options.DryRun {
		return nil
	}

	if len(*diffs) > 0 {
		klog.Infof("\n\n")
		klog.Infof("Copied files:\n")
		for _, diff := range *diffs {
			switch diff.Direction {
			case DIRECTION_SRC_TO_DST:
				klog.Infof("[SRC -> DST] Copied %s\n", diff.SrcFile.RelPath)
			case DIRECTION_DST_TO_SRC:
				klog.Infof("[DST -> SRC] Copied %s\n", diff.DstFile.RelPath)
			default:
				klog.Errorf("Unknown direction: %d\n", diff.Direction)
				return ErrUnknownDirection
			}
		}
	}
	return nil
}

func (d *Diff) getHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := make([]byte, 8194)
	hash := sha256.New()
	for {
		srcN, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		dstN, err := hash.Write(buf[:srcN])
		if err != nil {
			return "", err
		}
		if srcN != dstN {
			return "", ErrDiffSizeCopied
		}
	}

	return base64.URLEncoding.EncodeToString(hash.Sum(nil)), nil
}

func (d *Diff) buildDir(path string) error {
	dir := filepath.Dir(path)

	if d.options.DryRun {
		klog.V(3).Infof("DryRun: MkdirAll(%s)\n", dir)
		return nil
	}

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		klog.Errorf("MkdirAll failed. Err: %v\n", err)
		return err
	}

	return nil
}

func (d *Diff) copy(src, dst string) (int64, error) {
	if d.options.DryRun {
		klog.V(3).Infof("DryRun: copy(%s, %s)\n", src, dst)
		return 0, nil
	}

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
		return nBytes, fmt.Errorf("copy byte size mismatch. src: %d != dst: %d", sourceFileStat.Size(), nBytes)
	}
	return nBytes, err
}
