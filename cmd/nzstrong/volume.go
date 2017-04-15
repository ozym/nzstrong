package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ozym/nzstrong/strong"
)

type Volume struct {
	// Event details
	Origin time.Time
	// where to store the volume files
	Processed string
	// Build plots
	// ps plotter
	PS PS
	// pdf converter
	PDF PDF
	// Number
	Number string
	// suffix
	Suffix string
	// how to convert volume files
	Exec string
	// optional src directory enviornment
	Source string
}

func (v Volume) Path(file strong.File) string {

	return filepath.Join([]string{
		v.Processed,
		"Proc",
		v.Origin.Format("2006"),
		v.Origin.Format("01_Jan"),
		v.Origin.Format("2006-01-02_150405"),
		v.Number,
		"data",
		file.Id() + v.Suffix,
	}...)
}

func (v Volume) PSPath(volume Volume, file strong.File) string {
	return filepath.Join([]string{
		v.Processed,
		"Proc",
		v.Origin.Format("2006"),
		v.Origin.Format("01_Jan"),
		v.Origin.Format("2006-01-02_150405"),
		v.Number,
		"data",
		strings.TrimLeft(volume.Suffix, ".") + "_" + file.Id() + ".ps",
	}...)
}

func (v Volume) PDFPath(file strong.File) string {
	return filepath.Join([]string{
		v.Processed,
		"Proc",
		v.Origin.Format("2006"),
		v.Origin.Format("01_Jan"),
		v.Origin.Format("2006-01-02_150405"),
		v.Number,
		"plots",
		file.Id() + ".PDF",
	}...)
}

func (v Volume) Build(volume Volume, file strong.File) error {

	// directory is made up of new output with old file
	input := filepath.Join(filepath.Dir(v.Path(file)), filepath.Base(volume.Path(file)+"_"))

	// copy to _ version
	if err := copyFiles(volume.Path(file), input); err != nil {
		return err
	}
	defer os.Remove(input)

	// will have trouble if output _ file exists
	if _, err := os.Stat(v.Path(file) + "_"); err == nil {
		if err := os.Remove(v.Path(file) + "_"); err != nil {
			return err
		}
	}

	// run the builder
	cmd := exec.Command(v.Exec, filepath.Base(input))
	cmd.Dir = filepath.Dir(input)
	if v.Source != "" {
		cmd.Env = []string{"SRC=" + v.Source}
	}
	if _, err := cmd.Output(); err != nil {
		return err
	}
	defer os.Remove(v.Path(file) + "_")

	if err := copyFiles(v.Path(file)+"_", v.Path(file)); err != nil {
		return err
	}

	return nil
}

func (v Volume) Plot(volume Volume, file strong.File) error {

	// output files
	ps, pdf := volume.PSPath(v, file), v.PDFPath(file)
	defer os.Remove(ps)

	// will have trouble if the ps file already exists
	if _, err := os.Stat(ps); err == nil {
		if err := os.Remove(ps); err != nil {
			return err
		}
	}

	// make a copy of the v?a file ...
	if err := copyFiles(volume.Path(file), volume.Path(file)+"_"); err != nil {
		return err
	}
	defer os.Remove(volume.Path(file) + "_")

	// build ps file
	if err := v.PS.Plot(volume.Path(file) + "_"); err != nil {
		return err
	}

	// convert postscript to pdf
	if err := v.PDF.Convert(ps, pdf); err != nil {
		return err
	}

	return nil
}
