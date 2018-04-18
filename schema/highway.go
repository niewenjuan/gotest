package highwayinvoke

import (
	ctx "golang.org/x/net/context"
	cse "gotest/schema/protoc"
	basic "gotest/basicUtil"
	"github.com/ServiceComb/go-chassis/core/config"
	"strconv"
	"time"
)
var Id = basic.Generate_RandString(7)

type HighwayInvocation struct {
}

func (s *HighwayInvocation) SayHello(ctx ctx.Context,request *cse.CseInvokeRequest)(*cse.CseInvokeReply,error){
	return &cse.CseInvokeReply{ResMessage:"my name is :"+request.Name+" and my age is :"+request.Age},nil
}

func (s *HighwayInvocation) LoadBalanceTest(ctx ctx.Context,request *cse.CseInvokeRequest)(*cse.CseInvokeReply,error){
	properity := config.MicroserviceDefinition.ServiceDescription.Properties
	instId := properity["flag"]
	basic.DelayTime(request.Delaytime)
	return &cse.CseInvokeReply{InstantId:instId,Version:config.MicroserviceDefinition.ServiceDescription.Version},nil
}

func (s *HighwayInvocation) LoadBalanceTest_weight(ctx ctx.Context,request *cse.CseInvokeRequest)(*cse.CseInvokeReply,error){
	properity := config.MicroserviceDefinition.ServiceDescription.Properties
	instId := properity["flag"]
	temp := properity["delay"]
	delaytime , _ :=strconv.Atoi(temp)
	basic.DelayTime(int64(delaytime))
	return &cse.CseInvokeReply{InstantId:instId},nil
}

func (s *HighwayInvocation) ConcurrentTest(ctx ctx.Context,request *cse.CseInvokeRequest)(*cse.CseInvokeReply,error){
	instId := Id
	basic.DelayTime(request.Delaytime)
	return &cse.CseInvokeReply{Number:request.Number,InstantId:instId},nil
}

func (s *HighwayInvocation) QpsTest(ctx ctx.Context,request *cse.CseInvokeRequest)(*cse.CseInvokeReply,error) {
	min := strconv.Itoa(time.Now().Minute())
	sec := strconv.Itoa(time.Now().Second())
	instId := min+":"+sec
	return &cse.CseInvokeReply{InstantId:instId},nil
}
