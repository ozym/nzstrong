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

type V3A struct {
	// Event details
	Origin time.Time

	// where to store the V3A files
	Processed string

	// how to convert into V3A file format
	Exec string
}

func (v V3A) Path(file strong.File) string {

	return filepath.Join([]string{
		v.Processed,
		"Proc",
		v.Origin.Format("2006"),
		v.Origin.Format("01_Jan"),
		v.Origin.Format("2006-01-02_150405"),
		"Vol3",
		"data",
		file.Id() + ".V3A",
	}...)
}

func (v V3A) Build(v2a V2A, file strong.File) error {

	input := filepath.Join(filepath.Dir(v.Path(file)), filepath.Base(v2a.Path(file)+"_"))
	log.Println(v2a.Path(file), input)

	if err := copyFiles(v2a.Path(file), input); err != nil {
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
