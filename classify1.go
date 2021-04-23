package classify

import (
	"fmt"

	"github.com/474420502/focus/tree/vbtkey"
)

// 分类
type Classify struct {
	Categorys []*Category
	Values    *vbtkey.Tree
}

type CategoryHandler func(value interface{}) interface{}

// Category 类别
type Category struct {
	Name      string
	Handler   CategoryHandler
	IsCollect bool
}

func New() *Classify {
	return &Classify{}
}

func (c *Classify) AddCategory(name string, handler CategoryHandler) *Classify {
	c.Categorys = append(c.Categorys, &Category{
		Name:    name,
		Handler: handler,
	})
	return c
}

func (c *Classify) Keys(paths ...interface{}) (result []interface{}) {
	var values *vbtkey.Tree = c.Values
	// var category *Category

	for _, p := range paths {
		// category = c.Categorys[i]
		if child, ok := values.Get(p); ok {
			values = child.(*vbtkey.Tree)
		} else {
			panic(fmt.Errorf("no key %v", p))
		}
	}

	// if category.IsCollect {

	// }

	values.Traversal(func(k, v interface{}) bool {
		result = append(result, k)
		return true
	})

	return
}
