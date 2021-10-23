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

type CountHandler func(counted interface{}, item interface{})
type CreateCountedHandler func(item interface{}) interface{}
