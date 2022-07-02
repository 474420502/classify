package classify

import (
	"reflect"
	"time"
)

var timeKind = reflect.TypeOf(time.Time{}).Kind()

// 1 是字段( 2 是方法< 3 结尾收集@
type MethodType int

const (
	MT_UNKNOWN MethodType = 0 // 0 未知
	MT_FIELD   MethodType = 1 // 1 字段 KEY a-zA-Z
	MT_METHOD  MethodType = 2 // 2 方法 函数 <>
	MT_COLLECT MethodType = 3 // 3 收集操作 @
)

type CountHandler[T any] func(counted T, item T)
type CreateCountedHandler[T any] func(item T) T

// CategoryHandler 处理结构体字段的返回值
type CategoryHandler[T any] func(item T) any
