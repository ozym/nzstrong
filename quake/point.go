package quake

import (
	"fmt"
	"math"
)

type Point struct {
	Latitude  float64
	Longitude float64
}

const RadiansToDegrees = 57.2957795
const DegreesToKm = math.Pi * 6371.0 / 180.0
const RadiansToKm = RadiansToDegrees * DegreesToKm

func (p Point) Distance(point Point) float64 {

	if (p.Latitude == point.Latitude) && (p.Longitude == point.Longitude) {
		return 0.0
	}

	esq := (1.0 - 1.0/298.25) * (1.0 - 1.0/298.25)
	alat3 := math.Atan(math.Tan(p.Latitude/RadiansToDegrees)*esq) * RadiansToDegrees
	alat4 := math.Atan(math.Tan(point.Latitude/RadiansToDegrees)*esq) * RadiansToDegrees

	rlat1 := alat3 / RadiansToDegrees
	rlat2 := alat4 / RadiansToDegrees
	rdlon := (point.Longitude - p.Longitude) / RadiansToDegrees

	clat1 := math.Cos(rlat1)
	clat2 := math.Cos(rlat2)
	slat1 := math.Sin(rlat1)
	slat2 := math.Sin(rlat2)
	cdlon := math.Cos(rdlon)

	cdel := slat1*slat2 + clat1*clat2*cdlon
	switch {
	case cdel > 1.0:
		cdel = 1.0
	case cdel < -1.0:
		cdel = -1.0
	}

	return RadiansToKm * math.Acos(cdel)
}

func (p Point) Azimuth(point Point) float64 {

	if (p.Latitude == point.Latitude) && (p.Longitude == point.Longitude) {
		return 0.0
	}

	esq := (1.0 - 1.0/298.25) * (1.0 - 1.0/298.25)
	alat3 := math.Atan(math.Tan(p.Latitude/RadiansToDegrees)*esq) * RadiansToDegrees
	alat4 := math.Atan(math.Tan(point.Latitude/RadiansToDegrees)*esq) * RadiansToDegrees

	rlat1 := alat3 / RadiansToDegrees
	rlat2 := alat4 / RadiansToDegrees
	rdlon := (point.Longitude - p.Longitude) / RadiansToDegrees

	clat1 := math.Cos(rlat1)
	clat2 := math.Cos(rlat2)
	slat1 := math.Sin(rlat1)
	slat2 := math.Sin(rlat2)
	cdlon := math.Cos(rdlon)
	sdlon := math.Sin(rdlon)

	yazi := sdlon * clat2
	xazi := clat1*slat2 - slat1*clat2*cdlon

	azi := RadiansToDegrees * math.Atan2(yazi, xazi)

	if azi < 0.0 {
		azi += 360.0
	}

	return azi
}

func (p Point) BackAzimuth(point Point) float64 {

	if (p.Latitude == point.Latitude) && (p.Longitude == point.Longitude) {
		return 0.0
	}

	esq := (1.0 - 1.0/298.25) * (1.0 - 1.0/298.25)
	alat3 := math.Atan(math.Tan(p.Latitude/RadiansToDegrees)*esq) * RadiansToDegrees
	alat4 := math.Atan(math.Tan(point.Latitude/RadiansToDegrees)*esq) * RadiansToDegrees

	rlat1 := alat3 / RadiansToDegrees
	rlat2 := alat4 / RadiansToDegrees
	rdlon := (point.Longitude - p.Longitude) / RadiansToDegrees

	clat1 := math.Cos(rlat1)
	clat2 := math.Cos(rlat2)
	slat1 := math.Sin(rlat1)
	slat2 := math.Sin(rlat2)
	cdlon := math.Cos(rdlon)
	sdlon := math.Sin(rdlon)

	ybaz := -sdlon * clat1
	xbaz := clat2*slat1 - slat2*clat1*cdlon

	baz := RadiansToDegrees * math.Atan2(ybaz, xbaz)

	if baz < 0.0 {
		baz += 360.0
	}

	return baz
}

func (p Point) Compass(point Point) string {
	azimuth := p.Azimuth(point) + 22.5

	for azimuth < 0.0 {
		azimuth += 360.0
	}
	for azimuth >= 360.0 {
		azimuth -= 360.0
	}

	switch int(math.Floor(azimuth / 45.0)) {
	case 0:
		return "north"
	case 1:
		return "north-east"
	case 2:
		return "east"
	case 3:
		return "south-east"
	case 4:
		return "south"
	case 5:
		return "south-west"
	case 6:
		return "west"
	default:
		return "north-west"
	}
}

func Bearing(azimuth float64) string {

	for azimuth < 0.0 {
		azimuth += 360.0
	}
	for azimuth >= 360.0 {
		azimuth -= 360.0
	}

	switch int(math.Floor(azimuth / 90.0)) {
	case 0:
		return fmt.Sprintf("N%02.0fE", azimuth)
	case 1:
		return fmt.Sprintf("S%02.0fE", 180.0-azimuth)
	case 2:
		return fmt.Sprintf("S%02.0fW", azimuth-180.0)
	default:
		return fmt.Sprintf("N%02.0fW", 360.0-azimuth)
	}
}
