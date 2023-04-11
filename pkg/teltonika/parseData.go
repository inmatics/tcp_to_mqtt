package teltonika

import (
	"bytes"
	"encoding/json"
	"github.com/inmatics/tcp_to_mqtt/pkg/streams"
	"log"
)

func parseData(data []byte, imei string) ([]Record, error) {
	reader := bytes.NewBuffer(data)

	// Header
	reader.Next(4)                                 // 4 Zero Bytes
	reader.Next(4)                                 // Header
	reader.Next(1)                                 // CodecID
	recordNumber := streams.ToInt8(reader.Next(1)) // Number of Records

	elements := make([]Record, recordNumber)

	for i := 0; i < int(recordNumber); i++ {
		timestamp := streams.ToTime(reader.Next(8)) // Timestamp
		reader.Next(1)                              // Priority

		// GPS Element
		longitude := float64(streams.ToInt32(reader.Next(4))) / PRECISION // Longitude
		latitude := float64(streams.ToInt32(reader.Next(4))) / PRECISION  // Latitude

		reader.Next(2)                           // Altitude
		angle := streams.ToInt16(reader.Next(2)) // Angle
		reader.Next(1)                           // Satellites
		speed := streams.ToInt16(reader.Next(2)) // Speed

		elements[i] = Record{
			Angle:     angle,
			Direction: angle,
			Imei:      imei,
			Lat:       latitude,
			Lng:       longitude,
			Speed:     speed,
			Timestamp: timestamp,
		}

		// IO Events Elements
		reader.Next(1) // ioEventID
		reader.Next(1) // total Elements

		for stage := 1; stage <= 4; stage++ {
			for j := 0; j < int(streams.ToInt8(reader.Next(1))); j++ {
				elementID, _ := streams.ParseBs2Uint16([]byte{0, 0, streams.ToUInt8(reader.Next(1))}, 1)
				decoder, err := NewDecoder(FMBXY)
				if err != nil {
					return nil, err
				}
				decoder.manageElementValue(elementID, reader.Next(1<<(stage-1)), &elements[i])
			}
		}
	}

	reader.Next(1) // Number of Records
	reader.Next(4) // CRC
	return elements, nil
}

type Decoder struct {
	elements map[uint16]AvlEncodeKey
}

type AvlEncodeKey struct {
	Bytes           string
	Description     string
	FinalConversion string
	HWSupport       string
	Max             string
	Min             string
	Multiplier      string
	No              string
	ParametrGroup   string
	PropertyName    string
	Type            string
	Units           string
}

func NewDecoder(fmbxyJSON string) (*Decoder, error) {
	fmbxy := make(map[uint16]AvlEncodeKey)
	err := json.Unmarshal([]byte(fmbxyJSON), &fmbxy)
	if err != nil {
		return nil, err
	}

	return &Decoder{
		elements: map[uint16]AvlEncodeKey{},
	}, nil
}

func (d *Decoder) manageElementValue(key uint16, value []byte, el *Record) {
	avl, ok := d.elements[key]
	if !ok {
		log.Printf("Key not found: %d", key)
		return
	}

	switch avl.FinalConversion {
	case "toUint8":
		if avl.PropertyName == "Ignition" {
			el.Ignition = streams.ToInt8(value)
		}
	case "toUint16":
		if avl.PropertyName == "External Voltage" {
			el.Battery = float32(streams.ToInt16(value)) / 1000
		}
	case "toUint32":
		if avl.PropertyName == "Total Odometer" {
			el.Odometer = streams.ToInt32(value)
		}
	}
}
