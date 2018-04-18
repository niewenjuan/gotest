package schemaForClient

import (
	"gotest/schema/protoc"
	"github.com/ServiceComb/go-chassis/server/restful"
	"encoding/json"
	q "gotest/basicUtil"
	"net/http"

	"time"
	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core"
	//"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/metadata"
	"golang.org/x/net/context"
	"strconv"
	"log"
	"fmt"
	"strings"
)
var restInvoker = core.NewRestInvoker()
var ctx = context.Background()
//var ctx = metadata.NewContext(context.Background(), map[string]string{"user-name":"Hoy",})

const(
	provider_scname = "GoServerHoy"
	provider_Mid_scname = "GoMidServer"
	provider_schema = "hellworld"
	provider_operationId = "LoadBalanceTest"
	provider_operationId_weight = "LoadBalanceTest_weight"
)
type RestFulConsumer struct {
}

type token struct {
	token string
	time time.Time
}

func (s *RestFulConsumer) SimpleInvokeServer(b *restful.Context) {

	reply := q.HighwayInvoke(provider_scname, provider_schema, provider_operationId, &cse.CseInvokeRequest{Delaytime:1000}, &cse.CseInvokeReply{})
	resp , _ := json.Marshal(reply)
	b.Write(resp)
}

func (s *RestFulConsumer) CycleInvokeServer(b *restful.Context) {
	var temp string
	var url string
	var delaytime int
	cycleNum, _ :=strconv.Atoi(b.ReadPathParameter("num"))
	delaytime, _ =strconv.Atoi(b.ReadPathParameter("delayTime"))

	protocol := b.ReadPathParameter("protocol")
	if protocol == "rest" {
		url = q.CombinationURL("cse","%s%s",provider_scname,fmt.Sprintf("/demo/LoadBalanceTest/%s",strconv.Itoa(delaytime)))
		//restQuest, _ := rest.NewRequest(http.MethodGet,url,nil)
		for i := 0; i < cycleNum; i++ {
			result := q.RestfulInvoke(http.MethodGet, url, nil)
			temp = result +","+ temp
		}
	} else {
		for i := 0; i < cycleNum; i++ {
			resp, _ := json.Marshal(q.HighwayInvoke(provider_scname, provider_schema, provider_operationId, &cse.CseInvokeRequest{Delaytime:int64(delaytime)}, &cse.CseInvokeReply{}))
			dat := q.JsonExchange(resp)
			temp = dat["instant_id"] + "," + temp
		}
	}
	b.Write([]byte(temp))
}

//func (s *RestFulConsumer) CycleInvokeServerForQps(b *restful.Context) {
//	var temp string
//	var url string
//	var delaytime int
//	cycleNum, _ :=strconv.Atoi(b.ReadPathParameter("num"))
//	delaytime, _ =strconv.Atoi(b.ReadPathParameter("delayTime"))
//
//	protocol := b.ReadPathParameter("protocol")
//	if protocol == "rest" {
//		url = q.CombinationURL("cse","%s%s",provider_scname,"/demo/QpsTest")
//		restQuest, _ := rest.NewRequest(http.MethodGet,url, nil)
//		for i := 0; i < cycleNum; i++ {
//			result, err := restInvoker.ContextDo(ctx,restQuest)
//			if err == nil {
//			log.Println(string(result.ReadBody()))
//			temp =  temp+"\n"+string(result.ReadBody())
//			}else{
//				log.Println(err)
//			}
//		}
//	} else {
//		for i := 0; i < cycleNum; i++ {
//			resp, _ := json.Marshal(q.HighwayInvoke(provider_scname, provider_schema, provider_operationId, &cse.CseInvokeRequest{Delaytime:int64(delaytime)}, &cse.CseInvokeReply{}))
//			dat := q.JsonExchange(resp)
//			temp = dat["instant_id"] + "\n" + temp
//		}
//	}
//	b.Write([]byte(temp))
//}

