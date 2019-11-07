package xpoint

import (
	"fmt"
	"testing"
)

func TestGetAddPoint(t *testing.T) {
	var i uint16
	for i = 0; i < 10; i++ {
		fmt.Println(GetAddPoint(i))
	}
}
