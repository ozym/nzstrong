package main

import (
	"encoding/xml"
	"io"
	"os"
)

const RFC3339Micro = "2006-01-02T15:04:05.999999Z"

type Earthquake struct {
	XMLName    string `xml:"event"`
	AgencyID   string `xml:"creationInfo>agencyID"`
	UpdateTime string `xml:"creationInfo>updateTime"`

	PublicID string `xml:"publicID,attr"`
	UID      string `xml:"uid"`

	Type   string `xml:"type"`
	Status string `xml:"status"`

	MethodID         string `xml:"methodID"`
	EarthModelID     string `xml:"earthModelID"`
	EvaluationMode   string `xml:"evaluationMode"`
	EvaluationStatus string `xml:"evaluationStatus"`

	Origin    string  `xml:"preferredOrigin>time>value"`
	Latitude  float64 `xml:"preferredOrigin>latitude>value"`
	Longitude float64 `xml:"preferredOrigin>longitude>value"`
	Depth     float64 `xml:"preferredOrigin>depth>value"`

	UsedPhaseCount   int32   `xml:"preferredOrigin>quality>usedPhaseCount"`
	UsedStationCount int32   `xml:"preferredOrigin>quality>usedStationCount"`
	StandardError    float64 `xml:"preferredOrigin>quality>standardError"`
	AzimuthalGap     float64 `xml:"preferredOrigin>quality>azimuthalGap"`
	MinimumDistance  float64 `xml:"preferredOrigin>quality>minimumDistance"`

	Magnitude             float64 `xml:"preferredOrigin>preferredMagnitude>magnitude>value"`
	MagnitudeType         string  `xml:"preferredOrigin>preferredMagnitude>type"`
	MagnitudeStationCount int32   `xml:"preferredOrigin>preferredMagnitude>stationCount"`
}

func (e *Earthquake) Marshal() ([]byte, error) {

	res := ([]byte)(xml.Header)

	b, err := xml.MarshalIndent(e, "", "   ")
	if err != nil {
		return nil, err
	}

	return append(res, b...), nil
}

func (e *Earthquake) Write(wr io.Writer) error {

	b, err := e.Marshal()
	if err != nil {
		return err
	}

	_, err = wr.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func (e *Earthquake) WriteFile(outfile string) error {

	file, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return e.Write(file)
}
