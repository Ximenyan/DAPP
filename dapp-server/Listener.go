package main

import (
	"log"
)

var ALL_TYPE []string

func Listenerinit() {
	CreateONT()
	OntConnect()
	ALL_TYPE = []string{
		`_SELL___List_Tail_Order___ONG_ONT_`,
		`_BUY___List_Tail_Order___ONG_ONT_`,
	}
	go SyncingAllOrder()
}
func SyncingAllOrder() {
	i := 1
	for {
		o := GetOrdersByBytesId(IntToByte(int64(i)))
		if o == nil {
			return
		}
		log.Println("SyncingAllOrder::::", o.ToString())
		i++
	}
}
