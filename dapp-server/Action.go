package main

import (
	"encoding/hex"
	"strconv"

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
	for from, to := 0, len(res)-1; from < to; from, to = from+1, to-1 {
		res[from], res[to] = res[to], res[from]
	}
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
	if s, err := strconv.ParseUint(reverseString2(items[0].(string)), 16, 64); err == nil {
		o.Id = s
	}
	if bytes, err := hex.DecodeString(items[1].(string)); err == nil {
		o.Type = string(bytes)
	}
	if addr, err := common.AddressFromHexString(items[2].(string)); err == nil {
		o.Owner = addr.ToBase58()
	}
	//fmt.Println(items[3].(string))
	//fmt.Println(reverseString2(items[3].(string)))
	if s, err := strconv.ParseUint(reverseString2(items[3].(string)), 16, 64); err == nil {
		o.Price = s
	}
	if s, err := strconv.ParseUint(reverseString2(items[4].(string)), 16, 64); err == nil {
		o.Amount = s
	}
	if s, err := strconv.ParseUint(reverseString2(items[5].(string)), 16, 64); err == nil {
		o.State = s
	}
	if s, err := strconv.ParseUint(reverseString2(items[6].(string)), 16, 64); err == nil {
		o.PreId = s
	}
	if s, err := strconv.ParseUint(reverseString2(items[7].(string)), 16, 64); err == nil {
		o.NextId = s
	}
	if s, err := strconv.ParseUint(reverseString2(items[8].(string)), 16, 64); err == nil {
		o.UnAmount = s
	}
	return o
}