func (s *RestFulConsumer) CycleInvokeServerForQps(b *restful.Context) {
	var temp string
	//var url string
	//var delaytime int
	cycleNum, _ :=strconv.Atoi(b.ReadPathParameter("num"))
	//delaytime, _ =strconv.Atoi(b.ReadPathParameter("delayTime"))

	protocol := b.ReadPathParameter("protocol")
	if protocol == "rest" {
		//url = q.CombinationURL("cse","%s%s",provider_scname,"/demo/LoadBalanceTest/0")
		//restQuest, _ := rest.NewRequest(http.MethodGet,url, nil)
		for i := 0; i < cycleNum; i++ {
			result := q.RestfulInvoke(http.MethodGet, "cse://GoServerHoy/demo/QpsTest", nil)
			temp =  temp+","+ result
			log.Println(result)
		}
	} else {
		for i := 0; i < cycleNum; i++ {
			resp, _ := json.Marshal(q.HighwayInvoke(provider_scname, provider_schema, "QpsTest", &cse.CseInvokeRequest{}, &cse.CseInvokeReply{}))
			dat := q.JsonExchange(resp)
			temp = temp + "," +dat["instant_id"]
			log.Println(dat["instant_id"])
		}
	}
	b.Write([]byte(temp))
}

//会话
func (s *RestFulConsumer) CycleInvokeServer_session(b *restful.Context) {
	var temp string
	var url string
	var delaytime int
	var cookies string
	var cookieId []string
	cycleNum, _ :=strconv.Atoi(b.ReadPathParameter("num"))
	delaytime, _ =strconv.Atoi(b.ReadPathParameter("delayTime"))
	//protocol := b.ReadPathParameter("protocol")

	url = q.CombinationURL("cse","%s%s",provider_scname,fmt.Sprintf("/demo/LoadBalanceTest/%s",strconv.Itoa(delaytime)))
	restQuest, _ := rest.NewRequest(http.MethodGet,url, nil)
	result, _ := restInvoker.ContextDo(ctx,restQuest)
	cookies = string(result.GetCookie("ServiceCombLB"))
	log.Println("first get :",cookies)
	temp = string(result.ReadBody())+","+ temp
	restQuest.Close()
	result.Close()
	cookieId = strings.Split(cookies,"=")
	for i := 0; i < cycleNum; i++ {

		//if(cookies != ""){

			//restQuest.SetCookie("ServiceCombLB",cookieId[1])
			result := q.RestfulInvokeSession(cookieId[1],http.MethodGet,"cse://GoServerHoy/demo/LoadBalanceTest/0",nil)
			temp = result+","+ temp
		//}else {
		//	result, _ := restInvoker.ContextDo(ctx,restQuest)
		//	cookies = string(result.GetCookie("ServiceCombLB"))
		//	log.Println(cookies)
		//	temp = string(result.ReadBody())+"\n"+ temp
		//	result.Close()
		//}
	}
	b.Write([]byte(temp))
}

//权值
func (s *RestFulConsumer) CycleInvokeServer_weight(b *restful.Context) {
	var temp string
	var url string
	cycleNum, _ :=strconv.Atoi(b.ReadPathParameter("num"))

	protocol := b.ReadPathParameter("protocol")
	if protocol == "rest" {
		url = q.CombinationURL("cse","%s%s",provider_scname,fmt.Sprintf("/demo/%s","LoadBalanceTestWeight"))

		for i := 0; i < cycleNum; i++ {
			restQuest, _ := rest.NewRequest(http.MethodGet,url, nil)
			result, _ := restInvoker.ContextDo(ctx,restQuest)
			temp = string(result.ReadBody()) +","+ temp
			result.Close()
			restQuest.Close()
		}
	} else {
		for i := 0; i < cycleNum; i++ {
			resp, _ := json.Marshal(q.HighwayInvoke(provider_scname, provider_schema, provider_operationId_weight, &cse.CseInvokeRequest{}, &cse.CseInvokeReply{}))
			dat := q.JsonExchange(resp)
			temp = dat["instant_id"] + "," + temp
		}
	}
	b.Write([]byte(temp))
}


