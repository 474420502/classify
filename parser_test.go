package classify

import (
	"io"
	"log"
	"os"
	"testing"
)

func init() {
	log.SetFlags(log.Llongfile)
}

func TestParserCPath(t *testing.T) {
	f, err := os.Open("./test1.cpath")
	if err != nil {
		log.Panic(err)
	}

	src, err := io.ReadAll(f)
	if err != nil {
		log.Panic(err)
	}

	var parent *category = newCategory()

	cur := headerCompletion(src)

	// First paragraph 标签名检测
	extract(parent, cur)

	log.Println(parent)
	log.Println("")
}
