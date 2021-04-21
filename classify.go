package classify

import (
	"log"
)

type CollectHandler func(value interface{}) interface{}

type CollectType int

const (
	CollectField  CollectType = 1 // 字段
	CollectMethod CollectType = 2 // 方法
	CollectArray  CollectType = 3 // 并行数组

	CollectEnd     CollectType = 4 // 结束
	CollectDescEnd CollectType = 5 // 排序 降 结束
	CollectAscEnd  CollectType = 6 // 排序 升 结束

)

type Classify struct {
	Handlers   map[string]CollectHandler
	root       *category     // 类别
	collection []*collection // 数据集合
}

// cursor 游标
type cursor struct {
	Data    []byte
	Index   int
	Size    int
	Unknown int
}

type category struct {
	Name string

	CType   CollectType    // 提取类型
	CMethod string         // 提取方法
	Collect CollectHandler // 收集分类的方法
	// IsEnd   bool           // 是否结尾
	Next []*category
}

func newCategory() *category {
	return &category{}
}

func (csf *Classify) Build(path string) {
	csf.root = &category{}
	extract(csf.root, headerCompletion([]byte(path)))
}

func (csf *Classify) Put(item interface{}) {
	for i, cur := range csf.root.Next {
		csf.put(cur, csf.collection[i], item)
	}
}

func (csf *Classify) put(parent *category, collection *collection, item interface{}) {
	switch parent.CType {
	case CollectField: // 字段
		// itype := reflect.TypeOf(item)
		// ivalue := reflect.ValueOf(item)
		// v := ivalue.FieldByName(parent.CMethod).Interface()
		// csf.collection[v] =
	case CollectMethod: // 方法
	}
}

// extract 提取初始入口
func extract(parent *category, cur *cursor) {

	for ; cur.Index < cur.Size; cur.Index++ {
		c := cur.Data[cur.Index]
		switch c {
		case ' ', '\n':
			continue
		case '.':
			// 进入 Second paragraph

			cur.Index++
			c = cur.Data[cur.Index]
			switch c {
			case ' ':
				log.Panic("space behind '.'")
			case '[':
				for _, aCur := range extractCollectArray(cur) {
					log.Println(string(aCur.Data[aCur.Index:aCur.Size]))
					extract(parent, aCur)
				}
				cur.Index++
			case '@':
				// 提取排序字段
				cur.Index++
				classify := &category{}
				// classify.IsEnd = true

				c = cur.Data[cur.Index]
				if c == '!' {
					classify.CType = CollectAscEnd
					cur.Index++
				} else {
					classify.CType = CollectDescEnd
				}

				var label []byte
				for (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
					label = append(label, c)
					cur.Index++
					c = cur.Data[cur.Index]
				}
				classify.Name = string(label)
				log.Println("@" + classify.Name)
				parent.Next = append(parent.Next, classify)
				return
			default:
				// 提取字段名
				extractClassify(parent, cur)
			}

		default:
			// 提取字段名
			extractClassify(parent, cur)
		}
	}
}

func extractClassify(parent *category, cur *cursor) {
	c := cur.Data[cur.Index]

	if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
		var label []byte = []byte{c}
		cur.Index++
		c = cur.Data[cur.Index]
		for (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			label = append(label, c)
			cur.Index++
			c = cur.Data[cur.Index]
		}

		var classify *category = newCategory()
		classify.Name = string(label)
		log.Println(classify.Name)
		parent.Next = append(parent.Next, classify)

		// 进入 Second paragraph
		switch c {
		case '<':
			// 字段
			classify.CType = CollectField
			classify.CMethod = extractCollectField(cur)
			extract(classify, cur)
		case '(':
			classify.CType = CollectMethod
			classify.CMethod = extractCollectMethod(cur)
			extract(classify, cur)
			// 方法
		case '[':
			// 并行数组
			for _, aCur := range extractCollectArray(cur) {
				extract(classify, aCur)
			}
			cur.Index++
		}
	} else {
		log.Println(string(c))
	}
}

func extractCollectField(cur *cursor) string {
	var fieldname []byte
	cur.Index++

	for ; cur.Index < cur.Size; cur.Index++ {
		c := cur.Data[cur.Index]
		switch c {
		case ' ':
			continue
		case '>':
			cur.Index++
			return string(fieldname)
		default:
			fieldname = append(fieldname, c)
		}
	}
	log.Panic("can't find '>' end char")
	return ""
}

func extractCollectMethod(cur *cursor) string {
	var methodname []byte
	cur.Index++

	for ; cur.Index < cur.Size; cur.Index++ {
		c := cur.Data[cur.Index]
		switch c {
		case ' ':
			continue
		case ')':
			cur.Index++
			return string(methodname)
		default:
			methodname = append(methodname, c)
		}
	}
	log.Panic("can't find ')' end char")
	return ""
}

func extractCollectArray(cur *cursor) (Objects []*cursor) {

	cur.Index++
	start := cur.Index

	for ; cur.Index < cur.Size; cur.Index++ {
		c := cur.Data[cur.Index]
		switch c {
		case ' ':
			continue
		case ',':
			aCur := &cursor{
				Data:  cur.Data,
				Index: start,
				Size:  cur.Index,
			}
			Objects = append(Objects, aCur)
			start = aCur.Size + 1
		case ']':

			aCur := &cursor{
				Data:  cur.Data,
				Index: start,
				Size:  cur.Index,
			}
			Objects = append(Objects, aCur)
			cur.Index++
			return
		default:
			// methodname = append(methodname, c)
		}
	}

	log.Panic("can't find ']' end char")
	return
}

func headerCompletion(data []byte) *cursor {

	for i := 0; i < len(data); i++ {
		if data[i] == ' ' {
			continue
		} else {
			cur := &cursor{}
			cur.Data = append(cur.Data, data...)
			cur.Size = len(cur.Data)
			cur.Data = append(cur.Data, ' ')
			return cur
		}
	}

	log.Panic("data is nil")
	return nil
}
