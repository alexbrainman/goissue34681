// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goissue34681_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/alexbrainman/goissue34681"
)

func displayFile(t *testing.T, path string) {
	t.Helper()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Logf("display faild: %v", err)
		return
	}
	t.Logf("displaying %v", path)
	lines := bytes.Split(data, []byte{'\n'})
	for _, line := range lines {
		t.Logf("%q", line)
	}
}

func TestRotateFiles(t *testing.T) {
	temp, err := ioutil.TempDir("", "TestRotateFiles")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(temp)

	logpath := filepath.Join(temp, "log.txt")

	var renameCounter int

	backupName := func(counter int) string {
		return filepath.Join(temp, fmt.Sprintf("log%d.txt", counter))
	}

	rename := func() error {
		renameCounter++
		return os.Rename(logpath, backupName(renameCounter))
	}

	f1, err := goissue34681.OpenFile(logpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer f1.Close()

	write := func(f *os.File, s string) {
		_, err := f.Write([]byte(s))
		if err != nil {
			t.Fatal(err)
		}
	}

	write(f1, "line1\n")

	err = rename()
	if err != nil {
		t.Fatal(err)
	}

	write(f1, "line2\n")

	f2, err := goissue34681.OpenFile(logpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer f1.Close()

	write(f2, "line3\n")

	write(f1, "line4\n")

	err = rename()
	if err != nil {
		t.Fatal(err)
	}

	f3, err := goissue34681.OpenFile(logpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer f3.Close()

	write(f3, "liner5\n")

	f4, err := goissue34681.OpenFile(logpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer f4.Close()

	write(f4, "liner6\n")
	write(f3, "liner7\n")

	f5, err := goissue34681.Open(logpath)
	if err != nil {
		t.Fatal(err)
	}
	defer f5.Close()

	// display files

	displayFile(t, logpath)
	for i := 1; i <= renameCounter; i++ {
		displayFile(t, backupName(i))
	}
}
