package classify

import (
	"testing"
)

func Test1(t *testing.T) {
	var m map[interface{}]interface{} = make(map[interface{}]interface{})
	m[1] = 1
	m["123"] = 1
	m[1] = 2
	m["asd"] = 3
	m["1234"] = 5
	t.Errorf("%#v", m)
}
