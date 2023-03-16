package teltonika

import (
	"bytes"
	"fmt"
	"github.com/inmatics/tcp_to_mqtt/streams"
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
			imei,
			Location{
				"Point",
				[]float64{longitude, latitude},
			},
			timestamp,
			angle,
			speed,
		}

		// IO Events Elements

		reader.Next(1) // ioEventID
		reader.Next(1) // total Elements

		stage := 1
		for stage <= 4 {
			stageElements := streams.ToInt8(reader.Next(1))

			var j int8 = 0
			for j < stageElements {
				reader.Next(1) // elementID

				switch stage {
				case 1: // One byte IO Elements
					streams.ToInt8(reader.Next(1))
				case 2: // Two byte IO Elements
					streams.ToInt16(reader.Next(2))
				case 3: // Four byte IO Elements
					streams.ToInt32(reader.Next(4))
				case 4: // Eigth byte IO Elements
					streams.ToInt64(reader.Next(8))
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
