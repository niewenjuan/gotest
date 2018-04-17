package main

import (
	"github.com/ServiceComb/go-chassis"
	"github.com/ServiceComb/go-chassis/core/server"
	"gotest/basicUtil"
	qy "gotest/schema/protoc"
	q "gotest/schema/schemaForClient"
	"log"
	"net/http"
	"strconv"
	_ "github.com/ServiceComb/go-chassis/config-center"
)

const (
	serviceName = "GoServerHoy"
	schemaId    = "hellworld"
	loadBalance = "LoadBalanceTest"
	concurrent  = "ConcurrentTest"
	delaytime   = 0
)

func main() {
	chassis.RegisterSchema("rest", &q.RestFulConsumer{}, server.WithSchemaID("ClientRestful"))
	basicUtil.Isok(chassis.Init())
	chassis.Run()
	//rest_circuit()
}

func concurrent_test() {
	ch := make(chan interface{})
	sh := make(chan interface{})
	for i := 1; i < 16; i++ {
		go runtest_goroutine(i, ch, sh)
	}
	basicUtil.PrintlnGreen("Highway Invoke")
	for j := 0; j < 15; j++ {
		log.Println(<-ch)
	}
	//basicUtil.PrintlnCyan("rest Invoke")
	//for p :=0;p<15;p++ {
	//	log.Println(<-sh)
	//}
	close(ch)
	close(sh)
}

//（1）隔离超时验证，该接口延时为1s，保证该方法执行，可动态修改隔离配置，可同时验证服务级别和全局级别
func runTest() {
	for j := 0; j < 30; {
		//basicUtil.PrintlnGreen("Highway Invoke")
		//for sum := 0; sum < 2; sum++ {
		//	log.Println(basicUtil.HighwayInvoke(serviceName, schemaId, loadBalance, &qy.CseInvokeRequest{Delaytime: delaytime}, &qy.CseInvokeReply{}))
		//}
		basicUtil.PrintlnCyan("Rest Invoke")
		for i := 0; i < 1; i++ {
			log.Println(basicUtil.RestfulInvoke(http.MethodGet, "cse://GoServerHoy/demo/LoadBalanceTest/" + strconv.Itoa(delaytime), nil))
		}
	}

}

//容错验证
func runTest_failover() {
	for j := 0; j < 30; {
		basicUtil.PrintlnGreen("Highway Invoke")
		for sum := 0; sum < 1; sum++ {
			log.Println(basicUtil.HighwayInvoke(serviceName, schemaId, loadBalance, &qy.CseInvokeRequest{Delaytime: delaytime}, &qy.CseInvokeReply{}))
		}
		//basicUtil.PrintlnCyan("Rest Invoke")
		//for i := 0; i < 1; i++ {
		//	log.Println(basicUtil.RestfulInvoke(http.MethodGet, "cse://GoServerHoy/demo/LoadBalanceTest/" + strconv.Itoa(delaytime), nil))
		//}
	}
}

func runTest_router() {
	//basicUtil.PrintlnGreen("Highway Invoke")
	//for sum := 0; sum < 2;  {
	//	log.Println(basicUtil.HighwayInvoke(serviceName, schemaId, loadBalance, &qy.CseInvokeRequest{Delaytime: delaytime}, &qy.CseInvokeReply{}))
	//}
	basicUtil.PrintlnCyan("Rest Invoke")
	for i := 0; i < 1;  {
		log.Println(basicUtil.RestfulInvoke(http.MethodGet, "cse://GoServerHoy/demo/LoadBalanceTest/"+strconv.Itoa(delaytime), nil))
	}
}

func runtest_goroutine(flag int, ch, sh chan interface{}) {
	//highway
	reply := basicUtil.HighwayInvoke(serviceName, schemaId, concurrent, &qy.CseInvokeRequest{Delaytime: delaytime, Number: int64(flag)}, &qy.CseInvokeReply{})
	//replyone := basicUtil.RestfulInvoke(http.MethodGet, "cse://GoServerHoy/demo/ConcurrentTest/delaytime/"+strconv.Itoa(delaytime)+"/number/"+strconv.Itoa(flag), nil)
	ch <- reply
	//sh<- replyone
}

//rest熔断验证
func rest_circuit() {
	for j := 0; j < 1; j++ {
		basicUtil.PrintlnCyan("rest Invoke")
		for sum := 0; sum < 25; sum++ {
			basicUtil.PrintlnRed("第", sum, "次")
			if sum > 10 {
				log.Println(basicUtil.RestfulInvoke(http.MethodGet, "cse://GoServerHoy/demo/LoadBalanceTest/"+strconv.Itoa(100), nil))
			} else {
				log.Println(basicUtil.RestfulInvoke(http.MethodGet, "cse://GoServerHoy/demo/LoadBalanceTest/"+strconv.Itoa(500), nil))
			}
		}
		basicUtil.PrintlnGreen("等待恢复时间5秒")
		basicUtil.DelayTime(5000)

		for sum := 0; sum < 5; sum++ {
			basicUtil.PrintlnCyan("第", sum, "次")
			log.Println(basicUtil.RestfulInvoke(http.MethodGet, "cse://GoServerHoy/demo/LoadBalanceTest/"+strconv.Itoa(500), nil)) //恢复后第一次失败，再次熔断
		}
		basicUtil.PrintlnGreen("再次等待恢复时间5秒")
		basicUtil.DelayTime(5000)
		for sum := 0; sum < 5; sum++ {
			basicUtil.PrintlnRed("第", sum, "次")
			log.Println(basicUtil.RestfulInvoke(http.MethodGet, "cse://GoServerHoy/demo/LoadBalanceTest/"+strconv.Itoa(100), nil)) //恢复后第一次成功，提示退出熔断模式
		}
	}
}

//HIGHWAY 熔断验证
func highway_circuit() {
	for j := 0; j < 1; j++ {
		basicUtil.PrintlnGreen("Highway Invoke")
		for sum := 0; sum < 25; sum++ {
			basicUtil.PrintlnRed("第", sum, "次")
			if sum > 10 {
				log.Println(basicUtil.HighwayInvoke(serviceName, schemaId, loadBalance, &qy.CseInvokeRequest{Delaytime: 100}, &qy.CseInvokeReply{}))
			} else {
				log.Println(basicUtil.HighwayInvoke(serviceName, schemaId, loadBalance, &qy.CseInvokeRequest{Delaytime: 500}, &qy.CseInvokeReply{}))
			}
		}
		basicUtil.PrintlnGreen("等待恢复时间5秒")
		basicUtil.DelayTime(5000)

		for sum := 0; sum < 5; sum++ {
			basicUtil.PrintlnCyan("第", sum, "次")
			log.Println(basicUtil.HighwayInvoke(serviceName, schemaId, loadBalance, &qy.CseInvokeRequest{Delaytime: 500}, &qy.CseInvokeReply{})) //恢复后第一次失败，再次熔断
		}
		basicUtil.PrintlnGreen("再次等待恢复时间5秒")
		basicUtil.DelayTime(5000)
		for sum := 0; sum < 5; sum++ {
			basicUtil.PrintlnRed("第", sum, "次")
			log.Println(basicUtil.HighwayInvoke(serviceName, schemaId, loadBalance, &qy.CseInvokeRequest{Delaytime: 100}, &qy.CseInvokeReply{})) //恢复后第一次成功，提示退出熔断模式
		}
	}
}
