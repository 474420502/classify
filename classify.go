package classify

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	treequeue "github.com/474420502/structure/queue/priority"
)

// kingTime 时间
var kingTime = reflect.TypeOf(time.Time{}).Kind()

var defaultCategoryHandler CategoryHandler = func(value interface{}) interface{} {
	return 0
}

// 分类
type Classify struct {
	categorys []*hCategory
	values    *treequeue.Queue
}

// CategoryHandler 处理结构体字段的返回值
type CategoryHandler func(item interface{}) interface{}

// hCategory 类别
type hCategory struct {
	Name      string
	Handler   CategoryHandler
	IsCollect bool
}

// New 创建一个分类器
func New() *Classify {
	return &Classify{}
}

// NewWithMode 先New 后 Build的一体处理过程
func NewWithMode(mode string, handlers ...CategoryHandler) *Classify {
	c := New()
	c.Build(mode, handlers...)
	return c
}

// Build 根据简单的表达式分类. 对结构体的处理
func (clsfy *Classify) Build(mode string, handlers ...CategoryHandler) {

	for _, token := range bytes.Split([]byte(mode), []byte{'.'}) {
		var i = 0

		var cname []byte
		var cmethod []byte
		var methodType int = 0 // 1 是字段( 2 是方法< 3 结尾收集@

	CNAME:
		for ; i < len(token); i++ {
			c := token[i]
			switch c {
			case '@':
				methodType = 3
				i++
				break CNAME
			case '[':
				methodType = 2
				i++
				break CNAME
			case '<':
				methodType = 1
				i++
				break CNAME
			case ' ':
				continue
			default:
				cname = append(cname, c)
			}
		}

		if methodType == 0 {
			panic(fmt.Errorf("语法错误: %s", mode))
		}

	CMETHOD:
		for ; i < len(token); i++ {
			c := token[i]
			switch c {
			case ' ':
				continue
			case ']', '>':
				break CMETHOD
			default:
				cmethod = append(cmethod, c)
			}
		}

		switch methodType {
		case 1:
			// log.Println(string(cname), string(cmethod))
			clsfy.AddCategory(string(cname), func(value interface{}) interface{} {
				v := reflect.ValueOf(value)
				if v.Type().Kind() == reflect.Ptr {
					v = v.Elem()
				}
				return v.FieldByName(string(cmethod)).Interface()
			})
		case 2:
			// log.Println(string(cname), string(cmethod))
			fidx, err := strconv.Atoi(string(cmethod))
			if err != nil {
				panic(err)
			}
			clsfy.AddCategory(string(cname), handlers[fidx])
		case 3:
			if len(cmethod) == 0 {
				clsfy.Collect()
			} else {
				clsfy.CollectCategory(func(value interface{}) interface{} {
					v := reflect.ValueOf(value)
					if v.Type().Kind() == reflect.Ptr {
						v = v.Elem()
					}
					return v.FieldByName(string(cmethod)).Interface()
				})
			}

		default:
			panic("?")
		}

	}
}

// AddCategory 添加对每个分类的处理过程
func (clsfy *Classify) AddCategory(name string, handler CategoryHandler) *Classify {
	clsfy.categorys = append(clsfy.categorys, &hCategory{
		Name:      name,
		Handler:   handler,
		IsCollect: false,
	})
	return clsfy
}

// Categorys 分类后的类集合
func (clsfy *Classify) Categorys() string {
	if len(clsfy.categorys) == 0 {
		return ""
	}
	var content []byte
	for _, cate := range clsfy.categorys {
		content = append(content, []byte(cate.Name)...)
		if !cate.IsCollect {
			content = append(content, '.')
		}
	}

	if content[len(content)-1] == '.' {
		return string(content[:len(content)-1])
	}
	return string(content)
}

// Collect 默认的CollectCategory处理.
func (clsfy *Classify) Collect() {
	clsfy.categorys = append(clsfy.categorys, &hCategory{
		Name:      "@",
		Handler:   defaultCategoryHandler,
		IsCollect: true,
	})
}

// CollectCategory 收集数据的方法设置. 就是到最后item被添加到树过程处理
func (clsfy *Classify) CollectCategory(handler CategoryHandler) {
	clsfy.categorys = append(clsfy.categorys, &hCategory{
		Name:      "@",
		Handler:   handler,
		IsCollect: true,
	})
}

// Keys 获取所有分类的Key. 默认就返回第一分类keys
func (clsfy *Classify) Keys(paths ...interface{}) (keys []interface{}) {
	var values *treequeue.Queue = clsfy.values
	// var category *Category
	if len(paths) >= len(clsfy.categorys)-1 {
		panic(fmt.Sprintf("categorys len is %d only: %#v", len(clsfy.categorys)-1, clsfy.Categorys()))
	}
	for _, p := range paths {
		// category = clsfy.Categorys[i]
		if child := values.Get(p); child != nil {
			values = child.Value().(*treequeue.Queue)
		} else {
			log.Panic(ErrGetKeyNotExists, fmt.Errorf("no key %v", p))
		}
	}

	values.Traverse(func(s *treequeue.Slice) bool {
		keys = append(keys, s.Key())
		return true
	})

	return
}

