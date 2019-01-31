package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"log"
	"math/big"
	"strings"

	"github.com/ontio/ontology/common"
)

var tair_order_key = []string{
	"_BUY___List_Min_Order___ONG_ONT_",
	"_SELL___List_Min_Order___ONG_ONT_",
}

//代币小数精确到x分之1
var DECIMALS_MAPS = map[string]int{
	"_ONG_ONT_": 100000000}

type Order struct {
	Id       big.Int
	Type     string
	Business string
	Owner    string
	Price    big.Int
	Amount   big.Int
	State    big.Int
	PreId    big.Int
	NextId   big.Int
	UnAmount big.Int
}

func (o *Order) ToString() string {

	res, err := json.Marshal(o)
	if err != nil {
		log.Println(err, o)
		return ""
	}
	return string(res)
}
func (o *Order) ToBytes() []byte {

	res, err := json.Marshal(o)
	if err != nil {
		log.Println(err, o)
		return nil
	}
	return res
}

type Ranking struct {
	Price    big.Int
	UnAmount *big.Int
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

//根据BYTES 还原一个Order
func OrderDeseriallization(buf []byte) (order *Order, err error) {
	defer func() {
		if p := recover(); p != nil {
			order = nil
			err = errors.New("create order error!")
		}
	}()
	order = new(Order)
	order.Id.SetBytes(reverseBytes(buf[4 : 4+buf[3]]))
	tmp := 4 + buf[3]
	order.Type = string(buf[tmp+2 : tmp+2+buf[tmp+1]])
	tmp = tmp + 2 + buf[tmp+1]
	order.Business = string(buf[tmp+2 : tmp+2+buf[tmp+1]])
	tmp = tmp + 2 + buf[tmp+1]
	addr, _ := common.AddressParseFromBytes(buf[tmp+2 : tmp+22])
	order.Owner = addr.ToBase58()
	tmp = tmp + 22
	order.Price.SetBytes(reverseBytes(buf[tmp+2 : tmp+2+buf[tmp+1]]))
	tmp = tmp + 2 + buf[tmp+1]
	order.Amount.SetBytes(reverseBytes(buf[tmp+2 : tmp+2+buf[tmp+1]]))

	tmp = tmp + 2 + buf[tmp+1]
	order.State.SetBytes(reverseBytes(buf[tmp+2 : tmp+2+buf[tmp+1]]))

	tmp = tmp + 2 + buf[tmp+1]
	order.PreId.SetBytes(reverseBytes(buf[tmp+2 : tmp+2+buf[tmp+1]]))

	tmp = tmp + 2 + buf[tmp+1]
	order.NextId.SetBytes(reverseBytes(buf[tmp+2 : tmp+2+buf[tmp+1]]))

	tmp = tmp + 2 + buf[tmp+1]
	order.UnAmount.SetBytes(reverseBytes(buf[tmp+2 : tmp+2+buf[tmp+1]]))

	//fmt.Println(order)
	return order, nil
}

var orders_type_map = make(map[string][]*Order)

//通过RPC 和 订单ID 查询所有订单
func GetOrdersByIntId(id int64) *Order {
	key := append([]byte("__ORDER___"), reverseBytes(IntToByte(int64(id)))...)
	res, _ := ONT.GetStorage(CONTRACT_ADDR.ToHexString(), key)
	order, err := OrderDeseriallization(res)
	if err != nil {
		return nil
	}
	return order
}

//通过RPC 和 订单ID 查询所有订单
func GetOrdersByBytesId(id []byte) *Order {
	key := append([]byte("__ORDER___"), id...)
	res, _ := ONT.GetStorage(CONTRACT_ADDR.ToHexString(), key)
	order, err := OrderDeseriallization(res)
	if err != nil {
		return nil
	}
	BoltPushOrder(order)
	return order
}

//根据订单类型抓取所有订单
func GetAllOrdersByType(strType string) []*Order {
	res, _ := ONT.GetStorage(CONTRACT_ADDR.ToHexString(), []byte(strType))

	orders := []*Order{}
	for len(res) >= 1 {
		order := GetOrdersByBytesId(res)
		BoltPushOrder(order)
		log.Println(order)
		orders = append(orders, order)
		if order != nil {
			res = reverseBytes(order.PreId.Bytes())
		} else {
			res = []byte{}
		}
	}
	orders_type_map[strType] = orders
	return orders
}

//获取应该插入的位置，可能会有错
func GetIndexOrder(order_type string, iPrice uint64) (pre, next uint64) {
	price := big.NewInt(int64(iPrice))
	res, err := ONT.GetStorage(CONTRACT_ADDR.ToHexString(), []byte(order_type))
	if err != nil || len(res) == 0 {
		return 0, 0
	}
	for {
		order := GetOrdersByBytesId(res)
		if strings.Index(order_type, "_BUY_") >= 0 {
			if order.Price.Cmp(price) <= 0 {
				return order.Id.Uint64(), order.NextId.Uint64()
			}
		} else if strings.Index(order_type, "_SELL_") >= 0 {
			if order.Price.Cmp(price) >= 0 {
				return order.Id.Uint64(), order.NextId.Uint64()
			}
		}
		if order.PreId.Int64() == 0 {
			return 0, order.Id.Uint64()
		}
		res = reverseBytes(order.PreId.Bytes())
	}
}

//获取排行数据
func GetOrdersRankByType(strType string, top_num int) []Ranking {
	var orders []*Order
	orders = GetAllOrdersByType(strType)
	ranks := []Ranking{}
	if len(orders) < 1 {
		return ranks
	}
	var price big.Int = orders[0].Price
	var unamount *big.Int = big.NewInt(0)
	for i := 0; i < len(orders); i++ {
		o := orders[i]
		if o == nil {
			ranks = append(ranks, Ranking{price, unamount})
			continue
		}
		if price.Cmp(&o.Price) != 0 {
			ranks = append(ranks, Ranking{price, unamount})
			unamount = &o.UnAmount
			price = o.Price
		} else {
			unamount = unamount.Add(unamount, &o.UnAmount)
		}
		if len(ranks) >= top_num {
			return ranks
		}
		if i == len(orders)-1 {
			ranks = append(ranks, Ranking{price, unamount})
		}
	}
	return ranks
}