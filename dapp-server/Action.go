package main

import (
	"encoding/hex"
	//"fmt"

	"github.com/ontio/ontology/common"
)

type Event interface {
	GetEventName() string
}

// 反转HexString
func reverseString2(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+2, to-2 {
		runes[from], runes[to-1] = runes[to-1], runes[from]
		runes[from+1], runes[to] = runes[to], runes[from+1]
	}
	return string(runes)
}

// 反转Bytearray
func reverseBytes(res []byte) []byte {
	//fmt.Println(res)
	for from, to := 0, len(res)-1; from < to; from, to = from+1, to-1 {
		res[from], res[to] = res[to], res[from]
	}
	//fmt.Println(res)
	return res
}

type CreateOrderEvent []interface{}

func (this CreateOrderEvent) GetEventName() string {
	bytes, _ := hex.DecodeString(this[0].(string))
	return string(bytes)
}
func (this CreateOrderEvent) GetOrder() Order {
	o := Order{}
	items := this[1].([]interface{})
	o.Id.SetString(reverseString2(items[0].(string)), 16)
	if bytes, err := hex.DecodeString(items[1].(string)); err == nil {
		o.Type = string(bytes)
	}
	if addr, err := common.AddressFromHexString(items[2].(string)); err == nil {
		o.Owner = addr.ToBase58()
	}
	o.Price.SetString(reverseString2(items[3].(string)), 16)
	o.Amount.SetString(reverseString2(items[4].(string)), 16)
	o.State.SetString(reverseString2(items[5].(string)), 16)
	o.PreId.SetString(reverseString2(items[6].(string)), 16)
	o.NextId.SetString(reverseString2(items[7].(string)), 16)
	o.UnAmount.SetString(reverseString2(items[8].(string)), 16)
	return o
}
