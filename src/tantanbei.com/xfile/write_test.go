package xfile

import (
	"testing"
)

func TestWriteFile(t *testing.T) {
	WriteFile("./write.txt", []byte("tankaide"))
}

func TestWriteJsonFile(t *testing.T) {
	WriteJsonFile("./write_json.txt", []int{1, 2, 3, 4, 5, 67, 9, 6, 54564, 3, 56, 3, 6, 4})
}
