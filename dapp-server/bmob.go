package main

import (
	"encoding/json"
	"log"

	"github.com/bmob/bmob-go-sdk"
)

var (
	appConfig = bmob.RestConfig{"b18cda25d056ac3a6e22f6a304cb37b8",
		"f60c9cd10b04fa2a5fd6f914c64c6528"}
)

func PushOrder(order *Order) error {
	bytes, _ := json.Marshal(order)
	header, err := bmob.DoRestReq(appConfig,
		bmob.RestRequest{
			bmob.BaseReq{
				"POST",
				bmob.ApiRestURL("Orders") + "/",
				""},
			"application/json",
			bytes},
		nil)
	if err == nil {
		log.Println(header)
	} else {
		log.Println(err)
	}
	log.Println("****************************************")
	return err
}