// Put 把需要分类的数据加进分类器.
func (clsfy *Classify) Put(v interface{}) {
	if clsfy.values == nil {
		clsfy.values = treequeue.New(autoComapre)
	}
	put(clsfy.categorys, 0, clsfy.values, v)
}

func (clsfy *Classify) PutSlice(items interface{}) {

	if clsfy.values == nil {
		clsfy.values = treequeue.New(autoComapre)
	}
	vitems := reflect.ValueOf(items)
	if vitems.Type().Kind() != reflect.Slice {
		panic(" input must slice ")
	}
	for i := 0; i < vitems.Len(); i++ {
		put(clsfy.categorys, 0, clsfy.values, vitems.Index(i).Interface())
	}

}

func put(categorys []*hCategory, cidx int, Values *treequeue.Queue, v interface{}) {
	cate := categorys[cidx]
	if cate.IsCollect {
		Values.Put(cate.Handler(v), v)
		return
	} else {
		// 判断Values是否存在
		var NextValues *treequeue.Queue
		key := cate.Handler(v)
		if vs := Values.Get(key); vs != nil {
			NextValues = vs.Value().(*treequeue.Queue)
		} else {
			NextValues = treequeue.New(autoComapre)
			Values.Put(key, NextValues)
		}
		put(categorys, cidx+1, NextValues, v)
	}
}

// Get 获取 vpaths(keys) 的值. 如果存在就返回nil. 否则返回 ErrGetKeyNotExists
func (clsfy *Classify) Get(out interface{}, vPaths ...interface{}) error {

	var values *treequeue.Queue = clsfy.values
	if len(vPaths) >= len(clsfy.categorys) {
		panic(fmt.Sprintf("values keys deepth is %d only: %#v", len(clsfy.categorys), clsfy.Categorys()))
	}

	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		panic("out must ptr")
	}

	outv := reflect.ValueOf(out)
	result := outv.Elem()
	result = reflect.MakeSlice(result.Type(), 0, 0)
	// log.Println(result)

	var cidx = 0
	for ; cidx < len(vPaths); cidx++ {
		vp := vPaths[cidx]

		if child := values.Get(vp); child != nil {
			values = child.Value().(*treequeue.Queue)
		} else {
			return ErrGetKeyNotExists
		}
	}

	var getValues func(cidx int, values *treequeue.Queue)
	getValues = func(cidx int, values *treequeue.Queue) {
		category := clsfy.categorys[cidx]
		if category.IsCollect {
			values.Traverse(func(s *treequeue.Slice) bool {
				result = reflect.Append(result, reflect.ValueOf(s.Value()))
				return true
			})
		} else {
			values.Traverse(func(s *treequeue.Slice) bool {
				getValues(cidx+1, s.Value().(*treequeue.Queue))
				return true
			})
		}
	}
	getValues(cidx, values)
	outv.Elem().Set(result)
	return nil
}

func autoComapre(k1, k2 interface{}) int {

	t1 := reflect.TypeOf(k1)
	t2 := reflect.TypeOf(k2)

	if t1.Kind() != t2.Kind() {
		log.Panicf("value1 %v, value2 %v is not same type. please check keys input: v1:%v v2:%v", t1.Kind(), t2.Kind(), k1, k2)
	}

	rv1 := reflect.ValueOf(k1)
	rv2 := reflect.ValueOf(k2)

	if t1.Kind() == reflect.Ptr {
		t1 = t1.Elem()
		t2 = t2.Elem()
		rv1 = rv1.Elem()
		rv2 = rv2.Elem()
	}

	switch t1.Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v1 := rv1.Int()
		v2 := rv2.Int()
		switch {
		case v1 > v2:
			return 1
		case v1 < v2:
			return -1
		default:
			return 0
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v1 := rv1.Uint()
		v2 := rv2.Uint()
		switch {
		case v1 > v2:
			return 1
		case v1 < v2:
			return -1
		default:
			return 0
		}
	case reflect.Float32, reflect.Float64:
		v1 := rv1.Float()
		v2 := rv2.Float()
		switch {
		case v1 > v2:
			return 1
		case v1 < v2:
			return -1
		default:
			return 0
		}
	case reflect.String:
		v1 := rv1.String()
		v2 := rv2.String()
		return strings.Compare(v1, v2)
	case kingTime:
		v1 := rv1.Interface().(time.Time).UnixNano()
		v2 := rv2.Interface().(time.Time).UnixNano()
		switch {
		case v1 > v2:
			return 1
		case v1 < v2:
			return -1
		default:
			return 0
		}
	default:

		panic(fmt.Sprintf("%v kind not handled", t1.Kind()))
	}

}

func (clsfy *Classify) Clear() {
	clsfy.values.Clear()
}
