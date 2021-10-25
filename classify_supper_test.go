package classify

import (
	"log"
	"testing"
)

func TestCompare(t *testing.T) {
	clsfy := New()
	clsfy.Build(`region<RegionRFC>.country<RFC>.@Coin`)
	streamer := NewClassifyEx("region<RegionRFC>.country<RFC>")

	// items := loadGiftItemsJson()

	items := loadGiftItems()

	for _, item := range items {
		clsfy.Put(item)
		streamer.Put(item)
	}

	var classifylen int = 0
	var getitems []*GiftItem
	for _, region := range clsfy.Keys() {
		for _, country := range clsfy.Keys(region) {
			var gitems []*GiftItem
			clsfy.Get(&gitems, region, country)
			getitems = append(getitems, gitems...)
			classifylen++
		}
	}

	var classifylen2 int = 1
	var i = 0
	streamer.RangeItems(func(iter *Iterator) bool {
		item := iter.Value().(*GiftItem)

		if item.RFC != getitems[i].RFC {
			panic("")
		} else if item.RegionRFC != getitems[i].RegionRFC {
			panic("")
		}

		mark := iter.Mark().(*GiftItem)
		if item.RFC != mark.RFC {
			// log.Println(item.RegionRFC, item.SFC)
			classifylen2++
			iter.SetMark()
		} else if item.RegionRFC != mark.RegionRFC {
			// log.Println(item.RegionRFC, item.SFC)
			classifylen2++
			iter.SetMark()
		}
		i++
		return true
	})

	if classifylen != classifylen2 {
		t.Error("classifylen", classifylen, classifylen2)
	}

	streamer.SeekGE(&GiftItem{
		RegionRFC: "Region_Turkey",
	}, func(iter *Iterator) bool {

		item, isMark := iter.ValueMark("RegionRFC", "RFC") //  RegionRFC 0 RFC 1
		if isMark > 0 {                                    // 只要 RegionRFC 有变化就返回 isMark 对应 ValueMark参数的索引返回
			log.Println("---------------------")
			return true
		}
		log.Println(item.(*GiftItem).RegionRFC, item.(*GiftItem).RFC)

		return true
	})

	// streamer.SeekGE(&GiftItem{
	// 	RegionRFC: "default",
	// }, func(iter *Iterator) bool {

	// 	item := iter.Value().(*GiftItem)

	// 	mark := iter.Mark().(*GiftItem)
	// 	if mark.RegionRFC != item.RegionRFC {
	// 		iter.SetMark()
	// 		return false
	// 	}
	// 	log.Println(item)

	// 	return true
	// })

}
