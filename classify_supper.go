package classify

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/474420502/structure/search/treelist"
)

type ClassifyEx struct {
	categorys   []CategoryHandler
	bytesdict   *treelist.Tree
	uniqueCount uint64
}

// NewClassifyEx CreateCountedHandler CountHandler Add 都必须要使用地址传入
func NewClassifyEx(mode string, handlers ...CategoryHandler) *ClassifyEx {
	s := &ClassifyEx{bytesdict: treelist.New()}
	s.Build(mode, handlers...)
	return s
}

// AddCategory 添加类别的处理方法
func (stream *ClassifyEx) AddCategory(handler CategoryHandler) *ClassifyEx {
	stream.categorys = append(stream.categorys, handler)
	return stream
}

// Add 添加到处理队列处理.  通过 Seek RangeItem获取结果
func (stream *ClassifyEx) Add(item interface{}) {
	var skey = stream.getEncodeKey(item)

	if !stream.bytesdict.Put(skey, item) {
		log.Println("Warnning key is Conflict")
	}
}

// getEncodeKey 序列化 item的所有key
func (stream *ClassifyEx) getEncodeKey(item interface{}) []byte {
	// var skey []byte
	var skey = bytes.NewBuffer(nil)

	for _, handler := range stream.categorys {
		skey.Write(handlerbytes(handler(item)))
	}
	err := binary.Write(skey, binary.BigEndian, stream.uniqueCount)
	if err != nil {
		panic(err)
	}
	stream.uniqueCount++
	return skey.Bytes()
}

// SeekGE 定位到 item 字节序列后的点. 然后从小到大遍历
// [1 2 3] 参数为2 则 第一个item为2
// [1 3] 参数为2 则 第一个item为3
func (stream *ClassifyEx) SeekGE(key interface{}, iterfunc func(item interface{}) bool) {
	skey := stream.getEncodeKey(key)
	iter := stream.bytesdict.Iterator()
	iter.SeekGE(skey)
	for iter.Valid() {
		if !iterfunc(iter.Value()) {
			break
		}
		iter.Next()
	}
}

// SeekGEReverse 定位到 item 字节序列后的点. 然后从大到小遍历
// [1 2 3] 参数为2 则 第一个item为2
// [1 3] 参数为2 则 第一个item为1.
func (stream *ClassifyEx) SeekGEReverse(item interface{}, iterfunc func(item interface{}) bool) {
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

// RangeItems 从小到大遍历 item 对象
func (stream *ClassifyEx) RangeItems(do func(item interface{}) bool) {
	stream.bytesdict.Traverse(func(s *treelist.Slice) bool {
		return do(s.Value)
	})
}

func (stream *ClassifyEx) Build(mode string, handlers ...CategoryHandler) {

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
