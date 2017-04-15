package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

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

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Manage strong motion earthquake processing\n")
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

	// process events ...
	for _, feature := range search.Features {
		if event, err := feature.Earthquake(&agency); err == nil {
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
		} else {
			log.Printf("error: %v", err)
		}
	}
}
