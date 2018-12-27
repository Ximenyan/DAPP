package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/ontio/ontology/common"
)
var tair_order_key = []string{
	"_BUY___List_Min_Order___ONG_ONT_",
	"_SELL___List_Min_Order___ONG_ONT_",
}
//代币小数精确到x分之1
var DECIMALS_MAPS = map[string]int{
	"_ONG_ONT_":100000000}

type Order struct {
	Id       uint64
	Type     string
	Owner    string
	Price    uint64
	Amount   uint64
	State    uint64
	PreId    uint64
	NextId   uint64
	UnAmount uint64
}

type Ranking struct {
	Price    uint64
	UnAmount uint64
}

func IntToByte(num int64) []byte {
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, num)
	var start int
	for i := 0; i < len(buffer.Bytes()); i++ {
		start = i
		if buffer.Bytes()[i] != 0 {
			break
		}
	}
	return buffer.Bytes()[start:]
}

func bufUint64(buf []byte) []byte {
	res := []byte{}
	for i := 0; i < 8-len(buf); i++ {
		res = append(res, 0)
	}
	for j := len(buf); j >= 1; j-- {
		res = append(res, buf[j-1])
	}
	return res
}
func CreateOrder(buf []byte) (order *Order, err error) {
	defer func() {
		if p := recover(); p != nil {
			order = nil
			err = errors.New("create order error!")
		}
	}()
	order = new(Order)
	if buf[3] == 1 {
		order.Id = binary.BigEndian.Uint64(bufUint64(buf[4 : 4+buf[3]]))
	}
	tmp := 4 + buf[3]
	order.Type = string(buf[tmp+2 : tmp+2+buf[tmp+1]])
	tmp = tmp + 2 + buf[tmp+1]
	addr, _ := common.AddressParseFromBytes(buf[tmp+2 : tmp+22])
	order.Owner = addr.ToBase58()
	tmp = tmp + 22
	order.Price = binary.BigEndian.Uint64(bufUint64(buf[tmp+2 : tmp+2+buf[tmp+1]]))
	tmp = tmp + 2 + buf[tmp+1]
	order.Amount = binary.BigEndian.Uint64(bufUint64(buf[tmp+2 : tmp+2+buf[tmp+1]]))
	tmp = tmp + 2 + buf[tmp+1]
	order.State = binary.BigEndian.Uint64(bufUint64(buf[tmp+2 : tmp+2+buf[tmp+1]]))
	tmp = tmp + 2 + buf[tmp+1]
	order.PreId = binary.BigEndian.Uint64(bufUint64(buf[tmp+2 : tmp+2+buf[tmp+1]]))
	tmp = tmp + 2 + buf[tmp+1]
	order.NextId = binary.BigEndian.Uint64(bufUint64(buf[tmp+2 : tmp+2+buf[tmp+1]]))
	tmp = tmp + 2 + buf[tmp+1]
	order.UnAmount = binary.BigEndian.Uint64(bufUint64(buf[tmp+2 : tmp+2+buf[tmp+1]]))
	return order, nil
}
var orders_type_map = make(map[string][]*Order)

//通过RPC 和 订单ID 查询所有订单
func GetOrdersByIntId(id int64) *Order{
		key := append([]byte("__ORDER___"), reverseBytes(IntToByte(int64(id)))...)
		//key = []byte("_BUY___List_Min_Order___ONG_ONT_")
		//fmt.Println(key)
		res, _ := ONT.GetStorage(CONTRACT_ADDR.ToHexString(), key)
		//fmt.Println(res)
		order, err := CreateOrder(res)
		//fmt.Println(order)
		if err != nil {
			return nil
		}
		return order
}

//通过RPC 和 订单ID 查询所有订单
func GetOrdersByBytesId(id []byte) *Order{
	key := append([]byte("__ORDER___"), id...)
	res, _ := ONT.GetStorage(CONTRACT_ADDR.ToHexString(), key)
	order, err := CreateOrder(res)
	if err != nil {
		return nil
	}
	return order
}
//根据订单类型抓取所有订单
func GetAllOrdersByType(strType string) []*Order {
	res, _ := ONT.GetStorage(CONTRACT_ADDR.ToHexString(), []byte(strType))
	orders := []*Order{}
	for len(res)>=1 {
		order := GetOrdersByBytesId(res)
		orders = append(orders,order)
		if order != nil{
			res = reverseBytes(IntToByte(int64(order.PreId)))
		}else{
			res = []byte{}
		}
	}
	orders_type_map[strType] = orders
	return orders
}

//获取排行数据
func GetOrdersRankByType(strType string,top_num int) []Ranking {
	var orders []*Order
	if _,ok := orders_type_map[strType]; !ok{
		orders = GetAllOrdersByType(strType)
	}else{
		orders = orders_type_map[strType]
	}

	ranks := []Ranking{}
	if len(orders) < 1{
		return ranks
	}
	var price uint64 = orders[0].Price
	var  unamount uint64= 0
	for i:=0;i<len(orders);i++{
		o := orders[i]
		if o == nil{
			ranks = append(ranks,Ranking{price,unamount})
			continue
		}
		if price != o.Price{
			ranks = append(ranks,Ranking{price,unamount})
			unamount = o.UnAmount
			price = o.Price
		}else{
			unamount += o.UnAmount
		}
		if len(ranks) >= top_num{
			return ranks
		}
	}
	return ranks
}