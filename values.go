package classify

// values 值+
type collection struct {
	valuesdict map[interface{}]*collection
	values     []interface{}
}
