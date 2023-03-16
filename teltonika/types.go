package teltonika

import "time"

const PRECISION = 1_0000_000.0

// Struct for Mongo GeoJSON
type Location struct {
	Type        string
	Coordinates []float64
}

// Record Schema
type Record struct {
	Imei     string
	Location Location
	Time     time.Time
	Angle    int16
	Speed    int16
}
