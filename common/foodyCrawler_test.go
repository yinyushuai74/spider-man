package common

import (
	"testing"
	"fmt"
)

func TestLocation(t *testing.T) {
	resp := QueryLocationList()
	fmt.Println(resp)
}