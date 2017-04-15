package strong

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(buf []byte) error {
	u, err := strconv.Unquote(string(buf))
	if err != nil {
		return err
	}
	v, err := time.Parse("2006-01-02T15:04:05 MST", u)
	if err != nil {
		return err
	}
	*t = Time{v.UTC()}
	return nil
}

type Bool bool

func (b *Bool) UnmarshalJSON(buf []byte) error {
	u, err := strconv.Unquote(string(buf))
	if err != nil {
		return err
	}
	switch u {
	case "yes":
		*b = Bool(true)
	default:
		*b = Bool(false)
	}
	return nil
}

type Pin struct {
	No       int     `json:"no"`
	Location string  `json:"location"`
	Channel  string  `json:"channel"`
	Peak     float64 `json:"peak"`
	Azimuth  float64 `json:"azimuth"`
	Dip      float64 `json:"dip"`
	Reversed Bool    `json:"reversed"`
}

type File struct {
	Filename  string  `json:"filename"`
	Site      string  `json:"site"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Elevation float64 `json:"height"`
	Pins      []Pin   `json:"pin"`
	Trigger   Time    `json:"trigger"`
}

type Files map[string][]File

func DecodeFiles(buf []byte, files *Files) error {
	return json.Unmarshal(buf, files)
}

func (f File) Location() string {
	for _, p := range f.Pins {
		return p.Location
	}
	return ""
}

func (f File) Label() string {
	return strings.Join([]string{
		f.Trigger.Format("20060102"),
		f.Trigger.Format("150405"),
		f.Site,
	}, "_")
}

func (f File) Id() string {
	return strings.Join([]string{
		f.Trigger.Format("20060102"),
		f.Trigger.Format("150405"),
		f.Site,
		f.Location(),
	}, "_")
}
