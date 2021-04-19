package classify

type CollectHandler func(value interface{}) interface{}

type Category struct {
	Name    string
	Collect CollectHandler
	IsEnd   bool
}

type Classify struct {
	Category Category
	Next     map[string]*Classify
	Values   []interface{}
}

func (csf *Classify) BuildClassifyMode(path string) {

}
