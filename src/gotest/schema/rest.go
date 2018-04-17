package highwayinvoke

import (
	"github.com/ServiceComb/go-chassis/server/restful"
	"net/http"
	p "gotest/basicUtil"
	"strconv"
	"github.com/ServiceComb/go-chassis/core/config"
	"log"
	"time"
)

type RestfulInvoke struct {
}

func (s *RestfulInvoke) SimpleInvoke(b *restful.Context) {
	name := b.ReadPathParameter("user")
	age := b.ReadPathParameter("age")
	b.Write([]byte("my name is :"+name+" and my age is :"+age))
}

func (s *RestfulInvoke) QpsTest(b *restful.Context) {
	min := strconv.Itoa(time.Now().Minute())
	sec := strconv.Itoa(time.Now().Second())
	log.Println(min+":"+sec)
	b.Write([]byte(min+":"+sec))
}

func (s *RestfulInvoke) LoadBalanceTest(b *restful.Context) {
	properity := config.MicroserviceDefinition.ServiceDescription.Properties
	temp := b.ReadPathParameter("delaytime")
	delaytime , _ :=strconv.Atoi(temp)
	p.DelayTime(int64(delaytime))
	res := properity["flag"]
	b.Write([]byte(res))
	return
}

func (s *RestfulInvoke) LoadBalanceTest_Weight(b *restful.Context) {
	properity := config.MicroserviceDefinition.ServiceDescription.Properties
	res := properity["flag"]
	temp := properity["delay"]
	delaytime , _ :=strconv.Atoi(temp)
	log.Println(delaytime)
	p.DelayTime(int64(delaytime))
	b.Write([]byte(res))
	return
}

func (s *RestfulInvoke) ConcurrentTest(b *restful.Context) {
	temp := b.ReadPathParameter("delaytime")
	number := b.ReadPathParameter("number")
	//log.Println(temp)
	//log.Println(number)
	delaytime , _ :=strconv.Atoi(temp)
	p.DelayTime(int64(delaytime))
	b.Write([]byte("instant_id:\""+Id+"\",Number:"+number))
	return
}

func (s *RestfulInvoke) URLPatterns() []restful.Route {
	return []restful.Route{
		{http.MethodGet, p.SimpleInvoke_url, p.SimpleInvoke_operationId},
		{http.MethodGet, p.Qps_url, p.Qps_operationId},
		{http.MethodGet, p.LoadBalanceTest_url, p.LoadBalanceTest_operationId},
		{http.MethodGet, p.LoadBalanceTest_Weight_url, p.LoadBalanceTest_Weight_operationId},
		{http.MethodGet, p.ConcurrentTest_url, p.ConcurrentTest_operationId},
	}
}

type RestfulSpecial struct {
}

func (s *RestfulSpecial) CircuitFailTest(b *restful.Context) {
	properity := config.MicroserviceDefinition.ServiceDescription.Properties
	res := properity["fail"]
	b.WriteHeader(http.StatusInternalServerError)
	b.Write([]byte(res))
}

func (s *RestfulSpecial) URLPatterns() []restful.Route {
	return []restful.Route{
		{http.MethodGet, p.CircuitFailInvoke_url, p.CircuitFailInvoke_operationId},

	}
}
