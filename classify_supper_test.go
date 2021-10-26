package classify

// func TestCompare(t *testing.T) {
// 	// clsfy := New()
// 	// clsfy.Build(`region<RegionRFC>.country<RFC>.@Coin`)
// 	// streamer := NewClassifyEx("region<RegionRFC>.country<RFC>")

// 	// // items := loadGiftItemsJson()

// 	// items := loadGiftItems()

// 	// for _, item := range items {
// 	// 	clsfy.Put(item)
// 	// 	streamer.Put(item)
// 	// }

// 	// var classifylen int = 0
// 	// var getitems []*GiftItem
// 	// for _, region := range clsfy.Keys() {
// 	// 	for _, country := range clsfy.Keys(region) {
// 	// 		var gitems []*GiftItem
// 	// 		clsfy.Get(&gitems, region, country)
// 	// 		getitems = append(getitems, gitems...)
// 	// 		classifylen++
// 	// 	}
// 	// }

// 	// var classifylen2 int = 1
// 	// var i = 0
// 	// streamer.RangeItems(func(iter *Iterator) bool {
// 	// 	item := iter.Value().(*GiftItem)

// 	// 	if item.RFC != getitems[i].RFC {
// 	// 		panic("")
// 	// 	} else if item.RegionRFC != getitems[i].RegionRFC {
// 	// 		panic("")
// 	// 	}

// 	// 	mark := iter.Mark().(*GiftItem)
// 	// 	if item.RFC != mark.RFC {
// 	// 		// log.Println(item.RegionRFC, item.SFC)
// 	// 		classifylen2++
// 	// 		iter.SetMark()
// 	// 	} else if item.RegionRFC != mark.RegionRFC {
// 	// 		// log.Println(item.RegionRFC, item.SFC)
// 	// 		classifylen2++
// 	// 		iter.SetMark()
// 	// 	}
// 	// 	i++
// 	// 	return true
// 	// })

// 	// if classifylen != classifylen2 {
// 	// 	t.Error("classifylen", classifylen, classifylen2)
// 	// }

// 	// var (
// 	// 	compareitems1 []interface{}
// 	// 	compareitems2 []interface{}
// 	// )

// 	// streamer.SeekGE(&GiftItem{
// 	// 	RegionRFC: "Region_Turkey",
// 	// }, func(iter *Iterator) bool {

// 	// 	item := iter.Value().(*GiftItem)

// 	// 	mark := iter.Mark().(*GiftItem)

// 	// 	if mark.RFC != item.RFC {
// 	// 		iter.SetMark()
// 	// 		compareitems2 = append(compareitems2, item)
// 	// 	} else if mark.RegionRFC != item.RegionRFC {
// 	// 		iter.SetMark()
// 	// 		compareitems2 = append(compareitems2, item)
// 	// 		return true
// 	// 	}
// 	// 	return true
// 	// })

// 	// if !reflect.DeepEqual(compareitems1, compareitems2) {
// 	// 	t.Error("compareitems1 != compareitems2")
// 	// }

// 	// compareitems1 = nil
// 	// compareitems2 = nil

// 	// streamer.SeekGEReverse(&GiftItem{
// 	// 	RegionRFC: "default",
// 	// }, func(iter *Iterator) bool {

// 	// 	item := iter.Value().(*GiftItem)

// 	// 	mark := iter.Mark().(*GiftItem)

// 	// 	if mark.RFC != item.RFC {
// 	// 		iter.SetMark()
// 	// 		compareitems2 = append(compareitems2, item)
// 	// 	} else if mark.RegionRFC != item.RegionRFC {
// 	// 		iter.SetMark()
// 	// 		compareitems2 = append(compareitems2, item)
// 	// 		return true
// 	// 	}
// 	// 	return true
// 	// })

// 	// if !reflect.DeepEqual(compareitems1, compareitems2) {
// 	// 	t.Error("compareitems1 != compareitems2")
// 	// }

// }

// type TestClassify struct {
// 	CreateAt time.Time
// 	Label    string
// 	Name     string
// 	Value    int32
// }

// func TestClassifyExForce(t *testing.T) {
// 	rand := random.New(1635218347457601795)

// 	random.Use(random.DataCityChina)
// 	random.Use(random.DataNameChina)

// 	streamer := NewClassifyEx("[0].<Label>.<Name>.<Value>", func(item interface{}) interface{} {
// 		return item.(*TestClassify).CreateAt.Truncate(time.Minute * 5)
// 	})

// 	clsfy := New()
// 	clsfy.Build(`CreateAt[0].<Label>.<Name>.@Value`, func(item interface{}) interface{} {
// 		return item.(*TestClassify).CreateAt.Truncate(time.Minute * 5)
// 	})

// 	now := time.Unix(1635215668, 0)
// 	for i := 0; i < 1000; i++ {

// 		for n := 0; n < rand.Intn(64)+1; n++ {
// 			item := &TestClassify{
// 				CreateAt: now,
// 				Label:    rand.Extend().City(),
// 				Name:     rand.Extend().FullName(),
// 				Value:    rand.Int31n(9999999),
// 			}
// 			streamer.Put(item)
// 			clsfy.Put(item)
// 			now = now.Add(time.Minute * time.Duration(rand.Intn(15)+1))
// 		}

// 		var (
// 			compareitems1 []interface{}
// 			compareitems2 []interface{}
// 		)

// 		for _, createAt := range clsfy.Keys() {
// 			for _, label := range clsfy.Keys(createAt) {
// 				var gitems []*TestClassify
// 				clsfy.Get(&gitems, createAt, label)
// 				compareitems1 = append(compareitems1, gitems[0])
// 			}
// 		}

// 		var count = 0
// 		streamer.RangeItems(func(iter *Iterator) bool {

// 			item := iter.Value().(*TestClassify)
// 			mark := iter.Mark().(*TestClassify)

// 			if count == 0 {
// 				iter.SetMark()
// 				compareitems2 = append(compareitems2, item)
// 				count++
// 				return true
// 			}

// 			if mark.Label != item.Label {
// 				iter.SetMark()
// 				compareitems2 = append(compareitems2, item)
// 			} else if !mark.CreateAt.Truncate(time.Minute * 5).Equal(item.CreateAt.Truncate(time.Minute * 5)) {
// 				iter.SetMark()
// 				compareitems2 = append(compareitems2, item)
// 			}
// 			return true
// 		})

// 		if !reflect.DeepEqual(compareitems1, compareitems2) {
// 			t.Error("compareitems1 != compareitems2")
// 		}

// 		clsfy.Clear()
// 		streamer.Clear()
// 	}

// }