func (s *RestFulConsumer) RestSimpleInvokeServer(b *restful.Context) {
	delaytime, _ :=strconv.Atoi(b.ReadPathParameter("delayTime"))
	u := q.CombinationURL("cse","%s%s",provider_scname,fmt.Sprintf("/demo/LoadBalanceTest/%s",strconv.Itoa(delaytime)))
	//log.Println("url is:",u)
	restQuest, _ := rest.NewRequest(http.MethodGet,u, nil)
	defer restQuest.Close()
	result, err1 := restInvoker.ContextDo(ctx,restQuest)
	defer result.Close()
	q.Isok(err1)
	b.WriteHeader(result.GetStatusCode())
	b.Write(result.ReadBody())
	return
}

func (s *RestFulConsumer) RestSimpleInvokeServerCrossApp(b *restful.Context) {
	serviceName := b.ReadPathParameter("serviceName")
	appId := b.ReadPathParameter("appId")
	delaytime, _ :=strconv.Atoi(b.ReadPathParameter("delayTime"))
	u := q.CombinationURL("cse","%s%s",serviceName,fmt.Sprintf("/demo/LoadBalanceTest/%s",strconv.Itoa(delaytime)))
	log.Println("url is:",u)
	log.Println("appid is:",appId)
	restQuest, _ := rest.NewRequest(http.MethodGet,u, nil)
	defer restQuest.Close()
	result, err1 := restInvoker.ContextDo(ctx,restQuest,core.WithAppID(appId))
	defer result.Close()
	q.Isok(err1)
	b.Write(result.ReadBody())
}

func (s *RestFulConsumer) CycleInvokeServerForQpsCrossApp(b *restful.Context) {
	var temp string
	cycleNum, _ :=strconv.Atoi(b.ReadPathParameter("num"))
	url1 := q.CombinationURL("cse","%s%s","GoMidServer","/demo/QpsTest")
	protocol := b.ReadPathParameter("protocol")
	if protocol == "rest" {
		//url = q.CombinationURL("cse","%s%s",provider_scname,"/demo/LoadBalanceTest/0")
		//restQuest, _ := rest.NewRequest(http.MethodGet,url, nil)
		for i := 0; i < cycleNum; i++ {
			result := q.RestfulInvoke(http.MethodGet, url1, nil,core.WithAppID("go"))
			temp =  temp+","+ result
			log.Println(result)
		}
	} else {
		for i := 0; i < cycleNum; i++ {
			resp, _ := json.Marshal(q.HighwayInvoke(provider_Mid_scname, provider_schema, "QpsTest", &cse.CseInvokeRequest{}, &cse.CseInvokeReply{},core.WithAppID("go")))
			dat := q.JsonExchange(resp)
			temp = temp + "," +dat["instant_id"]
			log.Println(dat["instant_id"])
		}
	}
	b.Write([]byte(temp))
}

func (s *RestFulConsumer) SimpleInvokeServerCrossApp(b *restful.Context) {
	serviceName := b.ReadPathParameter("serviceName")
	appId := b.ReadPathParameter("appId")
	schemaid := b.ReadPathParameter("schemaId")
	operationid := b.ReadPathParameter("operationId")
	reply := q.HighwayInvoke(serviceName, schemaid, operationid, &cse.CseInvokeRequest{Delaytime:1000}, &cse.CseInvokeReply{},core.WithAppID(appId))
	resp , _ := json.Marshal(reply)
	b.Write(resp)
}

func (s *RestFulConsumer) QueryToken(b *restful.Context) {
	mytoken := new(token)
	resp , _ := json.Marshal(mytoken)
	b.Write(resp)
}

func (s *RestFulConsumer) AutoCircuitRestFulTest(b *restful.Context){
	protocol := b.ReadPathParameter("procotol")
	if protocol == "rest" {
		b.Write([]byte(restAutoCircuitTest()))
	}else {
		b.Write([]byte(highwayAutoCircuit()))
	}

}


