package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ozym/nzstrong/quake"
	"github.com/ozym/nzstrong/strong"
)

type Altus struct {
	Id      string
	Config  string
	Builder string
	Match   string
	Archive string
}

func (a Altus) Files(epoch time.Time) ([]string, error) {

	dir := strings.Replace(epoch.Format(a.Archive), "%j", fmt.Sprintf("%03d", epoch.YearDay()), -1)
	match := strings.Replace(epoch.Format(a.Match), "%j", fmt.Sprintf("%03d", epoch.YearDay()), -1)

	fmt.Println(a.Archive, a.Match)
	fmt.Println(dir, match)

	var files []string
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		m, err := filepath.Match(match, filepath.Base(path))
		if err != nil || m == false {
			return err
		}
		files = append(files, path)

		return nil
	}); err != nil {
		return nil, err
	}
	return files, nil
}

func (a Altus) Gather(lat, lon, within float64, paths []string) ([]strong.File, error) {

	point := quake.Point{lat, lon}

	var files []strong.File
	for _, path := range paths {
		bin := strings.Fields(a.Id)[0]
		var args []string
		if a.Config != "" {
			args = append(args, []string{"-c", a.Config}...)
		}
		if len(strings.Fields(a.Id)) > 1 {
			args = append(args, strings.Fields(a.Id)[1:]...)
		}
		args = append(args, path)

		fmt.Println(bin, args)
		out, err := exec.Command(bin, args...).Output()
		if err != nil {
			return nil, err
		}

		var list strong.Files
		if err := strong.DecodeFiles(out, &list); err != nil {
			return nil, err
		}
		for _, v := range list {
			for _, f := range v {
				if point.Distance(quake.Point{f.Latitude, f.Longitude}) > within {
					continue
				}
				files = append(files, f)
			}
		}
	}

	return files, nil
}
