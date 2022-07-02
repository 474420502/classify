package classify

import (
	"reflect"
	"time"
)

// kingTime 时间
var kingTime = reflect.TypeOf(time.Time{}).Kind()

var defaultCategoryHandler CategoryHandler[any] = func(value any) any {
	return 0
}

// 分类
// type Classify[T any] struct {
// 	categorys []*hCategory[T]
// 	Values    *treelist.Tree

// 	defaultCategoryHandler CategoryHandler[T]
// }

// // hCategory 类别
// type hCategory[T any] struct {
// 	Name      string
// 	Handler   CategoryHandler[T]
// 	IsCollect bool
// }

// func New[T any]() *Classify[T] {
// 	return &Classify[T]{defaultCategoryHandler: func(value T) any {
// 		return 0
// 	}}
// }

// func NewWithMode[T any](mode string, handlers ...CategoryHandler[T]) *Classify[T] {
// 	c := New[T]()
// 	c.Build(mode, handlers...)
// 	return c
// }

// func (clsfy *Classify[T]) Build(mode string, handlers ...CategoryHandler[T]) {

// 	for _, token := range bytes.Split([]byte(mode), []byte{'.'}) {
// 		var i = 0

// 		var cname []byte
// 		var cmethod []byte
// 		var methodType int = 0 // 1 是字段( 2 是方法< 3 结尾收集@

// 	CNAME:
// 		for ; i < len(token); i++ {
// 			c := token[i]
// 			switch c {
// 			case '@':
// 				methodType = 3
// 				i++
// 				break CNAME
// 			case '[':
// 				methodType = 2
// 				i++
// 				break CNAME
// 			case '<':
// 				methodType = 1
// 				i++
// 				break CNAME
// 			case ' ':
// 				continue
// 			default:
// 				cname = append(cname, c)
// 			}
// 		}

// 		if methodType == 0 {
// 			panic(fmt.Errorf("语法错误: %s", mode))
// 		}

// 	CMETHOD:
// 		for ; i < len(token); i++ {
// 			c := token[i]
// 			switch c {
// 			case ' ':
// 				continue
// 			case ']', '>':
// 				break CMETHOD
// 			default:
// 				cmethod = append(cmethod, c)
// 			}
// 		}

// 		switch methodType {
// 		case 1:
// 			log.Println(string(cname), string(cmethod))
// 			clsfy.AddCategory(string(cname), func(value T) any {
// 				v := reflect.ValueOf(value)
// 				if v.Type().Kind() == reflect.Ptr {
// 					v = v.Elem()
// 				}
// 				return v.FieldByName(string(cmethod)).Interface()
// 			})
// 		case 2:
// 			// log.Println(string(cname), string(cmethod))
// 			fidx, err := strconv.Atoi(string(cmethod))
// 			if err != nil {
// 				panic(err)
// 			}
// 			clsfy.AddCategory(string(cname), handlers[fidx])
// 		case 3:
// 			if len(cmethod) == 0 {
// 				clsfy.Collect()
// 			} else {
// 				clsfy.CollectCategory(func(value T) any {
// 					v := reflect.ValueOf(value)
// 					if v.Type().Kind() == reflect.Ptr {
// 						v = v.Elem()
// 					}
// 					return v.FieldByName(string(cmethod)).Interface()
// 				})
// 			}

// 		default:
// 			panic("?")
// 		}

// 	}
// }

// func (clsfy *Classify[T]) AddCategory(name string, handler CategoryHandler[T]) *Classify[T] {
// 	clsfy.categorys = append(clsfy.categorys, &hCategory[T]{
// 		Name:      name,
// 		Handler:   handler,
// 		IsCollect: false,
// 	})
// 	return clsfy
// }

// func (clsfy *Classify[T]) Categorys() string {
// 	if len(clsfy.categorys) == 0 {
// 		return ""
// 	}
// 	var content []byte
// 	for _, cate := range clsfy.categorys {
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

// func (clsfy *Classify[T]) Collect() {
// 	clsfy.categorys = append(clsfy.categorys, &hCategory[T]{
// 		Name:      "@",
// 		Handler:   clsfy.defaultCategoryHandler,
// 		IsCollect: true,
// 	})
// }

// func (clsfy *Classify[T]) CollectCategory(handler CategoryHandler[T]) {
// 	clsfy.categorys = append(clsfy.categorys, &hCategory[T]{
// 		Name:      "@",
// 		Handler:   handler,
// 		IsCollect: true,
// 	})
// }

// func (clsfy *Classify[T]) Keys(paths ...interface{}) (result []interface{}) {
// 	var values *treelist.Tree = clsfy.Values
// 	// var category *Category
// 	if len(paths) >= len(clsfy.categorys)-1 {
// 		panic(fmt.Sprintf("categorys len is %d only: %#v", len(clsfy.categorys)-1, clsfy.Categorys()))
// 	}
// 	for _, p := range paths {
// 		// category = clsfy.Categorys[i]
// 		if child, ok := values.Get(p); ok {
// 			values = child.(*treelist.Tree)
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

// func (clsfy *Classify[T]) Put(v T) {
// 	if clsfy.Values == nil {
// 		clsfy.Values = vbtkey.New(autoComapre)
// 	}
// 	put[T](clsfy.categorys, 0, clsfy.Values, v)
// }

// func (clsfy *Classify[T]) PutSlice(items []T) {

// 	if clsfy.Values == nil {
// 		clsfy.Values = vbtkey.New(autoComapre)
// 	}
// 	vitems := reflect.ValueOf(items)
// 	if vitems.Type().Kind() != reflect.Slice {
// 		panic(" input must slice ")
// 	}
// 	for i := 0; i < len(items); i++ {
// 		put[T](clsfy.categorys, 0, clsfy.Values, items[i])
// 	}

// }

// func put[T any](categorys []*hCategory[T], cidx int, Values *treelist.Tree, v T) {
// 	cate := categorys[cidx]
// 	if cate.IsCollect {
// 		Values.Put(cate.Handler(v), v)
// 		return
// 	} else {
// 		// 判断Values是否存在
// 		var NextValues *treelist.Tree
// 		key := cate.Handler(v)
// 		if vs, ok := Values.Get(key); ok {
// 			NextValues = vs.(*treelist.Tree)
// 		} else {
// 			NextValues = vbtkey.New(autoComapre)
// 			Values.Put(key, NextValues)
// 		}
// 		put(categorys, cidx+1, NextValues, v)
// 	}
// }

// func (clsfy *Classify[T]) Get(out interface{}, vPaths ...interface{}) {

// 	var values *treelist.Tree = clsfy.Values
// 	if len(vPaths) >= len(clsfy.categorys) {
// 		panic(fmt.Sprintf("values keys deepth is %d only: %#v", len(clsfy.categorys), clsfy.Categorys()))
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
// 			values = child.(*treelist.Tree)
// 		}
// 	}

// 	var getValues func(cidx int, values *treelist.Tree)
// 	getValues = func(cidx int, values *treelist.Tree) {
// 		category := clsfy.categorys[cidx]
// 		if category.IsCollect {
// 			values.Traversal(func(k, v interface{}) bool {
// 				result = reflect.Append(result, reflect.ValueOf(v))
// 				return true
// 			})
// 		} else {
// 			values.Traversal(func(k, v interface{}) bool {
// 				getValues(cidx+1, v.(*treelist.Tree))
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
