package strong

import (
	"log"

	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

type Trigger struct {
	Builder       string
	Config        string
	Path          string
	Id            string
	Label         string
	Origin        time.Time
	Latitude      float64
	Longitude     float64
	Depth         float64
	Magnitude     float64
	MagnitudeType string
	Location      string
}

func (t Trigger) V1A(processed string) string {

	return filepath.Join([]string{
		processed,
		"Proc",
		t.Origin.Format("2006"),
		t.Origin.Format("01_Jan"),
		t.Origin.Format("2006-01-02_150405"),
		"Vol1",
		"data",
		t.Id + ".V1A",
	}...)
}

func (t Trigger) V2A(processed string) string {

	return filepath.Join([]string{
		processed,
		"Proc",
		t.Origin.Format("2006"),
		t.Origin.Format("01_Jan"),
		t.Origin.Format("2006-01-02_150405"),
		"Vol2",
		"data",
		t.Id + ".V2A",
	}...)
}

func (t Trigger) V3A(processed string) string {

	return filepath.Join([]string{
		processed,
		"Proc",
		t.Origin.Format("2006"),
		t.Origin.Format("01_Jan"),
		t.Origin.Format("2006-01-02_150405"),
		"Vol3",
		"data",
		t.Id + ".V3A",
	}...)
}

func (t Trigger) PSPath(processed string, volume string) string {
	return filepath.Join([]string{
		processed,
		"Proc",
		t.Origin.Format("2006"),
		t.Origin.Format("01_Jan"),
		t.Origin.Format("2006-01-02_150405"),
		volume,
		"data",
		"V" + volume[len(volume)-1:] + "A_" + t.Id + ".ps",
	}...)
}

func (t Trigger) PDFPath(processed string, volume string) string {
	return filepath.Join([]string{
		processed,
		"Proc",
		t.Origin.Format("2006"),
		t.Origin.Format("01_Jan"),
		t.Origin.Format("2006-01-02_150405"),
		volume,
		"plots",
		t.Id + ".PDF",
	}...)
}

func (t Trigger) Build(processed string) error {

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
	vargs = append(vargs, []string{"-o", t.V1A(processed)}...)
	vargs = append(vargs, t.Path)

	if err := os.MkdirAll(filepath.Dir(t.V1A(processed)), 0755); err != nil {
		return err
	}

	log.Println(t.Builder, vargs)
	if _, err := exec.Command(t.Builder, vargs...).Output(); err != nil {
		return err
	}

	if err := os.Chmod(t.V1A(processed), 0644); err != nil {
		return err
	}

	return nil
}
