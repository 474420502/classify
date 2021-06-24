package classify

import (
	"testing"
)

type TestStruct struct {
	Name  string
	Label string
	Value int
}

func TestMode(t *testing.T) {

	streamer := NewStreamer("<Name>.<Label>")
	streamer.SetCreateCountedHandler(func(passitem interface{}) interface{} {
		item := passitem.(*TestStruct)
		return &TestStruct{
			Name:  item.Name,
			Label: item.Label,
			Value: item.Value}
	})
	streamer.SetCountHandler(func(counted, item interface{}) {
		c := counted.(*TestStruct)
		c.Value += item.(*TestStruct).Value
	})

	streamer.Add(&TestStruct{
		Name:  "day",
		Label: "xi",
		Value: 1,
	})

	streamer.Add(&TestStruct{
		Name:  "month",
		Label: "ha",
		Value: 1,
	})

	streamer.Add(&TestStruct{
		Name:  "day",
		Label: "xi",
		Value: 2,
	})

	streamer.RangeItems(func(item interface{}) bool {
		t.Error(item)
		return true
	})

}
