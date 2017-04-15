package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func equalFiles(f1, f2 string) bool {
	sf, err := os.Open(f1)
	if err != nil {
		return false
	}
	defer sf.Close()

	df, err := os.Open(f2)
	if err != nil {
		return false
	}
	defer df.Close()

	sscan := bufio.NewScanner(sf)
	dscan := bufio.NewScanner(df)

	for sscan.Scan() {
		dscan.Scan()
		if !bytes.Equal(sscan.Bytes(), dscan.Bytes()) {
			return false
		}
	}

	return true
}

func copyFiles(f1, f2 string) error {
	sf, err := os.Open(f1)
	if err != nil {
		return err
	}
	defer sf.Close()

	if err := os.MkdirAll(filepath.Dir(f2), 0755); err != nil {
		return err
	}

	file, err := ioutil.TempFile(filepath.Dir(f2), "xxx")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	if _, err := io.Copy(file, sf); err != nil {
		return err
	}

	if err := os.Rename(file.Name(), f2); err != nil {
		return err
	}

	if err := os.Chmod(f2, 0644); err != nil {
		return err
	}

	return nil
}
