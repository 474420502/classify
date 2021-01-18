package classify

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/474420502/focus/tree/vbtkey"
)

// kingTime 时间
var kingTime = reflect.TypeOf(time.Time{}).Kind()

var defaultCategoryHandler CategoryHandler = func(value interface{}) interface{} {
	return 0
}

// CategoryHandler 还回识别类别的keys. 返回的是子兄弟的类别.
type CategoryHandler func(value interface{}) interface{}

// Classify 分类
type Classify struct {
	// Name     string // 分类器的总名
	CategoryPath map[string]*CPath
	CategoryData *CData
}

// CData 分类数据
type CData struct {
	IsCollect bool
	Name      string
	Values    *vbtkey.Tree
}

// CPath 类别路径
type CPath struct {
	Name     string
	Handler  CategoryHandler
	IsValues bool
	Children map[string]*CPath
}

// New 创建新分类器
func New() *Classify {
	c := &Classify{}
	return c
}

// Categorys 返回所有数据的类别
func (c *Classify) Categorys() (result map[string]interface{}) {
	result = make(map[string]interface{})
	var categorys func(Category map[string]*CPath, root map[string]interface{})
	categorys = func(Category map[string]*CPath, root map[string]interface{}) {
		for _, p := range Category {
			if p.IsValues {
				return
			}
			// log.Println(p.Name)
			var child = make(map[string]interface{})
			root[p.Name] = child
			categorys(p.Children, child)
		}
	}
	categorys(c.CategoryPath, result)
	return
}

// Keys 返回所有数据的类别的key
func (c *Classify) Keys(paths ...interface{}) []interface{} {

	if c.CategoryData == nil {
		return nil
	}

	var outkeys []interface{}
	var get func(index int, Data *CData)
	get = func(index int, Data *CData) {

		if index >= len(paths) {
			Data.Values.Traversal(func(k, v interface{}) bool {
				outkeys = append(outkeys, k)
				return true
			})
			return
		}

		key := paths[index]
		// log.Println(Data.Values.String())
		if d, ok := Data.Values.Get(key); ok {
			next := d.(*CData)
			get(index+1, next)
		}
	}
	get(0, c.CategoryData)
	return outkeys
}

// Put 把数据压进分类器
func (c *Classify) Put(v interface{}) {
	if c.CategoryData == nil { //主要为了NewClassify 不添加其他属性. 使用nil指针
		c.CategoryData = &CData{}
	}

	put(c.CategoryPath, c.CategoryData, v)

}

// Get 获取路径的数据. 如果paths为nil 没输入则全部
func (c *Classify) Get(out interface{}, paths ...interface{}) {

	if c.CategoryData == nil {
		return
	}

	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		panic("out must ptr")
	}

	outv := reflect.ValueOf(out)
	result := outv.Elem()
	result = reflect.MakeSlice(result.Type(), 0, 0)
	// log.Println(result)

	var outcp []*CData
	var get func(index int, Data *CData)
	get = func(index int, Data *CData) {

		if index >= len(paths) {
			outcp = append(outcp, Data)
			return
		}

		key := paths[index]
		if d, ok := Data.Values.Get(key); ok {
			next := d.(*CData)
			get(index+1, next)
		}
	}
	get(0, c.CategoryData)

	var outdata func(Data *CData)
	outdata = func(Data *CData) {
		for _, v := range Data.Values.Values() {
			if cdata, ok := v.(*CData); ok {
				if cdata.IsCollect {
					for _, item := range cdata.Values.Values() {
						result = reflect.Append(result, reflect.ValueOf(item))
					}
				} else {
					outdata(cdata)
				}
			} else { // 可能拿到的是最终item值
				result = reflect.Append(result, reflect.ValueOf(v))
			}
		}
	}

	for _, data := range outcp {
		outdata(data)
	}

	// log.Println(result.Type())
	outv.Elem().Set(result)
}

// func get(index int, Category map[string]*ClassifyPath, Data *ClassifyData) {

// }

// DebugKeys 用于debug打印Keys
func (c *Classify) DebugKeys() {
	data, err := json.Marshal(c.Categorys())
	if err != nil {
		log.Panic(err)
	}
	log.Println(data)
}

func (c *Classify) debugPrint(limit int) {
	out("", c.CategoryPath, c.CategoryData, limit)
}

func out(parentName string, Category map[string]*CPath, Data *CData, limit int) {
	for _, p := range Category {
		// key := p.Handler(v)

		if p.IsValues {
			log.Println(fmt.Sprintf("data(%s) size=%d :", parentName, Data.Values.Size()))
			Data.Values.Traversal(func(k, v interface{}) bool {
				log.Println(v)
				limit--
				if limit <= 0 {
					log.Println("... .. .")
					return false
				}
				return true
			})
			return
		}

		log.Println(p.Name)
		Data.Values.Traversal(func(k, v interface{}) bool {
			next := v.(*CData)
			out(fmt.Sprint(k), p.Children, next, limit)
			return true
		})

		// 进入下一个类别
		// var next *ClassifyData

		// for _, inext := range Data.Values.Values() {

		// }
	}
}

func put(Category map[string]*CPath, Data *CData, v interface{}) {

	if Category != nil {
		for _, p := range Category {
			key := p.Handler(v)

			if Data.Values == nil {
				Data.Values = vbtkey.New(autoComapre)
			}

			if p.IsValues {
				if !Data.IsCollect {
					Data.IsCollect = true
				}
				Data.Values.Put(key, v)
				return
			}

			// 进入下一个类别
			var next *CData
			// var ok bool
			// log.Println(key, v.(*database.PayItem).CreateAt)
			if inext, ok := Data.Values.Get(key); ok {
				next = inext.(*CData)
			} else {
				next = &CData{Name: p.Name}
				Data.Values.Put(key, next)
			}

			put(p.Children, next, v)
		}
	}
}

// AddCategory 设置类别的处理句柄. 返回分类的key
func (c *Classify) AddCategory(name string, handler CategoryHandler) *Classify {

	if c.CategoryPath == nil {
		c.CategoryPath = make(map[string]*CPath)
	}

	path := &CPath{}
	path.Name = name
	c.CategoryPath[name] = path
	path.Handler = handler
	path.Children = make(map[string]*CPath)

	next := New()
	next.CategoryPath = path.Children

	// c.Children.Put(name, nc)
	return next
}

// Collect 设置类别的处理句柄. 区别于CollectCategory 没handler处理排序的返回key
func (c *Classify) Collect() {

	if c.CategoryPath == nil {
		c.CategoryPath = make(map[string]*CPath)
	}

	path := &CPath{}
	path.IsValues = true
	path.Name = "collect"
	c.CategoryPath[path.Name] = path
	path.Handler = defaultCategoryHandler
	path.Children = make(map[string]*CPath)

}

// CollectCategory 设置类别的处理句柄. 返回分类CategoryHandler的key. 用于排序
func (c *Classify) CollectCategory(handler CategoryHandler) {

	if c.CategoryPath == nil {
		c.CategoryPath = make(map[string]*CPath)
	}

	path := &CPath{}
	path.IsValues = true
	path.Name = "collect"
	c.CategoryPath[path.Name] = path
	path.Handler = handler
	path.Children = make(map[string]*CPath)

}

func autoComapre(k1, k2 interface{}) int {

	t1 := reflect.TypeOf(k1)
	t2 := reflect.TypeOf(k2)

	if t1.Kind() != t2.Kind() {
		log.Panicf("value1 %v, value2 %v is not same type. please check keys input", t1.Kind(), t2.Kind())
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
