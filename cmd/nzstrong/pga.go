package main

import (
	"path/filepath"
	"time"
)

type Pga struct {
	Origin    time.Time
	PublicID  string
	Processed string
}

func (p Pga) Path() string {
	return filepath.Join([]string{
		p.Processed,
		"Summary",
		p.Origin.Format("2006"),
		p.PublicID + "_pga.csv",
	}...)
}
