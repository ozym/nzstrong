package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

type PDF struct {
	Exec string
}

func (p PDF) Convert(ps, pdf string) error {

	if err := os.MkdirAll(filepath.Dir(pdf), 0755); err != nil {
		return err
	}
	cmd := exec.Command(p.Exec, ps, pdf)
	if _, err := cmd.Output(); err != nil {
		return err
	}

	return nil
}

type PS struct {
	Dir  string
	Exec string
}

func (p PS) Plot(file string) error {

	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return err
	}

	cmd := exec.Command(p.Exec, filepath.Base(file))
	cmd.Dir = filepath.Dir(file)
	cmd.Env = []string{"PGPLOT_DIR=" + p.Dir}

	if _, err := cmd.Output(); err != nil {
		return err
	}

	return nil
}