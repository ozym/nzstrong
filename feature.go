package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Feature struct {
	Geometry struct {
		Coordinates [2]float64 `json:"coordinates"`
	} `json:"geometry"`

	Properties struct {
		EventType             *string    `json:"eventtype"`
		PublicID              *string    `json:"publicid"`
		ModificationTime      *time.Time `json:"modificationtime"`
		OriginTime            *time.Time `json:"origintime"`
		OriginError           *float64   `json:"originerror"`
		EarthModel            *string    `json:"earthmodel"`
		EvaluationMethod      *string    `json:"evaluationmethod"`
		EvaluationStatus      *string    `json:"evaluationstatus"`
		EvaluationMode        *string    `json:"evaluationmode"`
		Latitude              *float64   `json:"latitude"`
		Longitude             *float64   `json:"longitude"`
		Depth                 *float64   `json:"depth"`
		DepthType             *string    `json:"depthtype"`
		UsedPhaseCount        *int32     `json:"usedphasecount"`
		UsedStationCount      *int32     `json:"usedstationcount"`
		AzimuthalGap          *float64   `json:"azimuthalgap"`
		MinimumDistance       *float64   `json:"minimumdistance"`
		Magnitude             *float64   `json:"magnitude"`
		MagnitudeType         *string    `json:"magnitudetype"`
		MagnitudeStationCount *int32     `json:"magnitudestationcount"`
		MagnitudeUncertainty  *float64   `json:"magnitudeuncertainty"`
	} `json:"properties"`
}

type Search struct {
	Features []Feature `json:"features"`
}

func (f *Feature) Earthquake(agency *string) (*Earthquake, error) {

	switch {
	case f.Properties.Magnitude == nil:
		return nil, fmt.Errorf("no magnitude found")
	case f.Properties.MagnitudeType == nil:
		return nil, fmt.Errorf("no magnitude type found")
	case f.Properties.ModificationTime == nil:
		return nil, fmt.Errorf("no modification time found")
	case f.Properties.Depth == nil:
		return nil, fmt.Errorf("no depth found")
	case f.Properties.PublicID == nil:
		return nil, fmt.Errorf("no public id found")
	case f.Properties.EventType == nil:
		return nil, fmt.Errorf("no event type found")
	case f.Properties.EvaluationStatus == nil:
		return nil, fmt.Errorf("no evaluation status found")
	case f.Properties.OriginTime == nil:
		return nil, fmt.Errorf("no origin time found")
	case f.Properties.OriginError == nil:
		return nil, fmt.Errorf("no origin error found")
	case f.Properties.UsedPhaseCount == nil:
		return nil, fmt.Errorf("no used phase count found")
	case f.Properties.UsedStationCount == nil:
		return nil, fmt.Errorf("no used station count found")
	case f.Properties.AzimuthalGap == nil:
		return nil, fmt.Errorf("no used azimuthal gap found")
	case f.Properties.MinimumDistance == nil:
		return nil, fmt.Errorf("no minimum distance gap found")
	case f.Properties.MagnitudeStationCount == nil:
		return nil, fmt.Errorf("no magnitude station count found")
	case f.Properties.EarthModel == nil:
		return nil, fmt.Errorf("no earth model found")
	case f.Properties.EvaluationMethod == nil:
		return nil, fmt.Errorf("no evaluation method found")
	case f.Properties.EvaluationMode == nil:
		return nil, fmt.Errorf("no evaluation mode found")
	}

	e := Earthquake{
		PublicID: *f.Properties.PublicID,

		UID: strings.Join([]string{
			*f.Properties.PublicID,
			*f.Properties.EvaluationStatus,
			f.Properties.OriginTime.Format(RFC3339Micro),
			strconv.FormatFloat(f.Geometry.Coordinates[1], 'f', -1, 64),
			strconv.FormatFloat(f.Geometry.Coordinates[0], 'f', -1, 64),
			strconv.FormatFloat(*f.Properties.Depth, 'f', -1, 64),
			strconv.FormatFloat(*f.Properties.Magnitude, 'f', -1, 64),
			*f.Properties.MagnitudeType,
		}, ":"),

		AgencyID: *agency,

		Type:          *f.Properties.EventType,
		Status:        *f.Properties.EvaluationStatus, // for want of an alternative
		Origin:        f.Properties.OriginTime.Format(RFC3339Micro),
		StandardError: *f.Properties.OriginError,

		Latitude:  f.Geometry.Coordinates[1],
		Longitude: f.Geometry.Coordinates[0],
		Depth:     *f.Properties.Depth,

		UsedPhaseCount:        *f.Properties.UsedPhaseCount,
		UsedStationCount:      *f.Properties.UsedStationCount,
		AzimuthalGap:          *f.Properties.AzimuthalGap,
		MinimumDistance:       *f.Properties.MinimumDistance,
		Magnitude:             *f.Properties.Magnitude,
		MagnitudeType:         *f.Properties.MagnitudeType,
		MagnitudeStationCount: *f.Properties.MagnitudeStationCount,
		EarthModelID:          *f.Properties.EarthModel,
		MethodID:              *f.Properties.EvaluationMethod,
		EvaluationMode:        *f.Properties.EvaluationMode,
		EvaluationStatus:      *f.Properties.EvaluationStatus,

		UpdateTime: f.Properties.ModificationTime.Format(RFC3339Micro),
	}

	return &e, nil
}
