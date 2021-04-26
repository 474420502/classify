package classify

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"sort"
	"testing"
	"time"

	database "git.nonolive.co/eson.hsm/databoard-database-myrocks"
)

func init() {
	log.SetFlags(log.Llongfile)
}

var db = database.DB

// GiftItem 单条充值数据
type GiftItem struct {
	GiftID        int64     `json:"gift_id"`         // 礼物ID
	GiftItemID    string    `json:"gift_item_id"`    // gitf item 的id
	Name          string    `json:"gift_name"`       // gitf item 的 name
	Sender        string    `json:"sender"`          // 送礼者名称
	Receive       string    `json:"receive"`         // 收礼者名称
	SenderUserID  int64     `json:"sender_user_id"`  // 送礼物者 用户 id
	ReceiveUserID int64     `json:"receive_user_id"` // 接收礼物 用户 id\
	ItemType      int64     `json:"item_type"`       // 礼物类型
	Picture       string    `json:"pic"`             // 图片路径
	Coin          int64     `json:"coin"`            // 金币的价值
	RealCoin      float64   `json:"real_coin"`       // 主播收到的金币
	Quantity      int64     `json:"quantity"`        // 数量 number
	SFC           string    `json:"sfc"`             // 送礼者来自地区
	RegionSFC     string    `json:"region_sfc"`      //送礼者来自区域
	RFC           string    `json:"rfc"`             // 收礼者来自地区
	RegionRFC     string    `json:"region_rfc"`      //收礼来自区域
	CreateAt      time.Time `json:"create_at"`       // 送礼时间
}

// QueryGiftItems 查询礼物的数据
func QueryGiftItems(start, end time.Time) (result []*GiftItem) {

	db.Do(func(db *sql.DB) {
		ssql := fmt.Sprintf("select gift_id,gift_item_id,gift_name,sender,receive,sender_user_id,receive_user_id,item_type,pic,coin,real_coin,quantity,sfc,region_sfc,rfc,region_rfc,create_at  from %s where create_at >= ? and create_at < ? and mod(gift_item_id, 50) = 0 ", "gift_items")
		rows, err := db.Query(ssql, start, end)
		if err != nil {
			log.Println(err)
			return
		}

		for rows.Next() {
			item := &GiftItem{}

			rows.Scan(
				&item.GiftID,
				&item.GiftItemID,
				&item.Name,
				&item.Sender,
				&item.Receive,
				&item.SenderUserID,
				&item.ReceiveUserID,
				&item.ItemType,
				&item.Picture,
				&item.Coin,
				&item.RealCoin,
				&item.Quantity,
				&item.SFC,
				&item.RegionSFC,
				&item.RFC,
				&item.RegionRFC,
				&item.CreateAt,
			)
			result = append(result, item)
			// log.Println(item.CreateAt)
		}
	})

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreateAt.Before(result[j].CreateAt)
	})

	return
}

func loadGiftItems() (result []*GiftItem) {
	f, err := os.Open("./data.gob")
	if err != nil {
		log.Panic(err)
	}
	reader, err := gzip.NewReader(f)
	if err != nil {
		log.Panic(err)
	}

	err = gob.NewDecoder(reader).Decode(&result)
	if err != nil {
		log.Panic(err)
	}

	return
}

