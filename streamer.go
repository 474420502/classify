package classify

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"reflect"
	"strconv"
)

// Streamer 流计算
type Streamer struct {
	categorys []*sCategory
	bytesdict map[string]interface{}

	count         CountHandler
	createHandler CreateCountedHandler
}

type sCategory struct {
	Handler CategoryHandler
}

type CountHandler func(counted interface{}, item interface{})
type CreateCountedHandler func(item interface{}) interface{}

func handlerbytes(item interface{}) []byte {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(item)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func NewStreamer(mode string, handlers ...CategoryHandler) *Streamer {
	s := &Streamer{bytesdict: make(map[string]interface{})}
	s.Build(mode, handlers...)
	return s
}

func (stream *Streamer) SetCreateCountedHandler(createHandler CreateCountedHandler) *Streamer {
	test := reflect.ValueOf(createHandler(nil))
	if test.Kind() == reflect.Ptr {
		test = test.Elem()
	}
	if test.Kind() != reflect.Struct {
		panic("countInitHandler must be struct")
	}
	stream.createHandler = createHandler
	return stream
}

func (stream *Streamer) SetCountHandler(countHandler CountHandler) *Streamer {
	stream.count = countHandler
	return stream
}

func (stream *Streamer) AddCategory(handler CategoryHandler) *Streamer {
	stream.categorys = append(stream.categorys, &sCategory{
		Handler: handler,
	})
	return stream
}

func (stream *Streamer) Add(item interface{}) {

	var bkey []byte
	for _, cg := range stream.categorys {
		bkey = append(bkey, handlerbytes(cg.Handler(item))...)
	}

	var counted interface{}
	var ok bool

	if counted, ok = stream.bytesdict[string(bkey)]; !ok {
		counted = stream.createHandler(item)
	}

	vcounted := reflect.ValueOf(counted)
	if vcounted.Kind() != reflect.Ptr {
		counted = vcounted.Addr().Interface()
	}

	stream.bytesdict[string(bkey)] = counted
}

func (stream *Streamer) Build(mode string, handlers ...CategoryHandler) {

	for _, token := range bytes.Split([]byte(mode), []byte{'.'}) {
		var i = 0

		// var cname []byte
		var cmethod []byte
		var methodType int = 0 // 1 是字段( 2 是方法< 3 结尾收集@

	CNAME:
		for ; i < len(token); i++ {
			c := token[i]
			switch c {
			case '@':
				panic("@ 不存在该操作符")
				// methodType = 3
				// i++
				// break CNAME
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
				// cname = append(cname, c)
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
			log.Println(string(cmethod))
			stream.AddCategory(func(value interface{}) interface{} {
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
			stream.AddCategory(handlers[fidx])

		default:
			panic("?")
		}

	}
}
