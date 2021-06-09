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
		if passitem == nil {
			return TestStruct{}
		}
		item := passitem.(TestStruct)
		return TestStruct{
			Name:  item.Name,
			Label: item.Label,
			Value: item.Value}
	})
	streamer.SetCountHandler(func(counted, item interface{}) {
		c := counted.(TestStruct)
		c.Value += item.(TestStruct).Value
	})

	streamer.Add(TestStruct{
		Name:  "day",
		Label: "xi",
		Value: 1,
	})

	streamer.Add(TestStruct{
		Name:  "month",
		Label: "ha",
		Value: 1,
	})

	streamer.Add(TestStruct{
		Name:  "day",
		Label: "xi",
		Value: 2,
	})

	for k, v := range streamer.bytesdict {
		t.Error(k, v.(TestStruct).Value)
	}

}
