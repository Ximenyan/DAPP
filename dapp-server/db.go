package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func DBTest() {

	bi := big.Int{}
	bi.SetString("00989680", 16)
	bys, _ := json.Marshal([]big.Int{bi})
	fmt.Println(string(bys))
	//创建并打开数据库
	db, err := leveldb.OpenFile("./db", nil)
	if err != nil {
		panic(err)
	}

	defer db.Close() //关闭数据库

	//写入5条数据
	db.Put([]byte("key1"), []byte("value1"), nil)
	db.Put([]byte("key2"), []byte("value2"), nil)
	db.Put([]byte("key3"), []byte("value3"), nil)
	db.Put([]byte("key4"), []byte("value4"), nil)
	db.Put([]byte("key5"), []byte("value5"), nil)

	//循环遍历数据
	fmt.Println("循环遍历数据")
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		fmt.Printf("key:%s, value:%s\n", iter.Key(), iter.Value())
	}
	iter.Release()

	//读取某条数据
	data, _ := db.Get([]byte("key2"), nil)
	fmt.Printf("读取单条数据key2:%s\n", data)

	//批量写入数据
	batch := new(leveldb.Batch)
	batch.Put([]byte("key6"), []byte(strconv.Itoa(10000)))
	batch.Put([]byte("key7"), []byte(strconv.Itoa(20000)))
	batch.Delete([]byte("key4"))
	db.Write(batch, nil)

	//查找数据
	key := "key7"
	iter = db.NewIterator(nil, nil)
	for ok := iter.Seek([]byte(key)); ok; ok = iter.Next() {
		fmt.Printf("查找数据:%s, value:%s\n", iter.Key(), iter.Value())
	}
	iter.Release()

	//按key范围遍历数据
	fmt.Println("按key范围遍历数据")
	iter = db.NewIterator(&util.Range{Start: []byte("key3"), Limit: []byte("key7")}, nil)
	for iter.Next() {
		fmt.Printf("key:%s, value:%s\n", iter.Key(), iter.Value())
	}
	iter.Release()
}

//Push
func PushOrder(o *Order) {
	db, err := leveldb.OpenFile("./database/all", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close() //关闭数据库
	source, _ := json.Marshal(o)
	db.Put(o.Id.Bytes(), source, nil)
}

//Deleteb by ID
func DeleteOrder(id *big.Int) error {
	db, err := leveldb.OpenFile("./database/all", nil)
	if err != nil {
		return err
	}
	defer db.Close() //关闭数据库
	key := id.Bytes()

	ok := db.Delete(key, nil)
	return ok
}

//排序
func RankOrder(id *big.Int) {

}

//查询
func GetOrder(id *big.Int) error {
	db, err := leveldb.OpenFile("./database/all", nil)
	if err != nil {
		return err
	}
	defer db.Close() //关闭数据库
	key := id.Bytes()
	data, ok := db.Get(key, nil)
	return ok
}
