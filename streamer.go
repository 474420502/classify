package classify

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/474420502/structure/search/treelist"
)

// Streamer 流计算. 用于分段时间聚合, 分类聚合...等
type Streamer struct {
	categorys     []CategoryHandler
	bytesdict     *treelist.Tree
	count         CountHandler
	createHandler CreateCountedHandler
}

func handlerbytes(item interface{}) []byte {
	var buf bytes.Buffer

	fi := reflect.ValueOf(item)
	if fi.Kind() == reflect.Ptr {
		fi = fi.Elem()
	}

	switch fi.Kind() {
	case reflect.String:
		err := binary.Write(&buf, binary.BigEndian, []byte(fi.Interface().(string)))
		if err != nil {
			panic(err)
		}
	case timeKind:
		data, err := fi.Interface().(time.Time).MarshalBinary()
		if err != nil {
			panic(err)
		}
		err = binary.Write(&buf, binary.BigEndian, data)
		if err != nil {
			panic(err)
		}
	default:
		err := binary.Write(&buf, binary.BigEndian, fi.Interface())
		if err != nil {
			panic(err)
		}
	}

	return buf.Bytes()
}

// NewStreamer CreateCountedHandler CountHandler Add 都必须要使用地址传入
func NewStreamer(mode string, handlers ...CategoryHandler) *Streamer {
	s := &Streamer{bytesdict: treelist.New()}
	s.Build(mode, handlers...)
	return s
}

// SetCreateCountedHandler 设置被计算生成的对象 通过所有Add item汇聚成handler 返回的结果
func (stream *Streamer) SetCreateCountedHandler(createHandler CreateCountedHandler) *Streamer {
	stream.createHandler = createHandler
	return stream
}

// SetCountHandler  设置计算过程
func (stream *Streamer) SetCountHandler(countHandler CountHandler) *Streamer {
	stream.count = countHandler
	return stream
}

// AddCategory 添加类别的处理方法
func (stream *Streamer) AddCategory(handler CategoryHandler) *Streamer {
	stream.categorys = append(stream.categorys, handler)
	return stream
}

// Add 添加到处理队列处理. 汇聚成counted. 通过 Seek RangeCounted获取结果
func (stream *Streamer) Add(item interface{}) {
	skey := stream.getEncodeKey(item)

	// var skey = bkey
	var counted interface{}
	var ok bool

	if counted, ok = stream.bytesdict.Get(skey); !ok {
		if stream.createHandler != nil {
			counted = stream.createHandler(item)
		} else {
			counted = item
		}

		// 必须地址传入 所以counted必须地址
		// stream.bytesdict[skey] = counted
		stream.bytesdict.Put(skey, counted)
	} else {
		// 必须地址传入
		stream.count(counted, item)
	}
}

// AddSlice items添加到处理队列处理. 汇聚成counted. 通过 Seek RangeCounted获取结果
func (stream *Streamer) AddSlice(items interface{}) {
	vitems := reflect.ValueOf(items)
	if vitems.Type().Kind() != reflect.Slice {
		panic(" input must slice ")
	}
	for i := 0; i < vitems.Len(); i++ {
		stream.Add(vitems.Index(i).Interface())
	}
}

// getEncodeKey 序列化 item的所有key
func (stream *Streamer) getEncodeKey(item interface{}) []byte {
	var skey []byte
	for _, handler := range stream.categorys {
		skey = append(skey, handlerbytes(handler(item))...)
	}
	return skey
}

// SeekGE 定位到 item 字节序列后的点. 然后从小到大遍历
// [1 2 3] 参数为2 则 第一个item为2
// [1 3] 参数为2 则 第一个item为3
func (stream *Streamer) SeekGE(item interface{}, iterfunc func(counted interface{}) bool) {
	skey := stream.getEncodeKey(item)
	iter := stream.bytesdict.Iterator()
	iter.SeekGE(skey)

	for iter.Valid() {
		if !iterfunc(iter.Value()) {
			break
		}
		iter.Next()
	}
}

// Seek 定位到 item 字节序列后的点. 然后从大到小遍历
// [1 2 3] 参数为2 则 第一个item为2
// [1 3] 参数为2 则 第一个item为1.
func (stream *Streamer) SeekGEReverse(item interface{}, iterfunc func(counted interface{}) bool) {
	skey := stream.getEncodeKey(item)
	iter := stream.bytesdict.Iterator()
	iter.SeekLE(skey)

	for iter.Valid() {
		if !iterfunc(iter.Value()) {
			break
		}
		iter.Prev()
	}
}

// RangeCounted 从小到大遍历 counted 对象
func (stream *Streamer) RangeCounted(do func(counted interface{}) bool) {
	stream.bytesdict.Traverse(func(s *treelist.Slice) bool {
		return do(s.Value)
	})
}

func (stream *Streamer) Build(mode string, handlers ...CategoryHandler) {

	for _, token := range bytes.Split([]byte(mode), []byte{'.'}) {
		var i = 0

		// var cname []byte
		var cmethod []byte
		var mt MethodType = 0 // 1 是字段( 2 是方法< 3 结尾收集@ 4 是拼接 +

	CNAME:
		for ; i < len(token); i++ {
			c := token[i]
			switch c {
			case '@':
				panic("@ 不存在该操作符")
			case '+':
				break CNAME
			case '[':
				mt = MT_METHOD
				i++
				break CNAME
			case '<':
				mt = MT_FIELD
				i++
				break CNAME
			case ' ':
				continue
			default:
				// cname = append(cname, c)
			}
		}

		if mt == MT_UNKNOWN {
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

		switch mt {
		case MT_FIELD:
			// 添加处理字段的类别方法
			stream.AddCategory(func(value interface{}) interface{} {
				v := reflect.ValueOf(value)
				if v.Type().Kind() == reflect.Ptr {
					v = v.Elem()
				}
				return v.FieldByName(string(cmethod)).Interface()
			})
		case MT_METHOD:
			// 通过自定义函数处理字段返回的方法
			fidx, err := strconv.Atoi(string(cmethod))
			if err != nil {
				panic(err)
			}
			stream.AddCategory(handlers[fidx])
		default:
			panic(fmt.Errorf("MethodType %d is error", mt))
		}

	}
}

func (stream *Streamer) Clear() {
	stream.bytesdict.Clear()
}
