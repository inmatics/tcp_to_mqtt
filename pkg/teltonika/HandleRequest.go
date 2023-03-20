package teltonika

import (
	"encoding/hex"
	"encoding/json"
	"golang.org/x/exp/slog"
	"io"
	"net"
	"os"
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
				bytes, err := json.Marshal(element)
				if err != nil {
					return
				}
				messages <- string(bytes)

				logger.Debug("Parsed data",
					slog.String("parsed data", string(bytes)),
				)
			}

			_, err = conn.Write([]byte{0, 0, 0, uint8(len(elements))})
			if err != nil {
				logger.Error("Error writing data", err)
			}
		}

	}
}
