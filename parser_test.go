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

type Item struct {
	Name      string
	Region    string
	Country   string
	Coin      int64
	ExtraCoin int64
}

func TestPut(t *testing.T) {
	clsfy := New()
	clsfy.Build(`region<Region>.country<Country>.@Name,
	 `)
	clsfy.Put(&Item{
		Name:    "test",
		Region:  "Arab",
		Country: "USA",
	})
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
