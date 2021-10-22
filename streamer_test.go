package classify

import (
	"testing"
	"time"

	"github.com/474420502/random"
	"github.com/Pallinder/go-randomdata"
)

type TestStruct struct {
	Name  string
	Label string
	Type  int32
	Value int32
}

func TestMode(t *testing.T) {

	streamer := NewStreamer("<Name>.<Type>")
	streamer.SetCreateCountedHandler(func(passitem interface{}) interface{} {
		item := passitem.(*TestStruct)
		return &TestStruct{
			Name:  item.Name,
			Label: item.Label,
			Value: item.Value,
			Type:  item.Type,
		}
	})
	streamer.SetCountHandler(func(counted, item interface{}) {
		c := counted.(*TestStruct)
		c.Value += item.(*TestStruct).Value
	})

	end := 10
	for i := 0; i < 10000; i++ {

		streamer.Add(&TestStruct{
			Name:  "day",
			Label: "xi",
			Type:  int32(randomdata.Number(0, end)),
			Value: 1,
		})

		streamer.Add(&TestStruct{
			Name:  "month",
			Label: "ha",
			Type:  int32(randomdata.Number(0, end)),
			Value: 1,
		})

		streamer.Add(&TestStruct{
			Name:  "day",
			Label: "xi",
			Type:  int32(randomdata.Number(0, end)),
			Value: 2,
		})
	}
}

func TestSortRange(t *testing.T) {

	rand := random.New()
	streamer := NewStreamer("<Type>.<Name>.<Label>")

	// 创建 生成统计的类型.
	streamer.SetCreateCountedHandler(func(passitem interface{}) interface{} {
		item := passitem.(*TestStruct)
		return &TestStruct{
			Name:  item.Name,
			Label: item.Label,
			Value: item.Value,
			Type:  item.Type,
		}
	})

	//  统计的每一个Add item. 累加之类的. 作为流计算
	streamer.SetCountHandler(func(counted, item interface{}) {
		c := counted.(*TestStruct)
		c.Value += item.(*TestStruct).Value
	})

	random.Use(random.DataNameChina)
	random.Use(random.DataIdidomChina)

	for i := 0; i < 200; i++ {
		item := &TestStruct{
			Name:  rand.Extend().FullName(),
			Label: rand.Extend().Ididom().Name,
			Value: rand.Int31(),
			Type:  rand.Int31n(100),
		}
		streamer.Add(item)
	}

	// /home/eson/workspace/classfiy/streamer_test.go:100: &{麻致爱 服服帖帖 93 1270299833}
	// 	/home/eson/workspace/classfiy/streamer_test.go:97: &{黄镇旨 蛮烟瘴雨 33 1350982193}
	// /home/eson/workspace/classfiy/streamer_test.go:97: &{齐因泽 街头巷口 44 466969591}

	streamer.Seek(TestStruct{Type: 70}, func(item interface{}) bool {
		i := item.(*TestStruct)
		if i.Type < 70 {
			t.Error("Seek error")
			return false
		}
		return true
	})

	streamer.SeekReverse(TestStruct{Type: 70}, func(item interface{}) bool {
		i := item.(*TestStruct)
		if i.Type > 70 {
			t.Error("Seek error")
			return false
		}
		return true
	})

}

type TSTime struct {
	CreateAt time.Time
	Name     string
	Label    string
	Type     int32
	Value    int32
}

func TestSortRangeMethod(t *testing.T) {

	rand := random.New()
	streamer := NewStreamer("[0].<Type>.<Name>.<Label>", func(item interface{}) interface{} {
		return item.(*TSTime).CreateAt.Truncate(time.Hour)
	})

	// 创建 生成统计的类型.
	streamer.SetCreateCountedHandler(func(passitem interface{}) interface{} {
		item := passitem.(*TSTime)
		return &TSTime{
			CreateAt: item.CreateAt.Truncate(time.Hour),
			Name:     item.Name,
			Label:    item.Label,
			Value:    item.Value,
			Type:     item.Type,
		}
	})

	//  统计的每一个Add item. 累加之类的. 作为流计算
	streamer.SetCountHandler(func(counted, item interface{}) {
		c := counted.(*TSTime)
		c.Value += item.(*TSTime).Value
	})

	random.Use(random.DataNameChina)
	random.Use(random.DataIdidomChina)

	now := time.Unix(1634895792, 0)
	for i := 0; i < 200; i++ {
		item := &TSTime{
			CreateAt: now,
			Name:     rand.Extend().FullName(),
			Label:    rand.Extend().Ididom().Name,
			Value:    rand.Int31(),
			Type:     rand.Int31n(100),
		}
		streamer.Add(item)
		now = now.Add(time.Minute * 15)
	}

	now = time.Unix(1634895792, 0).Add(time.Hour * 5).Truncate(time.Hour)
	streamer.Seek(&TSTime{CreateAt: now}, func(item interface{}) bool {
		i := item.(*TSTime)
		if i.CreateAt.Before(now) {
			panic("time seek error")
		}
		return true
	})

	streamer.SeekReverse(&TSTime{Type: 70}, func(item interface{}) bool {
		i := item.(*TSTime)
		if i.CreateAt.Before(now) {
			panic("time seek error")
		}
		return true
	})

}
