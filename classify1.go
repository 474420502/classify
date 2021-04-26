package classify

// import (
// 	"fmt"
// 	"log"
// 	"reflect"
// 	"strings"
// 	"time"

// 	"github.com/474420502/focus/tree/vbtkey"
// )

// // kingTime 时间
// var kingTime = reflect.TypeOf(time.Time{}).Kind()

// var defaultCategoryHandler CategoryHandler = func(value interface{}) interface{} {
// 	return 0
// }

// // 分类
// type Classify struct {
// 	categorys []*Category
// 	Values    *vbtkey.Tree
// }

// type CategoryHandler func(value interface{}) interface{}

// // Category 类别
// type Category struct {
// 	Name      string
// 	Handler   CategoryHandler
// 	IsCollect bool
// }

// type CData struct {
// 	IsCollect bool
// 	Values    *vbtkey.Tree
// }

// func New() *Classify {
// 	return &Classify{}
// }

// func (c *Classify) AddCategory(name string, handler CategoryHandler) *Classify {
// 	c.categorys = append(c.categorys, &Category{
// 		Name:      name,
// 		Handler:   handler,
// 		IsCollect: false,
// 	})
// 	return c
// }
// func (c *Classify) Categorys() string {
// 	if len(c.categorys) == 0 {
// 		return ""
// 	}
// 	var content []byte
// 	for _, cate := range c.categorys {
// 		content = append(content, []byte(cate.Name)...)
// 		if !cate.IsCollect {
// 			content = append(content, '.')
// 		}
// 	}

// 	if content[len(content)-1] == '.' {
// 		return string(content[:len(content)-1])
// 	}
// 	return string(content)
// }

// func (c *Classify) Collect() {
// 	c.categorys = append(c.categorys, &Category{
// 		Name:      "@",
// 		Handler:   defaultCategoryHandler,
// 		IsCollect: true,
// 	})
// }

// func (c *Classify) CollectCategory(handler CategoryHandler) {
// 	c.categorys = append(c.categorys, &Category{
// 		Name:      "@",
// 		Handler:   handler,
// 		IsCollect: true,
// 	})
// }

// func (c *Classify) Keys(paths ...interface{}) (result []interface{}) {
// 	var values *vbtkey.Tree = c.Values
// 	// var category *Category
// 	if len(paths) >= len(c.categorys)-1 {
// 		panic(fmt.Sprintf("categorys len is %d only: %#v", len(c.categorys)-1, c.Categorys()))
// 	}
// 	for _, p := range paths {
// 		// category = c.Categorys[i]
// 		if child, ok := values.Get(p); ok {
// 			values = child.(*vbtkey.Tree)
// 		} else {
// 			panic(fmt.Errorf("no key %v", p))
// 		}
// 	}

// 	values.Traversal(func(k, v interface{}) bool {
// 		result = append(result, k)
// 		return true
// 	})

// 	return
// }

// func (c *Classify) Put(v interface{}) {
// 	if c.Values == nil {
// 		c.Values = vbtkey.New(autoComapre)
// 	}
// 	put(c.categorys, 0, c.Values, v)
// }

// func put(categorys []*Category, cidx int, Values *vbtkey.Tree, v interface{}) {
// 	cate := categorys[cidx]
// 	if cate.IsCollect {
// 		Values.Put(cate.Handler(v), v)
// 		return
// 	} else {
// 		// 判断Values是否存在
// 		var NextValues *vbtkey.Tree
// 		key := cate.Handler(v)
// 		if vs, ok := Values.Get(key); ok {
// 			NextValues = vs.(*vbtkey.Tree)
// 		} else {
// 			NextValues = vbtkey.New(autoComapre)
// 			Values.Put(key, NextValues)
// 		}
// 		put(categorys, cidx+1, NextValues, v)
// 	}
// }

// func (c *Classify) Get(out interface{}, vPaths ...interface{}) {

// 	var values *vbtkey.Tree = c.Values
// 	if len(vPaths) >= len(c.categorys) {
// 		panic(fmt.Sprintf("values keys deepth is %d only: %#v", len(c.categorys), c.Categorys()))
// 	}

// 	if reflect.TypeOf(out).Kind() != reflect.Ptr {
// 		panic("out must ptr")
// 	}

// 	outv := reflect.ValueOf(out)
// 	result := outv.Elem()
// 	result = reflect.MakeSlice(result.Type(), 0, 0)
// 	// log.Println(result)

// 	var cidx = 0
// 	for ; cidx < len(vPaths); cidx++ {
// 		vp := vPaths[cidx]

// 		if child, ok := values.Get(vp); ok {
// 			values = child.(*vbtkey.Tree)
// 		}
// 	}

// 	var getValues func(cidx int, values *vbtkey.Tree)
// 	getValues = func(cidx int, values *vbtkey.Tree) {
// 		category := c.categorys[cidx]
// 		if category.IsCollect {
// 			values.Traversal(func(k, v interface{}) bool {
// 				result = reflect.Append(result, reflect.ValueOf(v))
// 				return true
// 			})
// 		} else {
// 			values.Traversal(func(k, v interface{}) bool {
// 				getValues(cidx+1, v.(*vbtkey.Tree))
// 				return true
// 			})
// 		}
// 	}
// 	getValues(cidx, values)
// 	outv.Elem().Set(result)
// }

// func autoComapre(k1, k2 interface{}) int {

// 	t1 := reflect.TypeOf(k1)
// 	t2 := reflect.TypeOf(k2)

// 	if t1.Kind() != t2.Kind() {
// 		log.Panicf("value1 %v, value2 %v is not same type. please check keys input: v1:%v v2:%v", t1.Kind(), t2.Kind(), k1, k2)
// 	}

// 	rv1 := reflect.ValueOf(k1)
// 	rv2 := reflect.ValueOf(k2)

// 	if t1.Kind() == reflect.Ptr {
// 		t1 = t1.Elem()
// 		t2 = t2.Elem()
// 		rv1 = rv1.Elem()
// 		rv2 = rv2.Elem()
// 	}

// 	switch t1.Kind() {

// 	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
// 		v1 := rv1.Int()
// 		v2 := rv2.Int()
// 		switch {
// 		case v1 > v2:
// 			return 1
// 		case v1 < v2:
// 			return -1
// 		default:
// 			return 0
// 		}
// 	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
// 		v1 := rv1.Uint()
// 		v2 := rv2.Uint()
// 		switch {
// 		case v1 > v2:
// 			return 1
// 		case v1 < v2:
// 			return -1
// 		default:
// 			return 0
// 		}
// 	case reflect.Float32, reflect.Float64:
// 		v1 := rv1.Float()
// 		v2 := rv2.Float()
// 		switch {
// 		case v1 > v2:
// 			return 1
// 		case v1 < v2:
// 			return -1
// 		default:
// 			return 0
// 		}
// 	case reflect.String:
// 		v1 := rv1.String()
// 		v2 := rv2.String()
// 		return strings.Compare(v1, v2)
// 	case kingTime:
// 		v1 := rv1.Interface().(time.Time).UnixNano()
// 		v2 := rv2.Interface().(time.Time).UnixNano()
// 		switch {
// 		case v1 > v2:
// 			return 1
// 		case v1 < v2:
// 			return -1
// 		default:
// 			return 0
// 		}
// 	default:

// 		panic(fmt.Sprintf("%v kind not handled", t1.Kind()))
// 	}

// }
