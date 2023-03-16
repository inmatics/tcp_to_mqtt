package teltonika

import (
	"encoding/hex"
	"golang.org/x/exp/slog"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func HandleRequest(conn net.Conn, messages chan string) {
	var b []byte
	var imei string
	step := 1
	textHandler := slog.NewTextHandler(os.Stdout)
	logger := slog.New(textHandler)
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			logger.Error("error closing connection", err)
		}
		logger.Info("connection closed successfully")
	}(conn)

	for {
		buf := make([]byte, 2048)

		size, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			logger.Error("error reading TCP connection", err)
			return
		}

		b = []byte{1} // 0x01 means we accept the message

		message := hex.EncodeToString(buf[:size])

		logger.Debug("new TCP packet",
			slog.String("remote address", conn.RemoteAddr().String()),
			slog.String("local address", conn.LocalAddr().String()),
			slog.Int("step", 2),
			slog.Group("message", slog.Int("size", size), slog.String("message", message)),
		)

		switch step {
		case 1:
			step = 2
			imei = message
			_, err := conn.Write(b)
			if err != nil {
				logger.Error("Error writing data", err)
			}
			logger.Debug("New accepted connection",
				slog.String("IMEI", imei),
			)
		case 2:
			elements, err := parseData(buf, imei)
			if err != nil {
				logger.Error("Error while parsing data", err)
				break
			}

			for i := 0; i < len(elements); i++ {
				element := elements[i]
				t := toString(element)
				messages <- t

				logger.Debug("Parsed data",
					slog.String("parsed data", t),
				)
			}

			_, err = conn.Write([]byte{0, 0, 0, uint8(len(elements))})
			if err != nil {
				logger.Error("Error writing data", err)
			}
		}

	}
}

func toString(element Record) string {
	var elements []string
	elements = append(elements, element.Imei)
	elements = append(elements, element.Time.Format(time.RFC3339))
	coords := element.Location.Coordinates
	elements = append(elements, strconv.FormatFloat(coords[0], 'f', -1, 64))
	elements = append(elements, strconv.FormatFloat(coords[1], 'f', -1, 64))
	elements = append(elements, strconv.Itoa(int(element.Speed)))
	return strings.Join(elements, ";")
}