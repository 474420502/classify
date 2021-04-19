package classify

import (
	"io"
	"log"
	"os"
	"testing"
)

func headerCompletion(data []byte) (newdata []byte, datalen int) {

	for i := 0; i < len(data); i++ {
		if data[i] == ' ' {
			continue
		} else {
			newdata = append(newdata, '.')
			newdata = append(newdata, data...)
			datalen = len(newdata)
			newdata = append(newdata, ' ')

			return
		}
	}

	log.Panic("data is nil")
	return
}

func TestParsePath(t *testing.T) {
	f, err := os.Open("./p.yaml")
	if err != nil {
		log.Panic(err)
	}

	src, err := io.ReadAll(f)
	if err != nil {
		log.Panic(err)
	}

	data, dlen := headerCompletion(src)

	var i int = 0

	for i < dlen {
		c := data[i]
		switch c {
		case ' ':
			continue
		case '[':
			continue
		case '.':
			continue
		}
	}
}