func XTestDownload(t *testing.T) {
	now := time.Now()
	today := now.Truncate(time.Hour * 24)
	items := QueryGiftItems(today.Add(-time.Hour*24*5), today)

	// var items []*GiftItem
	// for _, item := range items2 {

	// 	items = append(items, item)

	// }

	f, err := os.OpenFile("./data.gob", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	var buf bytes.Buffer // Stand-in for a network connection
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(items)
	if err != nil {
		log.Panic(err)
	}
	log.Println(len(items), len(buf.Bytes()))

	gz := gzip.NewWriter(f)
	_, err = gz.Write(buf.Bytes())
	if err != nil {
		log.Panic(err)
	}

	err = gz.Close()
	if err != nil {
		log.Panic(err)
	}
}

type Item struct {
	Name      string
	Region    string
	Country   string
	Coin      int64
	ExtraCoin int64
}

func Test0(t *testing.T) {
	c := New()
	c.AddCategory("region", func(value interface{}) interface{} {
		item := value.(*Item)
		return item.Region
	}).AddCategory("country", func(value interface{}) interface{} {
		item := value.(*Item)
		return item.Country
	}).Collect()

	c.Put(&Item{
		Name:    "test",
		Region:  "Arab",
		Country: "USA",
	})

	c.Put(&Item{
		Name:    "test",
		Region:  "Arab2",
		Country: "USA",
	})

	c.Put(&Item{
		Name:    "test",
		Region:  "Arab2",
		Country: "CN",
	})

	c.Put(&Item{
		Name:    "test",
		Region:  "Arab",
		Country: "USA2",
	})

	var result, expect string
	result = fmt.Sprintf("%#v", c.Keys())
	expect = `[]interface {}{"Arab", "Arab2"}`
	if result != expect {
		t.Error(result, "!=", expect)
	}

	defer func() {
		if err := recover(); err == nil {
			t.Error("check Keys()")
		}
	}()

	for _, key := range c.Keys() {
		// log.Println(c.Keys(key))
		for _, key2 := range c.Keys(key) {
			c.Keys(key, key2)
			panic(nil)
		}
	}
	// c.debugPrint(0)
}

func Test1(t *testing.T) {
	c := New()
	c.AddCategory("region", func(value interface{}) interface{} {
		item := value.(*Item)
		return item.Region
	}).AddCategory("country", func(value interface{}) interface{} {
		item := value.(*Item)
		return item.Country
	}).Collect()

	c.Put(&Item{
		Name:    "test",
		Region:  "Arab",
		Country: "USA",
	})

	c.Put(&Item{
		Name:    "test",
		Region:  "Arab2",
		Country: "USA",
	})

	c.Put(&Item{
		Name:    "test",
		Region:  "Arab2",
		Country: "CN",
	})

	c.Put(&Item{
		Name:    "test",
		Region:  "Arab",
		Country: "USA2",
		Coin:    1,
	})
	c.Put(&Item{
		Name:    "test",
		Region:  "Arab",
		Country: "USA2",
		Coin:    5,
	})

	c.Put(&Item{
		Name:    "test",
		Region:  "Arab",
		Country: "USA2",
		Coin:    1,
	})

	// var result, expect string
	var items []*Item
	c.Get(&items, "Arab", "USA2")
	if len(items) != 3 {
		t.Error("len != 3 is error.")
		for _, item := range items {
			log.Printf("%#v", item)
		}
	}

	// c.Get(&items, "Arab")
	// for _, item := range items {
	// 	log.Printf("%#v", item)
	// }

	c.Get(&items)
	if len(items) != 6 {
		t.Error("len != 6 is error.")
		for _, item := range items {
			log.Printf("%#v", item)
		}
	}

	// c.debugPrint(0)
}

func Test2(t *testing.T) {

	clsfy := New()
	clsfy.Build(`region<RegionRFC>.country<RFC>.@Coin`)

	// items := loadGiftItemsJson()

	items := loadGiftItems()

	for _, item := range items {
		clsfy.Put(item)
	}

	var getitems []*GiftItem
	for _, region := range clsfy.Keys() {
		for _, country := range clsfy.Keys(region) {
			var gitems []*GiftItem
			clsfy.Get(&gitems, region, country)
			getitems = append(getitems, gitems...)
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].GiftItemID > items[j].GiftItemID
	})

	sort.Slice(getitems, func(i, j int) bool {
		return getitems[i].GiftItemID > getitems[j].GiftItemID
	})

	if len(items) != len(getitems) {
		t.Error("items.len != getitems.len")
		return
	}

	for i, item := range items {

		if item.GiftItemID != getitems[i].GiftItemID {
			t.Error("items != getitems", i)
			return
		}
	}
}

func Test3(t *testing.T) {

	clsfy := New()
	clsfy.Build(` region[1]. country[0].@ `,
		func(value interface{}) interface{} {
			return value.(*GiftItem).RFC
		},
		func(value interface{}) interface{} {
			return value.(*GiftItem).RegionRFC
		},
	)

	// items := loadGiftItemsJson()

	items := loadGiftItems()

	// log.Println(clsfy.Categorys())
	if clsfy.Categorys() != "region.country.@" {
		t.Errorf("error output Categorys %#v", clsfy.Categorys())
	}
	for _, item := range items {
		clsfy.Put(item)
	}

	// var getitems []*GiftItem
	for _, region := range clsfy.Keys() {
		for _, country := range clsfy.Keys(region) {
			var gitems []*GiftItem
			clsfy.Get(&gitems, region, country)
			// log.Println(region, country)
			for _, item := range gitems {
				// log.Printf("%#v", item)
				if item.RegionRFC != region {
					t.Errorf("Errro Region %s != %s", item.RegionRFC, region)
				}

				if item.RFC != country {
					t.Errorf("Errro Country %s != %s", item.RFC, country)
				}
			}
		}
	}

}

// 测试@
func Test4(t *testing.T) {

}
