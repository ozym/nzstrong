package main

import (
	"log"

	"os"
	"os/exec"
	"path/filepath"
	//"strconv"
	"time"

	"github.com/ozym/nzstrong/strong"
)

type V2A struct {
	// Event details
	Origin time.Time

	// where to store the V2A files
	Processed string

	// how to convert into V2A file format
	Exec string
}

func (v V2A) Path(file strong.File) string {

	return filepath.Join([]string{
		v.Processed,
		"Proc",
		v.Origin.Format("2006"),
		v.Origin.Format("01_Jan"),
		v.Origin.Format("2006-01-02_150405"),
		"Vol2",
		"data",
		file.Id() + ".V2A",
	}...)
}

func (v V2A) Build(v1a V1A, file strong.File) error {

	input := filepath.Join(filepath.Dir(v.Path(file)), filepath.Base(v1a.Path(file)+"_"))
	log.Println(v1a.Path(file), input)

	if err := copyFiles(v1a.Path(file), input); err != nil {
		return err
	}
	defer os.Remove(input)

	cmd := exec.Command(v.Exec, filepath.Base(input))
	cmd.Dir = filepath.Dir(input)
	log.Println(cmd.Dir, v.Exec, filepath.Base(input))
	if _, err := cmd.Output(); err != nil {
		return err
	}
	defer os.Remove(v.Path(file) + "_")

	if err := copyFiles(v.Path(file)+"_", v.Path(file)); err != nil {
		return err
	}

	return nil
}
