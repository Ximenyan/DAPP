package main

import (
	"encoding/json"
	"log"
	"math/big"

	"github.com/boltdb/bolt"
)

var ALL_DB *bolt.DB

func BoltDBInit() {
	var err error
	ALL_DB, err = bolt.Open("./db/order.db", 0600, nil)
	if err != nil {
		log.Println(err)
	}
}

//Push Order
func BoltPushOrder(o *Order) {
	ALL_DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("all_order"))
		if err != nil {
			return err
		}
		if err = b.Put(o.Id.Bytes(), o.ToBytes()); err != nil {
			return err
		}
		b, err = tx.CreateBucketIfNotExists([]byte(o.Type + o.Business))
		if err = b.Put(o.Id.Bytes(), o.ToBytes()); err != nil {
			return err
		}
		b, err = tx.CreateBucketIfNotExists([]byte(o.Owner))
		if err = b.Put(o.Id.Bytes(), o.ToBytes()); err != nil {
			return err
		}
		return nil
	})
}

//Deleteb by ID
func BoltDeleteOrder(id *big.Int) {
	o := GetOrdersByIntId(id.Int64())
	//删除order
	ALL_DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(o.Type + o.Business))
		if err != nil {
			return err
		}
		if err = b.Delete(id.Bytes()); err != nil {
			return err
		}
		b, err = tx.CreateBucketIfNotExists([]byte(o.Owner))
		if err != nil {
			return err
		}
		if err = b.Delete(id.Bytes()); err != nil {
			return err
		}
		return nil
	})
}

//排序
func BoltRankOrder(id *big.Int) {

}

//查询 BoltGetOrderByID
func BoltGetOrderByID(id *big.Int) (*Order, error) {
	o := &Order{}
	var err error
	ALL_DB.View(func(tx *bolt.Tx) error {
		var b *bolt.Bucket
		b = tx.Bucket([]byte("all_order"))
		res := b.Get(id.Bytes())
		err = json.Unmarshal(res, o)
		return err
	})
	log.Println(err, o)
	return o, err
}

func BoltGetOrderByOwner(addr string) ([]*Order, error) {
	var err error
	orders := []*Order{}
	ALL_DB.View(func(tx *bolt.Tx) error {
		var b *bolt.Bucket
		b = tx.Bucket([]byte(addr))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			o := &Order{}
			err = json.Unmarshal(v, &o)
			orders = append(orders, o)
			log.Println(err, o)
		}

		return err
	})
	return orders, err
}

//查询指定类型的卖单
func BoltGetSellOrderByType(strType string) ([]*Order, error) {
	var err error
	orders := []*Order{}
	ALL_DB.View(func(tx *bolt.Tx) error {
		var b *bolt.Bucket
		b = tx.Bucket([]byte(strType + "_SELL_"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			o := &Order{}
			err = json.Unmarshal(v, &o)
			orders = append(orders, o)
			log.Println(err, o)
		}

		return err
	})
	return orders, err
}

//查询指定类型的买单
func BoltGetBuyOrderByType(strType string) ([]*Order, error) {
	var err error
	orders := []*Order{}
	ALL_DB.View(func(tx *bolt.Tx) error {
		var b *bolt.Bucket
		b = tx.Bucket([]byte(strType + "_BUY_"))
		if err != nil {
			return err
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			o := &Order{}
			err = json.Unmarshal(v, &o)
			orders = append(orders, o)
			log.Println(err, o)
		}

		return err
	})
	return orders, err
}
