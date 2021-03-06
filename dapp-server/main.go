package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ontio/ontology-go-sdk"
	sdk "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology/common"
)

//722da66717247fa6dd8dda4b5cef4535509855bb
var CONTRACT_ADDR, _ = common.AddressFromHexString("e2480be800fa3a1ad859d6518a115de02f7dceea")
var ONT *ontology_go_sdk.OntologySdk

func CreateONT() {
	ONT = ontology_go_sdk.NewOntologySdk()
	ONT.NewRpcClient().SetAddress("http://13.78.112.191:20336")
}

func OntConnect() {
	CL := ONT.NewWebSocketClient()
	CL.Connect("ws://13.78.112.191:20335")

	go func() {
		CL.AddContractFilter(CONTRACT_ADDR.ToHexString())
		CL.SubscribeEvent()
		for {
			select {
			case q := <-CL.GetActionCh():
				events := q.Result.(*sdk.SmartContactEvent)
				for _, event := range events.Notify {
					if event.ContractAddress == CONTRACT_ADDR.ToHexString() {
						create := CreateOrderEvent(event.States.([]interface{}))
						o := create.GetOrder()
						BoltPushOrder(o)
						fmt.Println(o)
						fmt.Println(create.GetEventName())
					}
				}
			}
		}
	}()
}

type ServerHandler struct {
}

func (th *ServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	switch r.Form.Get("req_type") {
	case "query_order_rank":
		orders := GetOrdersRankByType(r.Form.Get("order_type"), 10)
		bytes, _ := json.Marshal(orders)
		w.Write([]byte(bytes))
		break
	case "create_order":
		order_type := r.Form.Get("order_type")
		price, err := strconv.ParseUint(r.Form.Get("price"), 10, 64)
		if err != nil {
			return
		}
		pre, next := GetIndexOrder(order_type, price)
		bytes, _ := json.Marshal([]uint64{pre, next})
		w.Write(bytes)
		break
	default:
		break
	}
}

func main() {
	//DBTest()
	//fmt.Println(ONT.GetStorage(CONTRACT_ADDR.ToHexString(), append([]byte("__ORDER___"), 2)))
	//fmt.Println(ONT.Native.Ong.TotalSupply())
	//fmt.Println(ONT.GetStorage(CONTRACT_ADDR.ToHexString(), []byte("_BUY___List_Head_Order___ONG_ONT_")))
	//GetOrdersRankByType("_SELL___List_Tail_Order___ONG_ONT_", 10)
	mux := http.NewServeMux()
	th := &ServerHandler{}
	mux.Handle("/", th)
	http.ListenAndServe(":3030", mux)
}
