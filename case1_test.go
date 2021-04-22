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

// GiftItems 查询礼物的数据
func GiftItems(start, end time.Time) (result []*GiftItem) {

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

func TestDownload(t *testing.T) {
	now := time.Now()
	today := now.Truncate(time.Hour * 24)
	items := GiftItems(today.Add(-time.Hour*24*5), today)

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
		Region:  "Arab",
		Country: "USA2",
	})
	log.Println(c.Keys())
	c.debugPrint(0)
}

func Test2(t *testing.T) {

	clsfy := New()
	clsfy.Build(`region<RegionRFC>.country<RFC>.@Coin`)

	// items := loadGiftItemsJson()

	items := loadGiftItems()

	for _, item := range items {
		// log.Println(item.RegionRFC, item.SFC)
		clsfy.Put(item)
	}

	log.Println(clsfy.Categorys())
	log.Println(len(items), clsfy.Keys(), clsfy.Categorys())

	for _, region := range clsfy.Keys() {
		log.Println(region)
		for _, country := range clsfy.Keys(region) {
			log.Println(country)
			var items []*GiftItem
			clsfy.Get(&items, region, country)
			log.Println(items[0].RFC, items[0].RegionSFC)
		}
	}
}
