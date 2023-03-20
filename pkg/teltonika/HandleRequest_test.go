package teltonika

import (
	"fmt"
	"testing"
)

func Test_toString(t *testing.T) {

	record := Record{
		Imei: "356307042441013",
	}
	fmt.Println(toString(record))
}
