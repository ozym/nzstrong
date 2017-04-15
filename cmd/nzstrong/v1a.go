package main

import (
	"log"
	"strings"

	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ozym/nzstrong/strong"
)

type V1A struct {
	// Event details
	Origin        time.Time
	Latitude      float64
	Longitude     float64
	Depth         float64
	MagnitudeType string
	Magnitude     float64

	// where to store the V3A files
	Processed string
	// how to identify input format
	Match string

	// how to convert into V3A file format
	Exec string

	// converter configuration file
	Config string
}

func (v V1A) Path(file strong.File) string {

	return filepath.Join([]string{
		v.Processed,
		"Proc",
		v.Origin.Format("2006"),
		v.Origin.Format("01_Jan"),
		v.Origin.Format("2006-01-02_150405"),
		"Vol1",
		"data",
		file.Id() + ".V1A",
	}...)
}

func (v V1A) Build(file strong.File) (bool, error) {

	match, err := filepath.Match(v.Match, filepath.Base(file.Filename))
	if err != nil || !match {
		return false, err
	}

	var vargs []string

	if v.Config != "" {
		vargs = append(vargs, []string{"-c", v.Config}...)
	}
	vargs = append(vargs, []string{"-z", strconv.FormatFloat(v.Depth, 'g', -1, 64)}...)
	vargs = append(vargs, []string{"-e", v.Origin.Format("2006-01-02T15:04:05")}...)
	vargs = append(vargs, []string{"-n", file.Label()}...)
	vargs = append(vargs, []string{"-m", v.MagnitudeType}...)
	vargs = append(vargs, []string{"-x", strconv.FormatFloat(v.Longitude, 'g', -1, 64)}...)
	vargs = append(vargs, []string{"-y", strconv.FormatFloat(v.Latitude, 'g', -1, 64)}...)
	vargs = append(vargs, []string{"-l", strconv.FormatFloat(v.Magnitude, 'g', -1, 64)}...)
	vargs = append(vargs, []string{"-s", file.Location()}...)
	vargs = append(vargs, []string{"-o", v.Path(file)}...)
	vargs = append(vargs, file.Filename)

	if err := os.MkdirAll(filepath.Dir(v.Path(file)), 0755); err != nil {
		return false, err
	}

	log.Println(v.Exec, strings.Join(vargs, " "))
	if _, err := exec.Command(v.Exec, vargs...).Output(); err != nil {
		return false, err
	}

	if err := os.Chmod(v.Path(file), 0644); err != nil {
		return false, err
	}

	return true, nil
}
