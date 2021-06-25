package classify

import (
	"testing"

	"github.com/Pallinder/go-randomdata"
)

type TestStruct struct {
	Name  string
	Label string
	Type  int
	Value int
}

func TestMode(t *testing.T) {

	streamer := NewStreamer("<Name>.<Type>")
	streamer.SetCreateCountedHandler(func(passitem interface{}) interface{} {
		item := passitem.(*TestStruct)
		return &TestStruct{
			Name:  item.Name,
			Label: item.Label,
			Value: item.Value,
			Type:  item.Type,
		}
	})
	streamer.SetCountHandler(func(counted, item interface{}) {
		c := counted.(*TestStruct)
		c.Value += item.(*TestStruct).Value
	})

	end := 10
	for i := 0; i < 10000; i++ {

		streamer.Add(&TestStruct{
			Name:  "day",
			Label: "xi",
			Type:  randomdata.Number(0, end),
			Value: 1,
		})

		streamer.Add(&TestStruct{
			Name:  "month",
			Label: "ha",
			Type:  randomdata.Number(0, end),
			Value: 1,
		})

		streamer.Add(&TestStruct{
			Name:  "day",
			Label: "xi",
			Type:  randomdata.Number(0, end),
			Value: 2,
		})
	}

	streamer.RangeItems(func(key string, item interface{}) bool {
		t.Error(item.(*TestStruct).Name, item.(*TestStruct).Type)
		return true
	})

}
