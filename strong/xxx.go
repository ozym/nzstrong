package strong

/*

import (
	"log"

	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)


type Trigger struct {
	Origin        time.Time
	Magnitude     float64
	MagnitudeType string
	Latitude      float64
	Longitude     float64
	Depth         float64

	Id       string
	Label    string
	Location string
	Path     string

	Processed string
	Config    string

	Vol1 string
	Vol2 string
	Vol3 string
}

func (t Trigger) Volume(number int) string {

	return filepath.Join([]string{
		t.Processed,
		"Proc",
		t.Origin.Format("2006"),
		t.Origin.Format("01_Jan"),
		t.Origin.Format("2006-01-02_150405"),
		"Vol" + strconv.Itoa(number),
		"data",
		t.Id + ".V" + strconv.Itoa(number) + "A",
	}...)
}

func (t Trigger) V1A() (bool, error) {

	var vargs []string

	if t.Config != "" {
		vargs = append(vargs, []string{"-c", t.Config}...)
	}
	vargs = append(vargs, []string{"-z", strconv.FormatFloat(t.Depth, 'g', -1, 64)}...)
	vargs = append(vargs, []string{"-e", t.Origin.Format("2006-01-02T15:04:05")}...)
	vargs = append(vargs, []string{"-n", t.Label}...)
	vargs = append(vargs, []string{"-m", t.MagnitudeType}...)
	vargs = append(vargs, []string{"-x", strconv.FormatFloat(t.Longitude, 'g', -1, 64)}...)
	vargs = append(vargs, []string{"-y", strconv.FormatFloat(t.Latitude, 'g', -1, 64)}...)
	vargs = append(vargs, []string{"-l", strconv.FormatFloat(t.Magnitude, 'g', -1, 64)}...)
	vargs = append(vargs, []string{"-s", t.Location}...)
	vargs = append(vargs, []string{"-o", t.Volume(1)}...)
	vargs = append(vargs, t.Path)

	if err := os.MkdirAll(filepath.Dir(t.Volume(1)), 0755); err != nil {
		return false, err
	}

	log.Println(t.Vol1, strings.Join(vargs, " "))
	if _, err := exec.Command(t.Vol1, vargs...).Output(); err != nil {
		return false, err
	}

	if err := os.Chmod(t.Volume(1), 0644); err != nil {
		return false, err
	}

	return true, nil
}

func (t Trigger) V2A() error {

	input := t.Volume(1) + "_"
	log.Println(input)

	if err := copyFiles(t.Volume(1), input); err != nil {
		return err
	}
	defer os.Remove(input)

	cmd := exec.Command(t.Vol2, filepath.Base(input))
	cmd.Dir = filepath.Dir(input)
	log.Println(cmd.Dir, t.Vol2, filepath.Base(input))
	if _, err := cmd.Output(); err != nil {
		return err
	}
	defer os.Remove(t.Volume(2) + "_")

	if err := copyFiles(t.Volume(2)+"_", t.Volume(2)); err != nil {
		return err
	}

	return nil
}

func (t Trigger) V3A() error {

	input := t.Volume(2) + "_"
	log.Println(input)

	if err := copyFiles(t.Volume(2), input); err != nil {
		return err
	}
	defer os.Remove(input)

	cmd := exec.Command(t.Vol3, filepath.Base(input))
	cmd.Dir = filepath.Dir(input)
	log.Println(cmd.Dir, t.Vol3, filepath.Base(input))
	if _, err := cmd.Output(); err != nil {
		return err
	}
	defer os.Remove(t.Volume(3) + "_")

	if err := copyFiles(t.Volume(3)+"_", t.Volume(3)); err != nil {
		return err
	}

	return nil
}

*/
