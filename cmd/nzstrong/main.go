package main

import (
	"flag"
	"fmt"
	//"io/ioutil"
	"log"
	"os"
	//"regexp"
	"strconv"
	"strings"
	"time"
	//	"time"

	"github.com/ozym/nzstrong/quake"
)

func main() {

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "make noise")

	var service string
	flag.StringVar(&service, "service", "wfs.geonet.org.nz", "earthquake query service")
	var agency string
	flag.StringVar(&agency, "agency", "WEL", "earthquake agency service")

	var minmag float64
	flag.Float64Var(&minmag, "minmag", 3.0, "minimum magnitude to process, use 0.0 for no limit")
	var maxmag float64
	flag.Float64Var(&maxmag, "maxmag", 0.0, "maximum magnitude to process, use 0.0 for no limit")

	var since time.Duration
	flag.DurationVar(&since, "since", 30*time.Minute, "modified event search window, use 0 for no offset")

	var after time.Duration
	flag.DurationVar(&after, "after", 0, "modified event search window offset, use 0 for no offset")

	var eventType string
	flag.StringVar(&eventType, "type", "earthquake", "event type query parameter")

	var evaluationStatus string
	flag.StringVar(&evaluationStatus, "status", "confirmed", "event status query parameter")

	var evaluationMode string
	flag.StringVar(&evaluationMode, "mode", "manual", "event mode query parameter")

	var limit int
	flag.IntVar(&limit, "limit", 0, "maximum number of records to process before filters, use 0 for no limit")

	var spool string
	flag.StringVar(&spool, "spool", ".", "output spool directory")

	var match string
	flag.StringVar(&match, "match", ".xml$", "match event files")

	var remove bool
	flag.BoolVar(&remove, "remove", false, "remove event files")

	var age time.Duration
	flag.DurationVar(&age, "age", 0, "age after which events will be removed if set to")

	var outside time.Duration
	flag.DurationVar(&outside, "outside", 30*time.Second, "select triggers this time before an event")

	var inside time.Duration
	flag.DurationVar(&inside, "inside", 120*time.Second, "select triggers this time after an event")

	var within float64
	flag.Float64Var(&within, "within", 300, "quake requires to be within a given distance")

	var altus_config string
	flag.StringVar(&altus_config, "altus-config", "altus.xml", "altus config file")

	var altus_id string
	flag.StringVar(&altus_id, "altus-id", "esaltusid", "esaltusid program")

	var altus_builder string
	flag.StringVar(&altus_builder, "altus-builder", "/usr/bin/esaltus", "evt file converter")

	var altus_archive string
	flag.StringVar(&altus_archive, "altus-archive", "./test/2006/2006.%j", "where to find raw altus files")

	var altus_match string
	flag.StringVar(&altus_match, "altus-match", "2006.%j.1504.*.evt", "how to find altus files")

	var processed string
	flag.StringVar(&processed, "processed", "./processed", "provide the output processed directory")

	var std string
	flag.StringVar(&std, "std", ".", "esv3std.dat file directory")

	var pgplot string
	flag.StringVar(&pgplot, "pgplot", "/usr/local/share/pgplot", "esv3std.dat file directory")

	var ps2pdf string
	flag.StringVar(&ps2pdf, "ps2pdf", "ps2pdf", "convert postscript files to pdf")

	var plot1 string
	flag.StringVar(&plot1, "esplot_v1", "/usr/bin/esplot_v1", "vol1 plotter")

	var plot2 string
	flag.StringVar(&plot2, "esplot_v2", "/usr/bin/esplot_v2", "vol2 plotter")

	var plot3 string
	flag.StringVar(&plot3, "esplot_v3", "/usr/bin/esplot_v3", "vol3 plotter")

	var plot4 string
	flag.StringVar(&plot4, "esplot_v4", "/usr/bin/esplot_v4", "vol4 plotter")

	var vol2 string
	flag.StringVar(&vol2, "esvol2m", "/usr/bin/esvol2m", "vol2 file converter")

	var vol3 string
	flag.StringVar(&vol3, "esvol3m", "/usr/bin/esvol3m", "vol3 file converter")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Manage strong motion earthquake file processing\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	flag.Parse()

	q := quake.Query{
		Service: service,
		Limit:   limit,
	}

	// simple event and evaluation checks ...
	q.AddFilter("eventtype", "=", eventType)
	q.AddFilter("evaluationstatus", "=", evaluationStatus)
	q.AddFilter("evaluationmode", "=", evaluationMode)

	if flag.NArg() > 0 {
		var filters []string
		for _, id := range flag.Args() {
			filters = append(filters, "publicID+=+"+id)
		}
		q.AddFilter("(", strings.Join(filters, "+or+"), ")")

	} else {
		// check magnitudes are within scope
		if minmag > 0.0 {
			q.AddFilter("magnitude", ">=", strconv.FormatFloat(minmag, 'f', -1, 64))
		}
		if maxmag > 0.0 {
			q.AddFilter("magnitude", "<=", strconv.FormatFloat(maxmag, 'f', -1, 64))
		}

		// perhaps check whether it has been updated recently
		if since > 0 {
			q.AddFilter("modificationtime", ">=", quake.TimeOffsetNow(since))
		}
		if after > 0 {
			q.AddFilter("modificationtime", "<=", quake.TimeOffsetNow(after))
		}
	}

	if verbose {
		log.Println(q.URL().String())
	}

	// query the quake api
	search, err := q.Search()
	if err != nil {
		log.Fatal(err)
	}

	gather := Gather{
		altus: Altus{
			Id:      altus_id,
			Builder: altus_builder,
			Config:  altus_config,
			Match:   altus_match,
			Archive: altus_archive,
		},
		processed: processed,
		outside:   outside,
		inside:    inside,
		within:    within,
		pgplot:    pgplot,
		std:       std,
		ps2pdf:    ps2pdf,
		vol2:      vol2,
		vol3:      vol3,
		plot1:     plot1,
		plot2:     plot2,
		plot3:     plot3,
		plot4:     plot4,
	}

	/*
		check, err := regexp.Compile(match)
		if err != nil {
			log.Fatal(err)
		}

		if verbose {
			log.Println("check: ", spool)
		}

			files, err := ioutil.ReadDir(spool)
			if err != nil {
				log.Fatal(err)
			}

			for _, file := range files {
				if strings.HasPrefix(file.Name(), ".") {
					continue
				}
				if !check.MatchString(file.Name()) {
					continue
				}
				fmt.Println(file.Name())

				if err := gather.Process(file.Name()); err != nil {
					log.Printf("unable to process %s: %v", file.Name(), err)
				}

				if remove && time.Now().Sub(file.ModTime()) > age {
					log.Println("remove ", file.Name())
				}
			}
	*/

	// process events ...
	for _, feature := range search.Features {
		event, err := feature.Earthquake(&agency)
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}
		if err := gather.Build(event); err != nil {
			log.Printf("error: %v", err)
			continue
		}
		log.Println(event)
		/*
			// output xml formatted event files
			if err := os.MkdirAll(spool, 0755); err != nil {
				log.Fatal(err)
			}
			output := fmt.Sprintf("%s/%s-%s.xml", spool, event.PublicID, event.UpdateTime)
			if verbose {
				log.Printf("writing: %s", output)
			}
			if err := event.WriteFile(output); err != nil {
				log.Printf("error: %v", err)
			}
		*/
	}
}
