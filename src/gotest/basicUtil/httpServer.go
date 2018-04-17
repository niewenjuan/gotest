package basicUtil
//
import (
	"fmt"
	"net/http"
	//"encoding/json"
	qy1 "gotest/schema/protoc"
	"log"
	//"github.com/tedsuo/rata"
)

const (
	serviceName  = "GoServerHoy"
	schemaId = "hellworld"
	rpcOperationIdone = "LoadBalanceTest"
	delaytime = 1000
	listenPort = "2099"
	listenAddr = ":"+listenPort
)

func SimpleInvokeServer(w http.ResponseWriter, req *http.Request) {
	reply :=HighwayInvoke(serviceName, schemaId, rpcOperationIdone, &qy1.CseInvokeRequest{Delaytime:delaytime}, &qy1.CseInvokeReply{})
	PrintlnGreen("Highway Invoke")
	log.Println(reply)
	fmt.Fprint(w,string("test"))
}

func HttpGate(){

	http.HandleFunc("/hello", SimpleInvokeServer)

	//petRoutes := rata.Routes{
	//	{Name: "SimpleInvokeServer", Method: rata.GET, Path: "/testcommunication/GET/sayhello/:id"},
	//	//{Name: "sayhi", Method: rata.POST, Path: "/testcommunication/POST/sayhi"},
	//	//{Name: "sayjson", Method: rata.POST, Path: "/testcommunication/POST/sayjson"},
	//}
	//petHandlers := rata.Handlers{
	//	"SimpleInvokeServer": http.HandlerFunc(SimpleInvokeServer),
	//
	//}
	//router, err := rata.NewRouter(petRoutes, petHandlers)
	//if err != nil {
	//	panic(err)
	//}
	listenAndserve()
}

func listenAndserve(){
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}