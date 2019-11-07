package convert

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestConvertImage(t *testing.T) {
	bs, err := ioutil.ReadFile("logo.png")
	if err != nil {
		t.Fatal(err)
	}

	result := ConvertImage(bs, OUTPUT_WEBP)

	file, err := os.OpenFile("logo.webp", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		t.Fatal(err)
	}
	_, err = file.Write(result)
	if err != nil {
		t.Fatal(err)
	}
}
