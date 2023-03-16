package teltonika

import (
	"fmt"
	"testing"
)

func Test_toString(t *testing.T) {

	record := Record{
		Imei: "356307042441013",
		Location: Location{
			Type:        "",
			Coordinates: []float64{0, 0},
		},
	}
	fmt.Println(toString(record))
}
