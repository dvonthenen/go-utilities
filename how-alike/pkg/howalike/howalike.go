// Copyright 2023 dvonthenen/go-utilities contributors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package howalike

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	matchr "github.com/antzucaro/matchr"
)

func New(opts HowAlikeOptions) *HowAlike {
	dist := &HowAlike{
		options: opts,
	}
	return dist
}

func (h *HowAlike) Process() error {
	actualContents, err := h.readFile(h.options.ActualFile)
	if err != nil {
		return err
	}
	checkContents, err := h.readFile(h.options.CheckFile)
	if err != nil {
		return err
	}

	// err = h.dumpFile(h.options.ActualFile+".INT", actualContents)
	// if err != nil {
	// 	return err
	// }
	// err = h.dumpFile(h.options.CheckFile+".INT", checkContents)
	// if err != nil {
	// 	return err
	// }

	percent := matchr.JaroWinkler(checkContents, actualContents, false)
	fmt.Printf("JaroWinkler: %f\n", percent)

	return err
}

func (h *HowAlike) readFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	r, err := regexp.Compile("^(\\[SPEAKER [0-9]\\:\\] )")
	if err != nil {
		fmt.Printf("MatchString err: %v\n", err)
		return "", err
	}

	var sb strings.Builder

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sb.WriteString(r.ReplaceAllString(scanner.Text(), ""))
	}

	return sb.String(), nil
}

func (h *HowAlike) dumpFile(filename string, contents string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(contents)
	if err != nil {
		return err
	}

	return nil
}
