package streams

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"time"
)

func ToInt8(data []byte) int8 {
	var y int8
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &y)
	logFatal(err)
	return y
}

func ToInt16(data []byte) int16 {
	var y int16
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &y)
	logFatal(err)
	return y
}

func ToInt32(data []byte) int32 {
	var y int32
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &y)
	if y>>31 == 1 {
		y *= -1
	}
	logFatal(err)
	return y
}

func ToInt64(data []byte) int64 {
	var y int64
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &y)
	logFatal(err)
	return y
}

func ToUInt8(data []byte) uint8 {
	var y uint8
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &y)
	logFatal(err)
	return y
}

func ToTime(data []byte) time.Time {
	miliseconds := ToInt64(data)
	seconds := int64(float64(miliseconds) / 1000.0)
	nanoseconds := miliseconds % 1000
	return time.Unix(seconds, nanoseconds)
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ParseBs2Uint16(bs []byte, offset int) (uint16, error) {
	// error handling
	if len(bs) < offset+2 {
		return 0, fmt.Errorf("ParseBs2Uint16 invalid length of slice %#x , slice len %v , want %v", (bs), len(bs), offset+2)
	}
	var sum uint16
	var order uint32
	// convert hex byte slice to Uint64
	for i := offset + 1; i >= offset; i-- {
		// shift to the left by 8 bits every cycle
		sum += uint16((bs)[i]) << order
		order += 8
	}
	return sum, nil
}
