package teltonika

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/inmatics/tcp_to_mqtt/pkg/streams"
	"log"
)

func parseData(data []byte, imei string) (elements []Record, err error) {
	reader := bytes.NewBuffer(data)
	// fmt.Println("Reader Size:", reader.Len())

	// Header
	reader.Next(4)                                 // 4 Zero Bytes
	streams.ToInt32(reader.Next(4))                // Header
	reader.Next(1)                                 // CodecID
	recordNumber := streams.ToInt8(reader.Next(1)) // Number of Records

	elements = make([]Record, recordNumber)

	var i int8 = 0
	for i < recordNumber {
		timestamp := streams.ToTime(reader.Next(8)) // Timestamp
		reader.Next(1)                              // Priority

		// GPS Element
		longitudeInt := streams.ToInt32(reader.Next(4)) // Longitude
		longitude := float64(longitudeInt) / PRECISION
		latitudeInt := streams.ToInt32(reader.Next(4)) // Latitude
		latitude := float64(latitudeInt) / PRECISION

		reader.Next(2)                           // Altitude
		angle := streams.ToInt16(reader.Next(2)) // Angle
		reader.Next(1)                           // Satellites
		speed := streams.ToInt16(reader.Next(2)) // Speed

		elements[i] = Record{
			Angle:     angle,
			Direction: angle,
			Imei:      imei,
			Location: Location{
				"Point",
				[]float64{longitude, latitude},
			},
			Speed:     speed,
			Timestamp: timestamp,
		}
		// IO Events Elements

		reader.Next(1) // ioEventID
		reader.Next(1) // total Elements

		stage := 1
		for stage <= 4 {
			stageElements := streams.ToInt8(reader.Next(1))

			var j int8 = 0
			for j < stageElements {
				el := streams.ToUInt8(reader.Next(1))
				elementID, _ := streams.ParseBs2Uint16([]byte{0, 0, el}, 1)

				switch stage {
				case 1: // One byte IO Elements
					manageElementValue(elementID, reader.Next(1), &elements[i])
				case 2: // Two byte IO Elements
					manageElementValue(elementID, reader.Next(2), &elements[i])
				case 3: // Four byte IO Elements
					manageElementValue(elementID, reader.Next(4), &elements[i])
				case 4: // Eight byte IO Elements
					manageElementValue(elementID, reader.Next(8), &elements[i])
				}
				j++
			}
			stage++
		}

		if err != nil {
			fmt.Println("Error while reading IO Elements")
			break
		}

		i++
	}

	streams.ToInt8(reader.Next(1))  // Number of Records
	streams.ToInt32(reader.Next(4)) // CRC
	return
}

func manageElementValue(key uint16, value []byte, el *Record) {
	var h Decoder
	h.elements = make(map[string]map[uint16]AvlEncodeKey)
	// read our opened json as a byte array.
	byteValue := []byte(FMBXY)
	fmbxy := make(map[uint16]AvlEncodeKey)
	err := json.Unmarshal(byteValue, &fmbxy)
	if err != nil {
		log.Panic(err)
	}

	h.elements["FMBXY"] = fmbxy
	avl, _ := fmbxy[key]

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

type Decoder struct {
	elements map[string]map[uint16]AvlEncodeKey
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