func (s *RestFulConsumer) CircuitFailInvoke(b *restful.Context){
	serviceName := b.ReadPathParameter("serviceName")
	protocol := b.ReadPathParameter("protocol")
	var resbody []byte
	if protocol == "rest" {
		u := q.CombinationURL("cse","%s%s",serviceName,"/failInvoke")
		restQuest, _ := rest.NewRequest(http.MethodGet,u, nil)
		result, _ := restInvoker.ContextDo(ctx,restQuest)
		defer restQuest.Close()
		defer result.Close()
		resbody = result.ReadBody()
	} else{
		reply := q.HighwayInvoke(serviceName, provider_schema, provider_operationId, &cse.CseInvokeRequest{Delaytime:1000}, &cse.CseInvokeReply{})
		resp , _ := json.Marshal(reply)
		resbody = resp
	}
	b.Write(resbody)
}

func (s *RestFulConsumer) URLPatterns() []restful.Route {
	return []restful.Route{
		{http.MethodGet, q.C_SimpleInvoke_Url, q.C_SimpleInvoke_OperationId},
		{http.MethodGet, q.C_RestSimpleInvoke_Url, q.C_RestSimpleInvoke_OperationId},
		{http.MethodGet, q.C_CrossAppRestSimpleInvoke_Url, q.C_CrossAppRestSimpleInvoke_OperationId},
		{http.MethodGet, q.C_CrossAppQpsInvoke_Url, q.C_CrossAppQpsInvoke_OperationId},
		{http.MethodGet, q.C_CrossAppHighwaySimpleInvoke_Url, q.C_CrossAppHighwaySimpleInvoke_OperationId},
		{http.MethodGet, q.C_CycleInvokeQps_Url, q.C_CycleInvokeQps_OperationId},
		{http.MethodGet, q.C_CycleInvoke_Url, q.C_CycleInvoke_OperationId},
		{http.MethodGet, q.C_CycleInvoke_Weight_Url, q.C_CycleInvoke_Weight_OperationId},
		{http.MethodGet, q.C_CycleInvoke_Session_Url, q.C_CycleInvoke_Session_OperationId},
		{http.MethodGet, q.C_QueryToken_Url, q.C_QueryToken_OperationId},
		{http.MethodGet, q.C_AutoCircuit_Url, q.C_AutoCircuit_OperationId},
		{http.MethodGet, q.C_CircuitFailInvoke_Url, q.C_CircuitFailInvoke_OperationId},
	}
}

func restAutoCircuitTest()string{
	q.PrintlnCyan("rest Invoke")
	result := "test"
	for sum := 0; sum < 15; sum++ {
		q.PrintlnRed("第", sum, "次")
		if sum > 4 {
			result = result+ "@" + q.RestfulInvoke(http.MethodGet, "cse://GoServerHoy/demo/LoadBalanceTest/"+strconv.Itoa(100), nil)
		} else {
			result = result+ "@" + q.RestfulInvoke(http.MethodGet, "cse://GoServerHoy/demo/LoadBalanceTest/"+strconv.Itoa(500), nil)
		}
	}
	return result
}

func highwayAutoCircuit()string{
	q.PrintlnGreen("Highway Invoke")
	result := "test"
	for sum := 0; sum < 15; sum++ {
		q.PrintlnRed("第", sum, "次")
		if sum > 4 {
			maptest := q.JsonMarshalExhange(q.HighwayInvoke(provider_scname,provider_schema,provider_operationId, &cse.CseInvokeRequest{Delaytime: 100}, &cse.CseInvokeReply{}))
			result = result+ "@" + maptest["instant_id"]
		} else {
			maptest := q.JsonMarshalExhange(q.HighwayInvoke(provider_scname,provider_schema,provider_operationId, &cse.CseInvokeRequest{Delaytime: 500}, &cse.CseInvokeReply{}))
			result = result+ "@" + maptest["instant_id"]
		}
	}
	return result
}