package strong

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ozym/nzstrong/quake"
)

const Gravity = 9806.65

var CSVHeader = []string{
	"Earthquake Date (UT)",
	"Time (UT)",
	"ID",
	"Mag Type",
	"Magnitude",
	"Depth (km)",
	"R?",
	"Epic. Dist.(km)",
	"PGA Vertical (mm/s/s)",
	"PGA Horiz_1 (mm/s/s)",
	"PGA Horiz_2 (mm/s/s)",
	"PGA Vertical (%g)",
	"PGA Horiz_1 (%g)",
	"PGA Horiz_2 (%g)",
	"Accelerogram ID",
	"Site Code",
	"Name",
	"Site Latitude",
	"Site Longitude",
	"Site Elevation (m)",
}

type Element struct {
	Origin        time.Time
	PublicID      string
	MagnitudeType string
	Magnitude     float64
	Depth         float64
	Distance      float64
	Pga           [3]float64
	ID            string
	Site          string
	Name          string
	Latitude      float64
	Longitude     float64
	Elevation     float64
}

type Elements []Element

func (e Elements) Len() int {
	return len(e)
}
func (e Elements) Less(i, j int) bool {
	return e[i].Distance < e[j].Distance
}
func (e Elements) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

type CSV struct {
	// event parameters
	Origin        time.Time
	PublicID      string
	MagnitudeType string
	Magnitude     float64
	Latitude      float64
	Longitude     float64
	Depth         float64
	// where to store files
	Processed string
	// look before
	Before time.Duration
	// look after
	After time.Duration
}

func (c CSV) Path() string {

	return filepath.Join([]string{
		c.Processed,
		"Proc",
		c.Origin.Format("2006"),
		c.Origin.Format("01_Jan"),
		c.Origin.Format("2006-01-02_150405"),
		c.Origin.Format("20060102_150405") + ".CSV",
	}...)
}

func (c CSV) PGA() string {
	return filepath.Join([]string{
		c.Processed,
		"Summary",
		c.Origin.Format("2006"),
		c.PublicID + "_pga.csv",
	}...)
}

func (c CSV) Store(path string, elements []Element) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	file, err := ioutil.TempFile(filepath.Dir(path), "xxx")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	if err := c.Write(file, elements); err != nil {
		return err
	}

	if err := os.Rename(file.Name(), path); err != nil {
		return err
	}

	if err := os.Chmod(path, 0644); err != nil {
		return err
	}

	return nil
}

func (c *CSV) Files(files []File) ([]File, error) {

	var matched []File

	for _, f := range files {
		if f.Trigger.Before(c.Origin.Add(-c.Before)) {
			continue
		}
		if f.Trigger.After(c.Origin.Add(c.After)) {
			continue
		}

		matched = append(matched, f)
	}

	return matched, nil
}

func (c *CSV) Elements(files []File) ([]Element, error) {

	var elements []Element

	for _, f := range files {
		if f.Trigger.Before(c.Origin.Add(-c.Before)) {
			continue
		}
		if f.Trigger.After(c.Origin.Add(c.After)) {
			continue
		}

		var pga_v, pga_1, pga_2 float64
		for _, p := range f.Pins {
			switch p.Channel[len(p.Channel)-1] {
			case 'Z', 'z':
				pga_v = math.Abs(p.Peak)
			case '1', 'N', 'n', 'y', 'Y':
				pga_1 = math.Abs(p.Peak)
			case '2', 'E', 'e', 'x', 'X':
				pga_2 = math.Abs(p.Peak)
			}
		}

		elements = append(elements, Element{
			Origin:        c.Origin,
			PublicID:      c.PublicID,
			MagnitudeType: c.MagnitudeType,
			Magnitude:     c.Magnitude,
			Depth:         c.Depth,
			Distance:      quake.Point{c.Latitude, c.Longitude}.Distance(quake.Point{f.Latitude, f.Longitude}),
			Pga:           [3]float64{pga_v, pga_1, pga_2},
			ID:            f.Id(),
			Site:          f.Site,
			Name:          f.Name,
			Latitude:      f.Latitude,
			Longitude:     f.Longitude,
			Elevation:     f.Elevation,
		})
	}

	sort.Sort(Elements(elements))

	return elements, nil
}

func (c *CSV) Write(wr io.Writer, elements []Element) error {

	var header []string
	for _, h := range CSVHeader {
		header = append(header, strconv.Quote(h))
	}

	if _, err := fmt.Fprintln(wr, strings.Join(header, ",")); err != nil {
		return err
	}

	sort.Sort(Elements(elements))

	for _, e := range elements {
		var line = []string{
			e.Origin.Format("2006-01-02"),
			e.Origin.Format("15:04:05"),
			e.PublicID,
			e.MagnitudeType,
			strconv.FormatFloat(e.Magnitude, 'f', 2, 64),
			strconv.FormatFloat(e.Depth, 'f', 1, 64),
			"",
			strconv.FormatFloat(e.Distance, 'f', 0, 64),
			strconv.FormatFloat(e.Pga[0], 'f', 1, 64),
			strconv.FormatFloat(e.Pga[1], 'f', 1, 64),
			strconv.FormatFloat(e.Pga[2], 'f', 1, 64),
			strconv.FormatFloat(100.0*e.Pga[0]/Gravity, 'f', 4, 64),
			strconv.FormatFloat(100.0*e.Pga[1]/Gravity, 'f', 4, 64),
			strconv.FormatFloat(100.0*e.Pga[2]/Gravity, 'f', 4, 64),
			e.ID,
			e.Site,
			strconv.Quote(e.Name),
			strconv.FormatFloat(e.Latitude, 'g', -1, 64),
			strconv.FormatFloat(e.Longitude, 'g', -1, 64),
			func() string {
				if e.Elevation == 0.0 {
					return "0.0"
				}
				return strconv.FormatFloat(e.Elevation, 'g', -1, 64)
			}(),
		}
		if _, err := fmt.Fprintln(wr, strings.Join(line, ",")); err != nil {
			return err
		}
	}

	return nil
}
