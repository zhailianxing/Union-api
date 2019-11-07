package resize

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestResize(t *testing.T) {

	bs, err := ioutil.ReadFile("image.JPG")
	if err != nil {
		t.Fatal(err)
	}

	data := Resize(bs, 400, 400)

	ioutil.WriteFile("image2.JPG", data, 0644)
	fmt.Println(len(data))
}
