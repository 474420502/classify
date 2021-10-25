package classify

import (
	"reflect"
)

type Iterator struct {
	cur  interface{}
	mark interface{}
}

// Value 遍历的数据
func (iter *Iterator) Value() interface{} {
	return iter.cur
}

// ValueMark  根据字段, 如果是上个类型不同值 则标记 isMakr = true. 效率稍低
func (iter *Iterator) ValueMark(fields ...string) (value interface{}, isMark int) {
	c := reflect.ValueOf(iter.cur)
	m := reflect.ValueOf(iter.mark)

	if c.Kind() == reflect.Ptr {
		c = c.Elem()
		m = m.Elem()
	}

	for i := len(fields) - 1; i > 0; i-- {
		field := fields[i]
		cv := c.FieldByName(field)
		mv := m.FieldByName(field)

		if !reflect.DeepEqual(cv.Interface(), mv.Interface()) {
			iter.SetMark()
			return iter.cur, i
		}
	}

	return iter.cur, -1
}

//  Mark 第一条Value()数据作为mark值
func (iter *Iterator) Mark() interface{} {
	return iter.mark
}

//  SetMark 把当前数据Value设置为Mark. 用于标记类型
func (iter *Iterator) SetMark() {
	iter.mark = iter.cur
}
