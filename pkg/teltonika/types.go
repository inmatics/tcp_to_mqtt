package teltonika

import "time"

const PRECISION = 1_0000_000.0

// Record Schema
type Record struct {
	Angle      int16     `json:"angle"`
	Battery    float32   `json:"battery"`
	Direction  int16     `json:"direction"`
	Ignition   int8      `json:"ignition"`
	Imei       string    `json:"imei"`
	Lat        float64   `json:"lat" gorm:"type:decimal(9,6);"`
	Lng        float64   `json:"lng" gorm:"type:decimal(9,6);"`
	Odometer   int32     `json:"odometer,omitempty"`
	Speed      int16     `json:"speed"`
	Timestamp  time.Time `json:"timestamp"`
	RawMessage string    `json:"raw_message"`
}
