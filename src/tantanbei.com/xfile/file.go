package xfile

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func ReadLine(fileName string, handler func(string)) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if err == nil {
			line = strings.TrimSpace(line)
			handler(line)
		} else if err != io.EOF {
			return err
		} else {
			return nil
		}
	}
	return nil
}

func ReadLineHandlerFileName(fileName string, handler func(string, int, string)) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	buf := bufio.NewReader(f)
	for n := 0; ; n++ {
		line, err := buf.ReadString('\n')
		if err == nil {
			line = strings.TrimSpace(line)
			handler(fileName, n, line)
		} else if err != io.EOF {
			return err
		} else {
			return nil
		}
	}
	return nil
}
