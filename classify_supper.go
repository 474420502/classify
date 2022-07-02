package classify

// type ClassifyEx[T any] struct {
// 	categorys   []CategoryHandler[T]
// 	bytesdict   *treelist.Tree
// 	uniqueCount uint64
// }

// // NewClassifyEx CreateCountedHandler CountHandler Add 都必须要使用地址传入
// func NewClassifyEx[T any](mode string, handlers ...CategoryHandler[T]) *ClassifyEx[T] {
// 	s := &ClassifyEx[T]{bytesdict: treelist.New()}
// 	s.Build(mode, handlers...)
// 	return s
// }

// // AddCategory 添加类别的处理方法
// func (stream *ClassifyEx[T]) AddCategory(handler CategoryHandler[T]) *ClassifyEx[T] {
// 	stream.categorys = append(stream.categorys, handler)
// 	return stream
// }

// // Add 添加到处理队列处理. 汇聚成counted. 通过 Seek RangeCounted获取结果
// func (stream *ClassifyEx[T]) Add(item T) {
// 	var skey = stream.getEncodeKey(item)

// 	if !stream.bytesdict.Put(skey, item) {
// 		log.Println("Warnning key is Conflict")
// 	}
// }

// // getEncodeKey 序列化 item的所有key
// func (stream *ClassifyEx[T]) getEncodeKey(item T) []byte {
// 	// var skey []byte
// 	var skey = bytes.NewBuffer(nil)

// 	for _, handler := range stream.categorys {
// 		skey.Write(handlerbytes(handler(item)))
// 	}
// 	err := binary.Write(skey, binary.BigEndian, stream.uniqueCount)
// 	if err != nil {
// 		panic(err)
// 	}
// 	stream.uniqueCount++
// 	return skey.Bytes()
// }

// // Seek 定位到 item 字节序列后的点. 然后从小到大遍历
// // [1 2 3] 参数为2 则 第一个item为2
// // [1 3] 参数为2 则 第一个item为3
// func (stream *ClassifyEx[T]) Seek(key T, iterfunc func(item interface{}) bool) {
// 	skey := stream.getEncodeKey(key)
// 	iter := stream.bytesdict.Iterator()
// 	iter.SeekGE(skey)

// 	for iter.Valid() {
// 		if !iterfunc(iter.Value()) {
// 			break
// 		}
// 		iter.Next()
// 	}
// }

// // Seek 定位到 item 字节序列后的点. 然后从大到小遍历
// // [1 2 3] 参数为2 则 第一个item为2
// // [1 3] 参数为2 则 第一个item为1.
// func (stream *ClassifyEx[T]) SeekReverse(item T, iterfunc func(item interface{}) bool) {
// 	skey := stream.getEncodeKey(item)
// 	iter := stream.bytesdict.Iterator()
// 	iter.SeekLT(skey)

// 	for iter.Valid() {
// 		if !iterfunc(iter.Value()) {
// 			break
// 		}
// 		iter.Prev()
// 	}
// }

// // RangeItem 从小到大遍历 counted 对象
// func (stream *ClassifyEx[T]) RangeItem(do func(item interface{}) bool) {
// 	stream.bytesdict.Traverse(func(s *treelist.Slice) bool {
// 		return do(s.Value)
// 	})
// }

// func (stream *ClassifyEx[T]) Build(mode string, handlers ...CategoryHandler[T]) {

// 	for _, token := range bytes.Split([]byte(mode), []byte{'.'}) {
// 		var i = 0

// 		// var cname []byte
// 		var cmethod []byte
// 		var mt MethodType = 0 // 1 是字段( 2 是方法< 3 结尾收集@ 4 是拼接 +

// 	CNAME:
// 		for ; i < len(token); i++ {
// 			c := token[i]
// 			switch c {
// 			case '@':
// 				panic("@ 不存在该操作符")
// 			case '+':
// 				break CNAME
// 			case '[':
// 				mt = MT_METHOD
// 				i++
// 				break CNAME
// 			case '<':
// 				mt = MT_FIELD
// 				i++
// 				break CNAME
// 			case ' ':
// 				continue
// 			default:
// 				// cname = append(cname, c)
// 			}
// 		}

// 		if mt == MT_UNKNOWN {
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

// 		switch mt {
// 		case MT_FIELD:
// 			// 添加处理字段的类别方法
// 			stream.AddCategory(func(value T) interface{} {
// 				v := reflect.ValueOf(value)
// 				if v.Type().Kind() == reflect.Ptr {
// 					v = v.Elem()
// 				}
// 				return v.FieldByName(string(cmethod)).Interface()
// 			})
// 		case MT_METHOD:
// 			// 通过自定义函数处理字段返回的方法
// 			fidx, err := strconv.Atoi(string(cmethod))
// 			if err != nil {
// 				panic(err)
// 			}
// 			stream.AddCategory(handlers[fidx])
// 		default:
// 			panic(fmt.Errorf("MethodType %d is error", mt))
// 		}

// 	}
// }
