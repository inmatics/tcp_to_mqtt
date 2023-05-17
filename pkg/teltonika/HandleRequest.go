package teltonika

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

	"golang.org/x/exp/slog"
)

func HandleRequest(conn net.Conn, recordsCannel chan Record, logger *slog.Logger) {
	var b []byte
	var imei string
	step := 1
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			logger.Error("error closing connection", err)
		}
		logger.Debug("connection closed successfully")
	}(conn)

	for {
		buf := make([]byte, 2048)

		size, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			// https://github.com/golang/go/discussions/54763
			logger.Error("error reading TCP connection", err.Error(), imei)
			return
		}

		b = []byte{1} // 0x01 means we accept the message

		message := hex.EncodeToString(buf[:size])

		logPacket(conn, logger, size, message)

		switch step {
		case 1:
			imeiLength, _ := strconv.ParseInt(message[0:4], 16, 64)
			if imeiLength != 15 && imeiLength != 16 {
				log.Println("Error decoding IMEI length (length must be 15 or 16)", imeiLength)
				declineMessage(conn)
				break
			}
			imei, err = ParseIMEI(buf[2:size], int(imeiLength))

			if err != nil {
				log.Println("Error decoding IMEI:", message)
				declineMessage(conn)
				break
			} else {
				_, err := conn.Write(b)
				if err != nil {
					logger.Error("Error writing data", err)
				}
				logger.Debug("New accepted connection",
					slog.String("IMEI", imei),
				)
				step = 2
			}

		case 2:
			elements, err := parseData(buf, imei)
			if err != nil {
				logger.Error("Error while parsing data", err)
				break
			}

			for i := 0; i < len(elements); i++ {
				element := elements[i]
				element.RawMessage = message
				if err != nil {
					return
				}
				recordsCannel <- element

			}

			_, err = conn.Write([]byte{0, 0, 0, uint8(len(elements))})
			if err != nil {
				logger.Error("Error writing data", err)
			}
		}

	}
}

func logPacket(conn net.Conn, logger *slog.Logger, size int, message string) {
	logger.Debug("new TCP packet",
		slog.String("remote address", conn.RemoteAddr().String()),
		slog.String("local address", conn.LocalAddr().String()),
		slog.Int("step", 2),
		slog.Group("message", slog.Int("size", size), slog.String("message", message)),
	)
}

func declineMessage(conn net.Conn) {
	// 0x00 we decline the message
	_, err := conn.Write([]byte{0})
	if err != nil {
		log.Println("Error declining message")
	}
}

// ParseIMEI takes a pointer to a byte slice including IMEI number encoded as ASCII, IMEI length, offset and returns IMEI as string and error. If len is 15 chars, also do imei validation
func ParseIMEI(bs []byte, length int) (string, error) {
	// error handling
	if len(bs) < 15 {
		return "", fmt.Errorf("ParseIMEI invalid length of slice %#x , slice len %v , want %v", (bs), len(bs), 8)
	}
	// range over slice
	x := string((bs)[:length])

	if len(x) == 15 {
		if !ValidateIMEI(&x) {
			return "", fmt.Errorf("IMEI %v is invalid", x)
		}
	}

	return x, nil
}

// ValidateIMEI takes pointer to 15 digits long IMEI string, calculate checksum using the Luhn algorithm and return validity as bool
func ValidateIMEI(imei *string) bool {
	bs := []byte((*imei))

	if len(bs) != 15 {
		// log.Printf("Should validate only 15chars long Imei, got %v", len(bs))
		return false
	}

	parsed, err := strconv.ParseInt(string(bs[len(bs)-1]), 10, 8)
	if err != nil {
		// log.Printf("Unable to parse IMEI digits %v", err)
		return false
	}
	checkSumDigit := int8(parsed)
	var checkSum uint64

	// make buffer array for Luhn algorithm with len 14 bytes and cap 31 bytes
	digits := make([]uint8, 14, 31)

	// count Luhn algorithm
	for i := 0; i < 14; i++ {

		parsed, err = strconv.ParseInt(string(bs[i]), 10, 8)
		if err != nil {
			// log.Printf("Unable to parse IMEI digits %v", err)
			return false
		}

		digits[i] = uint8(parsed)
		if i%2 != 0 {
			digits[i] = digits[i] * 2
		}

		if digits[i] >= 10 {
			digits = append(digits, 1)
			digits[i] = digits[i] % 10
		}
	}

	for _, val := range digits {
		checkSum += uint64(val)
	}

	// when checkSum is 0, should use 0
	if checkSumDigit == 0 {
		return 0 == uint64(checkSumDigit)
	}

	// return true if divider to 10 is same as the checkSumDigit
	return (10 - checkSum%10) == uint64(checkSumDigit)
}
