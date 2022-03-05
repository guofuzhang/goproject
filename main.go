package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"strconv"
	"time"
)

// 商品表
type Goods struct {
	Id      uint   `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	Name    string `gorm:"column:name;type:varchar(50);NOT NULL" json:"name"`   // 名称
	Count   int    `gorm:"column:count;type:int(11);NOT NULL" json:"count"`     // 库存
	Sale    int    `gorm:"column:sale;type:int(11);NOT NULL" json:"sale"`       // 已售
	Version int    `gorm:"column:version;type:int(11);NOT NULL" json:"version"` // 乐观锁，版本号
}

// 订单表
type GoodsOrder struct {
	Id         uint      `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	Gid        int       `gorm:"column:gid;type:int(11);NOT NULL" json:"gid"`                                             // 库存ID
	Name       string    `gorm:"column:name;type:varchar(30);NOT NULL" json:"name"`                                       // 商品名称
	CreateTime time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_time"` // 创建时间
}

//实际表名
func (m *GoodsOrder) TableName() string {
	return "goods_order"
}

func main() {

	http.HandleFunc("/", addOrder)
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func getDb() *gorm.DB {
	connArgs := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", "guofu", "guofu", "localhost", 13306, "go-project")

	db, err := gorm.Open("mysql", connArgs)
	if err != nil {
		panic(err)
	}
	db.LogMode(false) //打印sql语句
	//开启连接池
	db.DB().SetMaxIdleConns(100)   //最大空闲连接
	db.DB().SetMaxOpenConns(100)   //最大连接数
	db.DB().SetConnMaxLifetime(30) //最大生存时间(s)
	return db
}

func addOrder(w http.ResponseWriter, r *http.Request) {
	value := GetUUID()
	key := "order"
	client := getRedis()
	defer client.Close()
	cmd := client.SetNX(key, value, time.Second*30) //这里会有一个问题,就是 我里面的程序执行过长,导致锁释放,那么程序末尾的删除锁 就会删除其他请求的锁,导致不可用
	if cmd.Val() == true {

		db := getDb()
		defer db.Close()

		// 先去查看商品表还有没有库存
		var goods Goods
		db.Where("id = ?", "1").First(&goods)
		fmt.Println(goods.Count)
		if goods.Count > 0 {
			tx := db.Begin()
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()

			goods.Sale += 1
			goods.Count -= 1
			//更新数据库
			if err := tx.Save(&goods).Error; err != nil {
				tx.Rollback()
				panic(err)
			}

			order := GoodsOrder{
				Gid:  1,
				Name: strconv.Itoa(int(time.Now().Unix())),
			}

			if err := tx.Create(&order).Error; err != nil {
				tx.Rollback()
				panic(err)
			}
			tx.Commit()
			w.Write([]byte(fmt.Sprintf("the count i read is %d", goods.Count)))
		} else {
			w.Write([]byte("我啥子都么抢到"))

		}
		if client.Get(key).Val() == value {
			client.Del(key)
		}
	}
}

func GetUUID() string {
	u2 := uuid.NewV4()
	return u2.String()
}

func getRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	return client
}
