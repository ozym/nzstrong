package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"time"

	"github.com/ozym/nzstrong/quake"
	"github.com/ozym/nzstrong/strong"
	/*
		"flag"
		"fmt"
		"log"
		"os"
		"regexp"
		//	"strconv"
		"strings"
		//	"time"
	*/)

type Gather struct {
	altus     Altus
	outside   time.Duration
	inside    time.Duration
	processed string
	within    float64
	ps2pdf    string
	plot1     string
	plot2     string
	plot3     string
	plot4     string
	pgplot    string
	vol2      string
	vol3      string
	std       string
}

func (g Gather) Process(path string) error {
	var eq quake.Earthquake

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := xml.Unmarshal(raw, &eq); err != nil {
		return nil
	}

	return g.Build(&eq)
}

func (g Gather) Build(eq *quake.Earthquake) error {

	origin, err := time.Parse("2006-01-02T15:04:05", eq.Origin[0:19])
	if err != nil {
		return err
	}

	list := make(map[string]interface{})
	for t := origin.Add(-g.inside); origin.Add(g.outside).After(t); t = t.Add(time.Minute) {
		files, err := g.altus.Files(t)
		if err != nil {
			return err
		}
		for _, f := range files {
			list[f] = true
		}
	}

	var files []string
	for f, _ := range list {
		files = append(files, f)
	}

	check, err := g.altus.Gather(eq.Latitude, eq.Longitude, g.within, files)
	if err != nil {
		return err
	}
	log.Println(check)

	altus := strong.Volume{
		Processed: g.processed,
		Source:    g.std,
		Ps2pdf:    g.ps2pdf,
		PgPlot:    g.pgplot,
		Vol2:      g.vol2,
		Vol3:      g.vol3,
		Plot1:     g.plot1,
		Plot2:     g.plot2,
		Plot3:     g.plot3,
		Plot4:     g.plot4,
	}
	for _, f := range check {
		if err := altus.Build(strong.Trigger{
			Origin:        origin,
			MagnitudeType: eq.MagnitudeType,
			Magnitude:     eq.Magnitude,
			Latitude:      eq.Latitude,
			Longitude:     eq.Longitude,
			Depth:         eq.Depth,
			Path:          f.Filename,
			Id:            f.Id(),
			Label:         f.Label(),
			Location:      f.Location(),
			Builder:       g.altus.Builder,
			Config:        g.altus.Config,
		}); err != nil {
			return err
		}
	}

	csv := strong.CSV{
		Processed: g.processed,

		Origin:        origin,
		PublicID:      eq.PublicID,
		MagnitudeType: eq.MagnitudeType,
		Magnitude:     eq.Magnitude,
		Latitude:      eq.Latitude,
		Longitude:     eq.Longitude,
		Depth:         eq.Depth,

		Before: g.outside,
		After:  g.inside,
	}

	matched, err := csv.Files(check)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Printf("matched %d files", len(matched))

	elements, err := csv.Elements(matched)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Printf("matched %d observations", len(elements))

	if err := csv.Store(csv.Path(), elements); err != nil {
		log.Fatalf("error: unable to store csv pga file: %v", err)
	}
	if err := csv.Store(csv.PGA(), elements); err != nil {
		log.Fatalf("error: unable to store csv summary file: %v", err)
	}

	return nil

	volumes := []Volume{
		Volume{
			Origin:    origin,
			Processed: g.processed,
			PS: PS{
				Exec: g.plot1,
				Dir:  g.pgplot,
			},
			PDF: PDF{
				Exec: g.ps2pdf,
			},
			Number: "Vol1",
			Suffix: ".V1A",
		},
		Volume{
			Origin:    origin,
			Processed: g.processed,
			PS: PS{
				Exec: g.plot2,
				Dir:  g.pgplot,
			},
			PDF: PDF{
				Exec: g.ps2pdf,
			},
			Number: "Vol2",
			Suffix: ".V2A",
			Exec:   g.vol2,
		},
		Volume{
			Origin:    origin,
			Processed: g.processed,
			PS: PS{
				Exec: g.plot3,
				Dir:  g.pgplot,
			},
			PDF: PDF{
				Exec: g.ps2pdf,
			},
			Number: "Vol3",
			Suffix: ".V3A",
			Exec:   g.vol3,
			Source: g.std,
		},
	}

	var v1a V1A

	for _, m := range matched {
		ok, err := v1a.Build(m)
		if err != nil {
			log.Fatalf("unable to build v1a file %s: %v", m.Filename, err)
		}
		if !ok {
			continue
		}

		for i := 1; i < len(volumes); i++ {
			if err := volumes[i].Build(volumes[i-1], m); err != nil {
				log.Fatalf("unable to build volume file %s (%s): %v", m.Filename, volumes[i].Number, err)
			}
		}
		for _, v := range volumes {
			if err := v.Plot(v, m); err != nil {
				log.Fatalf("unable to plot volume file %s (%s): %v", m.Filename, v.Number, err)
			}
		}
		/*
			if err := v4a.Plot(volumes[1], m); err != nil {
				log.Fatalf("unable to plot volume file %s (%s): %v", m.Filename, "Vol4", err)
			}
		*/
	}

	return nil
}

/*
package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	//	"time"
	//"strconv"
)

// Handle raw files ...
type Find struct {
	// Where to find files
	Spool string
	// How to identify them
	Match string

	// decoding program
	Exec string
	// configuration file
	Config string
}
*/

/*
func (g Gather) Files(event *Earthquake) ([]File, error) {

	// handle empty events
	if event == nil {
		return nil, nil
	}

	ctmpl, err := template.New("csd").Parse(f.Spool)
	if err != nil {
		return nil, err
	}
	var dir bytes.Buffer
	if err := ctmpl.Execute(&dir, event); err != nil {
		return nil, err
	}

	var files []File
	if err := filepath.Walk(dir.String(), func(path string, info os.FileInfo, err error) error {
		m, err := filepath.Match(f.Match, filepath.Base(path))
		if err != nil || m == false {
			return err
		}
		bin := strings.Fields(f.Exec)[0]
		var args []string
		if f.Config != "" {
			args = append(args, []string{"-c", f.Config}...)
		}
		if len(strings.Fields(f.Exec)) > 1 {
			args = append(args, strings.Fields(f.Exec)[1:]...)
		}
		args = append(args, path)

		out, err := exec.Command(bin, args...).Output()
		if err != nil {
			return err
		}

		var list Files
		if err := DecodeFiles(out, &list); err != nil {
			return err
		}
		for _, v := range list {
			files = append(files, v...)
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return files, nil
}
*/
