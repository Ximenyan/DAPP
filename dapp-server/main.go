package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ontio/ontology-go-sdk"
	sdk "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology/common"
)

var CONTRACT_ADDR, _ = common.AddressFromHexString("cf5e452798963fbb9f40652ffda7ddfcfd6db594")
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
						fmt.Println(create.GetEventName())
						fmt.Println(create.GetOrder())
					}
				}
			}
		}
	}()

}

type timeHandler struct {
}

func (th *timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	orders := GetOrdersRankByType(r.Form.Get("qurey_type"),10)
	bytes, _ := json.Marshal(orders)
	w.Write([]byte(bytes))
}

func main() {
	//&{d37b6553c7fdc61ebb10bb02ee968f845b66bd7902e70c9bff4bc4b0fcd43d89 1 11929500 [0xc4201387c0 0xc420138880 0xc420138940]}
	CreateONT()
	OntConnect()
	fmt.Println(ONT.GetTransaction("d37b6553c7fdc61ebb10bb02ee968f845b66bd7902e70c9bff4bc4b0fcd43d89"))
	mux := http.NewServeMux()
	th := &timeHandler{}
	mux.Handle("/", th)
	http.ListenAndServe(":3030", mux)
}
