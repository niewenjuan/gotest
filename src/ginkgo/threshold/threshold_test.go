package threshold_test

import (
	p "gotest/basicUtil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"io/ioutil"
	"net/http"
	"time"
	"strings"
)

var _ = Describe("GoSDKThreshold", func() {

	AfterEach(func() {
		By("TestcaseFail")
	})
	BeforeEach(func(){
		log.Println("\n\n\n")
		p.PrintlnRed("删除所有配置项")
		p.DeleteAllConfigForService(p.ServiceName)
	})
	BeforeSuite(func(){
		//p.PrintlnRed("用例执行前进行服务清理......")
		//p.DeleteAllConfigForService(p.ServiceName)
	})
	AfterSuite(func(){
		p.PrintlnRed("用例执行后进行服务清理......")
		p.DeleteAllConfigForService(p.ServiceName)
		time.Sleep(1*time.Second)
	})

	Context("Community",func(){
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.community.01",func(){
			p.PrintlnMegenta("REST:01 HTTP服务注册接入")
			sc :=p.ServiceCenter{}
			dat := p.JsonExchange(sc.QueryExistenceFromSC("default",p.ServiceName,"1.0.0"))
			res := p.JsonExchangeInter(sc.QueryInstanceIdFromSC(dat["serviceId"]))
			Expect(res["instances"]).NotTo(BeNil())
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.community.02",func(){
			p.PrintlnMegenta("REST:02 HTTP服务发现调用")
			url := p.ToRestInvoke("1000","temp")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err1 := ioutil.ReadAll(resp.Body)
			p.PrintlnCyan("body is:",string(body))
			Expect(err1).NotTo(HaveOccurred())
			flag := (string(body) == "a"||string(body) == "b")
			Expect(flag).To(Equal(true))
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.community.03",func(){
			p.PrintlnMegenta("REST:03 跨app服务发现调用")
			url := p.ToRestInvokeCrossApp(p.MidServerName,p.AppId_highwaySimpleInvoke,"1000","temp")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err1 := ioutil.ReadAll(resp.Body)
			p.PrintlnCyan(string(body))
			Expect(err1).NotTo(HaveOccurred())
			flag := (string(body) == "CrossApp")
			Expect(flag).To(Equal(true))
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.community.01",func(){
			p.PrintlnMegenta("HIGHWAY:01 RPC服务注册接入")
			sc :=p.ServiceCenter{}
			dat := p.JsonExchange(sc.QueryExistenceFromSC("default",p.ServerName,"1.0.1"))
			res := p.JsonExchangeInter(sc.QueryInstanceIdFromSC(dat["serviceId"]))
			Expect(res["instances"]).NotTo(BeNil())
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.community.02",func(){
			p.PrintlnMegenta("HIGHWAY:02 Highway服务发现调用")
			url := p.ToHighwayInvoke("1000")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err1 := ioutil.ReadAll(resp.Body)
			p.PrintlnCyan(string(body))
			Expect(err1).NotTo(HaveOccurred())
			dat := p.JsonExchange(body)
			Expect(dat["version"]).To(Equal("1.0.1"))
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.community.03",func(){
			p.PrintlnMegenta("HIGHWAY:03 跨app服务发现调用")
			url := p.ToHighwayInvokeCrossApp(p.MidServerName,p.AppId_highwaySimpleInvoke,p.SchemaId_highwaySimpleInvoke,p.OperationId_highwaySimpleInvoke,"1000")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err1 := ioutil.ReadAll(resp.Body)
			p.PrintlnCyan(string(body))
			Expect(err1).NotTo(HaveOccurred())
			dat := p.JsonExchange(body)
			Expect(dat["version"]).To(Equal("1.1.0"))
		})

	})
	Context("Loadbalance",func(){
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.loadbalance.01",func(){
			p.PrintlnMegenta("REST:01 轮询策略验证")
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_strategy] = p.Value_strategy_round
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.CycleInvokeUrl("rest","40","0")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			result := p.StrHandle(string(body),",",p.Value_strategy_round)

			//log.Println("删除配置项")
			//key := []string{p.Key_strategy}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(resp.StatusCode).To(Equal(200))
			Expect(result).To(Equal(p.Value_strategy_round))
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.loadbalance.02",func(){
			p.PrintlnMegenta("REST:02 权值策略验证")
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_strategy] = p.Value_strategy_weight
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.WeightInvokeUrl("rest","40")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			result := p.StrHandle(string(body),",",p.Value_strategy_weight)

			//log.Println("删除配置项")
			//key := []string{p.Key_strategy}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")

			Expect(result).To(Equal(p.Value_strategy_weight))
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.loadbalance.03",func(){
			p.PrintlnMegenta("REST:03 随机策略验证")
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_strategy] = p.Value_strategy_random
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.CycleInvokeUrl("rest","40","0")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			result := p.StrHandle(string(body),",",p.Value_strategy_random)

			//log.Println("删除配置项")
			//key := []string{p.Key_strategy}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")

			Expect(result).To(Equal(p.Value_strategy_random))
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.loadbalance.04",func(){
			p.PrintlnMegenta("REST:04 会话策略验证")
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_strategy] = p.Value_strategy_session
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(1*time.Second)

			log.Println("开始循环调用")
			url := p.SessionInvokeUrl("rest","40")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			//p.PrintlnCyan(string(body))
			result := p.StrHandle(string(body),",",p.Value_strategy_session)

			log.Println("删除配置项")
			key := []string{p.Key_strategy}
			ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			log.Println("删除完毕")

			Expect(result).To(Equal(p.Value_strategy_session))
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.loadbalance.05",func(){
			p.PrintlnMegenta("REST:05 负载策略优先级")
			//创建配置项
			log.Println("创建全局和指定服务级别策略")
			gt :=p.GinkgoTest{}
			data := make(map[string]interface{})
			data[p.Key_strategy] = p.Value_strategy_random
			data[gt.Get_Key_server_strategy()] = p.Value_strategy_round
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.CycleInvokeUrl("rest","40","0")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			result := p.StrHandle(string(body),",",p.Value_strategy_round)

			//log.Println("删除配置项")
			//key := []string{p.Key_strategy,gt.Get_Key_server_strategy()}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")

			Expect(result).To(Equal(p.Value_strategy_round))
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.loadbalance.01",func(){
			p.PrintlnMegenta("HIGHWAY:01 轮询策略验证")
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_strategy] = p.Value_strategy_round
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.CycleInvokeUrl("highway","40","0")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			result := p.StrHandle(string(body),",",p.Value_strategy_round)

			//log.Println("删除配置项")
			//key := []string{p.Key_strategy}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")

			Expect(result).To(Equal(p.Value_strategy_round))
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.loadbalance.02",func(){
			p.PrintlnMegenta("HIGHWAY:02 权值策略验证")
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_strategy] = p.Value_strategy_weight
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.WeightInvokeUrl("highway","40")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err1 := ioutil.ReadAll(resp.Body)
			Expect(err1).NotTo(HaveOccurred())
			result := p.StrHandle(string(body),",",p.Value_strategy_weight)

			//log.Println("删除配置项")
			//key := []string{p.Key_strategy}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")

			Expect(result).To(Equal(p.Value_strategy_weight))
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.loadbalance.03",func(){
			p.PrintlnMegenta("HIGHWAY:03 随机策略验证")
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_strategy] = p.Value_strategy_random
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			p.PrintlnYellow(string(ccbody))
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.CycleInvokeUrl("highway","40","0")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			result := p.StrHandle(string(body),",",p.Value_strategy_random)

			//log.Println("删除配置项")
			//key := []string{p.Key_strategy}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")

			Expect(result).To(Equal(p.Value_strategy_random))
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.loadbalance.04",func(){
			p.PrintlnMegenta("HIGHWAY:04 负载策略优先级")
			//创建配置项
			log.Println("创建全局和指定服务级别策略")
			gt :=p.GinkgoTest{}
			data := make(map[string]interface{})
			data[p.Key_strategy] = p.Value_strategy_random
			data[gt.Get_Key_server_strategy()] = p.Value_strategy_round
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.CycleInvokeUrl("highway","40","0")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			result := p.StrHandle(string(body),",",p.Value_strategy_round)

			//log.Println("删除配置项")
			//key := []string{p.Key_strategy,gt.Get_Key_server_strategy()}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")

			Expect(result).To(Equal(p.Value_strategy_round))
		})
	})
	Context("RouteRule",func(){
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.routerule.01",func(){
			gt :=p.GinkgoTest{}
			p.PrintlnMegenta("REST:01 路由权重策略多版本验证")
			var file []byte = []byte(`{"policyType": "RATE","ruleItems": [{"groupName": "s1","groupCondition": "version=1.0.1,1.0.0","policyCondition": "100"}]}`)
			cokey := gt.Get_Key_darklaunch_policy(p.ServerName)
			//创建配置项
			log.Println("创建全局和指定服务级别策略")
			data := make(map[string]interface{})
			data[p.Key_strategy] = p.Value_strategy_round
			data[cokey] = string(file)
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(1*time.Second)

			log.Println("开始循环调用")
			url := p.CycleInvokeUrl("rest","200","0")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			result1 := p.DarklanuchHandle(string(body),",")

			//log.Println("删除配置项")
			//key := []string{p.Key_strategy,cokey}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")

			Expect(result1).To(BeTrue())
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.routerule.02",func(){
			gt :=p.GinkgoTest{}
			p.PrintlnMegenta("REST:02 路由权重策略单版本验证")
			var file []byte = []byte(`{"policyType": "RATE","ruleItems": [{"groupName": "s1","groupCondition": "version=1.0.0","policyCondition": "50"}]}`)
			cokey := gt.Get_Key_darklaunch_policy(p.ServerName)
			//创建配置项
			log.Println("创建全局和指定服务级别策略")
			data := make(map[string]interface{})
			data[p.Key_strategy] = p.Value_strategy_round
			data[cokey] = string(file)
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.CycleInvokeUrl("rest","200","0")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			result1 := p.DarklanuchHandle(string(body),",")

			//log.Println("删除配置项")
			//key := []string{p.Key_strategy,cokey}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")

			Expect(result1).To(BeTrue())
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.routerule.01",func(){
			gt :=p.GinkgoTest{}
			p.PrintlnMegenta("HIGHWAY:01 路由权重策略多版本验证")
			var file []byte = []byte(`{"policyType": "RATE","ruleItems": [{"groupName": "s1","groupCondition": "version=1.0.1,1.0.0","policyCondition": "100"}]}`)
			cokey := gt.Get_Key_darklaunch_policy(p.ServerName)
			//创建配置项
			log.Println("创建全局和指定服务级别策略")
			data := make(map[string]interface{})
			data[p.Key_strategy] = p.Value_strategy_round
			data[cokey] = string(file)
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(1*time.Second)

			log.Println("开始循环调用")
			url := p.CycleInvokeUrl("highway","200","0")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			result1 := p.DarklanuchHandle(string(body),",")
			//
			//log.Println("删除配置项")
			//key := []string{p.Key_strategy,cokey}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")

			Expect(result1).To(BeTrue())

		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.routerule.02",func(){
			gt :=p.GinkgoTest{}
			p.PrintlnMegenta("HIGHWAY:02 路由权重策略单版本验证")
			var file []byte = []byte(`{"policyType": "RATE","ruleItems": [{"groupName": "s1","groupCondition": "version=1.0.0","policyCondition": "50"}]}`)
			cokey := gt.Get_Key_darklaunch_policy(p.ServerName)
			//创建配置项
			log.Println("创建全局和指定服务级别策略")
			data := make(map[string]interface{})
			data[p.Key_strategy] = p.Value_strategy_round
			data[cokey] = string(file)
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(1*time.Second)

			log.Println("开始循环调用")
			url := p.CycleInvokeUrl("highway","200","0")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			result1 := p.DarklanuchHandle(string(body),",")

			//log.Println("删除配置项")
			//key := []string{p.Key_strategy,cokey}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")

			Expect(result1).To(BeTrue())
		})
	})
	Context("QpsControl",func(){
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.qps.01",func(){
			p.PrintlnMegenta("REST:01 流控Consumer端global不生效验证")
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_qps_global_enable] = true
			data[p.Key_qps_global_limit] = "1"
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(1*time.Second)

			log.Println("开始循环调用")
			url := p.CycleInvokeQpsUrl("rest","40","0")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			result := p.NoQpsHandle(string(body),",")

			//log.Println("删除配置项")
			//key := []string{p.Key_qps_global_enable,p.Key_qps_global_limit}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(result).To(Equal(0))

		})
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.qps.02",func(){
			p.PrintlnMegenta("REST:02 流控Consumer端验证")
			gc := p.GinkgoTest{}
			limit := gc.Get_Key_qps_service_limit(p.ServerName)
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_qps_global_enable] = true
			data[limit] = "1"
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(1*time.Second)

			log.Println("开始循环调用")
			url := p.CycleInvokeQpsUrl("rest","40","0")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			result := p.QpsHandle(string(body),",")

			//log.Println("删除配置项")
			//key := []string{p.Key_qps_global_enable,limit}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(result).To(Equal(0))

		})
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.qps.03",func(){
			p.PrintlnMegenta("REST:03 流控Consumer优先级验证")
			gc := p.GinkgoTest{}
			limit := gc.Get_Key_qps_service_limit(p.ServerName)
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_qps_global_enable] = true
			data[p.Key_qps_global_limit] = "2"
			data[limit] = "1"
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(1*time.Second)

			log.Println("开始循环调用")
			url := p.CycleInvokeQpsUrl("rest","30","0")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			result := p.QpsHandle(string(body),",")

			//log.Println("删除配置项")
			//key := []string{p.Key_qps_global_enable,p.Key_qps_global_limit,limit}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(result).To(Equal(0))

		})
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.qps.04",func(){
			p.PrintlnMegenta("REST:04 流控provider端global验证")
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_qps_provider_global_enable] = true
			data[p.Key_qps_provider_global_limit] = "1"
			ccbody := p.CreateConfigToCC(p.MidServerName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(1*time.Second)

			log.Println("开始循环调用")
			url := p.CycleInvokeCrossAppQpsUrl("rest","40")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			result := p.QpsHandle(string(body),",")

			log.Println("删除配置项")
			key := []string{p.Key_qps_provider_global_enable,p.Key_qps_provider_global_limit}
			ccbody = p.DeleteConfigFromCC(p.MidServerName,key)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			log.Println("删除完毕")
			Expect(result).To(Equal(0))

		})
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.qps.05",func(){
			p.PrintlnMegenta("REST:05 流控provider端服务级别验证")
			gt :=p.GinkgoTest{}
			limit := gt.Get_Key_qps_provider_limit(p.MidServerName)
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_qps_provider_global_enable] = true
			data[limit] = "1"
			ccbody := p.CreateConfigToCC(p.MidServerName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(1*time.Second)

			log.Println("开始循环调用")
			url := p.CycleInvokeCrossAppQpsUrl("rest","40")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			result := p.QpsHandle(string(body),",")

			log.Println("删除配置项")
			key := []string{p.Key_qps_provider_global_enable,limit}
			ccbody = p.DeleteConfigFromCC(p.MidServerName,key)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			log.Println("删除完毕")
			Expect(result).To(Equal(0))

		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.qps.01",func(){
			p.PrintlnMegenta("HIGHWAY:01 流控流控Consumer端global不生效验证验证")
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_qps_global_enable] = true
			data[p.Key_qps_global_limit] = "1"
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(1*time.Second)

			log.Println("开始循环调用")
			url := p.CycleInvokeQpsUrl("highway","40","0")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			result := p.NoQpsHandle(string(body),",")

			//log.Println("删除配置项")
			//key := []string{p.Key_qps_global_enable,p.Key_qps_global_limit}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(result).To(Equal(0))

		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.qps.02",func(){
			p.PrintlnMegenta("HIGHWAY:02 流控Consumer端验证")
			gc := p.GinkgoTest{}
			limit := gc.Get_Key_qps_service_limit(p.ServerName)
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_qps_global_enable] = true
			data[limit] = "1"
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(1*time.Second)

			log.Println("开始循环调用")
			url := p.CycleInvokeQpsUrl("highway","40","0")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			result := p.QpsHandle(string(body),",")

			//log.Println("删除配置项")
			//key := []string{p.Key_qps_global_enable,limit}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(result).To(Equal(0))

		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.qps.03",func(){
			p.PrintlnMegenta("HIGHWAY:03 流控Consumer优先级验证")
			gc := p.GinkgoTest{}
			limit := gc.Get_Key_qps_service_limit(p.ServerName)
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_qps_global_enable] = true
			data[p.Key_qps_global_limit] = "2"
			data[limit] = "1"
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(1*time.Second)

			log.Println("开始循环调用")
			url := p.CycleInvokeQpsUrl("highway","20","0")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			result := p.QpsHandle(string(body),",")

			//log.Println("删除配置项")
			//key := []string{p.Key_qps_global_enable,p.Key_qps_global_limit,limit}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(result).To(Equal(0))

		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.qps.04",func(){
			p.PrintlnMegenta("HIGHWAY:04 流控provider验证")
			//gc := p.GinkgoTest{}
			//limit := gc.Get_Key_qps_service_limit(p.ServerName)
			//创建配置项
			data := make(map[string]interface{})
			data[p.Key_qps_provider_global_enable] = true
			data[p.Key_qps_provider_global_limit] = "1"
			ccbody := p.CreateConfigToCC(p.MidServerName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(1*time.Second)

			log.Println("开始循环调用")
			url := p.CycleInvokeCrossAppQpsUrl("highway","40")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			result := p.QpsHandle(string(body),",")

			log.Println("删除配置项")
			key := []string{p.Key_qps_provider_global_enable,p.Key_qps_provider_global_limit}
			ccbody = p.DeleteConfigFromCC(p.MidServerName,key)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			log.Println("删除完毕")
			Expect(result).To(Equal(0))

		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.qps.05",func(){
			p.PrintlnMegenta("HIGHWAY:05 流控provider服务级别验证")
			gc := p.GinkgoTest{}
			limit := gc.Get_Key_qps_provider_limit(p.MidServerName)
			//创建配置项
			data := make(map[string]interface{})
			data[p.Key_qps_provider_global_enable] = true
			data[limit] = "1"
			ccbody := p.CreateConfigToCC(p.MidServerName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(1*time.Second)

			log.Println("开始循环调用")
			url := p.CycleInvokeCrossAppQpsUrl("highway","40")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			result := p.QpsHandle(string(body),",")

			log.Println("删除配置项")
			key := []string{p.Key_qps_provider_global_enable,limit}
			ccbody = p.DeleteConfigFromCC(p.MidServerName,key)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			log.Println("删除完毕")
			Expect(result).To(Equal(0))

		})
	})
	Context("TimeOut",func(){
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.timeout.01",func(){
			p.PrintlnMegenta("REST:01 隔离超时全局验证")
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = true
			data[p.Key_isolation_global_timeout_able] = true
			data[p.Key_isolation_global_timeoutInMilliseconds] = 100
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始简单调用")
			url := p.ToRestInvoke("1000","rest")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			log.Println("body is",string(body))
			log.Println("statuscode is",resp.StatusCode)
			arr := p.Read3(p.LogPath)
			clog := p.LogReadInLinesOperation(arr,1,p.FallbackTimeOut)

			//log.Println("删除配置项")
			//key := []string{p.Key_isolation_global_timeout_able,p.Key_isolation_global_timeoutInMilliseconds}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(clog).To(BeTrue())
			Expect(resp.StatusCode).To(Equal(408))

		})
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.timeout.02",func(){
			p.PrintlnMegenta("REST:02 隔离超时服务级别验证")
			gt :=p.GinkgoTest{}
			enable := gt.Get_Key_isolation_service_timeout_able(p.ServerName)
			millseconds := gt.Get_Key_isolation_service_timeoutInMilliseconds(p.ServerName)
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = true
			data[enable] = true
			data[millseconds] = 100
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始简单调用")
			url := p.ToRestInvoke("1000","rest")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			log.Println("body is",string(body))
			log.Println("statuscode is",resp.StatusCode)
			arr := p.Read3(p.LogPath)
			clog := p.LogReadInLinesOperation(arr,1,p.FallbackTimeOut)

			//log.Println("删除配置项")
			//key := []string{enable,millseconds}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(clog).To(BeTrue())
			Expect(resp.StatusCode).To(Equal(408))

		})
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.timeout.03",func(){
			p.PrintlnMegenta("REST:03 隔离超时优先级验证")
			gt :=p.GinkgoTest{}
			enable := gt.Get_Key_isolation_service_timeout_able(p.ServerName)
			millseconds := gt.Get_Key_isolation_service_timeoutInMilliseconds(p.ServerName)
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[enable] = true
			data[p.Key_fallback_consumer_enable] = true
			data[millseconds] = 100
			data[p.Key_isolation_global_timeout_able] = false
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.ToRestInvoke("1000","rest")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			log.Println("body is",string(body))
			log.Println("statuscode is",resp.StatusCode)
			arr := p.Read3(p.LogPath)
			clog := p.LogReadInLinesOperation(arr,1,p.FallbackTimeOut)

			//log.Println("删除配置项")
			//key := []string{enable,millseconds,p.Key_isolation_global_timeout_able}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(clog).To(BeTrue())
			Expect(resp.StatusCode).To(Equal(408))

		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.timeout.01",func(){
			p.PrintlnMegenta("HIGHWAY:01 隔离超时全局验证")
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = false
			data[p.Key_isolation_global_timeout_able] = true
			data[p.Key_isolation_global_timeoutInMilliseconds] = 100
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.ToHighwayInvoke("1000")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			log.Println(string(body))

			//log.Println("删除配置项")
			//key := []string{p.Key_isolation_global_timeout_able,p.Key_isolation_global_timeoutInMilliseconds}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(string(body)).To(Equal("{\"Message\":\"timeout\"}"))

		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.timeout.02",func(){
			p.PrintlnMegenta("HIGHWAY:02 隔离超时服务级别验证")
			gt :=p.GinkgoTest{}
			enable := gt.Get_Key_isolation_service_timeout_able(p.ServerName)
			millseconds := gt.Get_Key_isolation_service_timeoutInMilliseconds(p.ServerName)
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[enable] = true
			data[millseconds] = 100
			data[p.Key_fallback_consumer_enable] = false
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.ToHighwayInvoke("1000")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			log.Println(string(body))

			//log.Println("删除配置项")
			//key := []string{enable,millseconds}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(string(body)).To(Equal("{\"Message\":\"timeout\"}"))
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.timeout.03",func(){
			p.PrintlnMegenta("HIGHWAY:03 隔离超时优先级验证")
			gt :=p.GinkgoTest{}
			enable := gt.Get_Key_isolation_service_timeout_able(p.ServerName)
			millseconds := gt.Get_Key_isolation_service_timeoutInMilliseconds(p.ServerName)
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = false
			data[p.Key_isolation_global_timeout_able] = false
			data[enable] = true
			data[millseconds] = 100
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.ToHighwayInvoke("1000")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			log.Println(string(body))

			//log.Println("删除配置项")
			//key := []string{enable,millseconds,p.Key_isolation_global_timeout_able}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(string(body)).To(Equal("{\"Message\":\"timeout\"}"))

		})
	})
	Context("Current",func(){
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.current.01",func(){
			p.PrintlnMegenta("REST:01 并发隔离全局验证")
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = true
			data[p.Key_isolation_global_maxConcurrentRequests] = 2
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(1*time.Second)

			log.Println("开始并发调用")
			flag := p.CurrentRequest("rest",10,2)
			arr := p.Read3(p.LogPath)
			clog := p.LogReadInLinesOperationMaxCurrent(arr,10,2,p.MaxCurrentError)

			//log.Println("删除配置项")
			//key := []string{p.Key_isolation_global_maxConcurrentRequests}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(flag).To(Equal(true))
			Expect(clog).To(BeTrue())

		})
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.current.02",func(){
			p.PrintlnMegenta("REST:02 并发隔离服务级别验证")
			gt := p.GinkgoTest{}
			maxCurrent := gt.Get_Key_isolation_service_maxConcurrentRequests(p.ServerName)
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = true
			data[maxCurrent] = 2
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始并发调用")
			flag := p.CurrentRequest("rest",10,2)
			arr := p.Read3(p.LogPath)
			clog := p.LogReadInLinesOperationMaxCurrent(arr,10,2,p.MaxCurrentError)

			//log.Println("删除配置项")
			//key := []string{maxCurrent}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(flag).To(Equal(true))
			Expect(clog).To(BeTrue())

		})
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.current.03",func(){
			p.PrintlnMegenta("REST:03 并发隔离优先级验证")
			gt := p.GinkgoTest{}
			maxCurrent := gt.Get_Key_isolation_service_maxConcurrentRequests(p.ServerName)
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = true
			data[p.Key_isolation_global_maxConcurrentRequests] = 5
			data[maxCurrent] = 1
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始并发调用")
			flag := p.CurrentRequest("rest",10,2)
			arr := p.Read3(p.LogPath)
			clog := p.LogReadInLinesOperationMaxCurrent(arr,10,2,p.MaxCurrentError)

			log.Println("删除配置项")
			key := []string{p.Key_isolation_global_maxConcurrentRequests,maxCurrent}
			ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			log.Println("删除完毕")
			Expect(flag).To(Equal(true))
			Expect(clog).To(BeTrue())

		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.current.01",func(){
			p.PrintlnMegenta("HIGHWAY:01 隔离并发全局验证")
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = false
			data[p.Key_isolation_global_maxConcurrentRequests] = 2
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始并发调用")
			flag := p.CurrentRequest("highway",10,2)

			//log.Println("删除配置项")
			//key := []string{p.Key_isolation_global_maxConcurrentRequests}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(flag).To(Equal(true))

		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.current.02",func(){
			p.PrintlnMegenta("HIGHWAY:02 隔离并发服务级别验证")
			gt := p.GinkgoTest{}
			maxCurrent := gt.Get_Key_isolation_service_maxConcurrentRequests(p.ServerName)
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = false
			data[maxCurrent] = 2
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始并发调用")
			flag := p.CurrentRequest("highway",10,2)

			//log.Println("删除配置项")
			//key := []string{maxCurrent}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(flag).To(Equal(true))

		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.current.03",func(){
			p.PrintlnMegenta("HIGHWAY:03 隔离并发优先级验证")
			gt := p.GinkgoTest{}
			maxCurrent := gt.Get_Key_isolation_service_maxConcurrentRequests(p.ServerName)
			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = false
			data[p.Key_isolation_global_maxConcurrentRequests] = 5
			data[maxCurrent] = 2
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始并发调用")
			flag := p.CurrentRequest("highway",10,2)

			//log.Println("删除配置项")
			//key := []string{p.Key_isolation_global_maxConcurrentRequests,maxCurrent}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(flag).To(Equal(true))

		})
	})
	Context("ManualCircuit",func(){
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.manualcircuit.01",func(){
			p.PrintlnMegenta("REST:01 手动熔断全局验证")

			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = false
			data[p.Key_circuit_consumer_enable] = true
			data[p.Key_circuit_manaul_forceOpen] = true
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(1*time.Second)

			log.Println("开始循环调用")
			url := p.ToRestInvoke("1000","rest")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			log.Println(string(body))

			//log.Println("删除配置项")
			//key := []string{p.Key_circuit_consumer_enable,p.Key_circuit_manaul_forceOpen}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(string(body)).To(Equal(""))
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.manualcircuit.02",func(){
			p.PrintlnMegenta("REST:21 手动熔断服务级别验证")
			gc := p.GinkgoTest{}
			enable := gc.Get_Key_circuit_service_consumer_enable(p.ServerName)
			forceOpen := gc.Get_Key_circuit_service_manaul_forceOpen(p.ServerName)

			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = false
			data[enable] = true
			data[forceOpen] = true
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(1*time.Second)

			log.Println("开始循环调用")
			url := p.ToRestInvoke("1000","rest")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			log.Println(string(body))

			//log.Println("删除配置项")
			//key := []string{enable,forceOpen}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(string(body)).To(Equal(""))
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.manualcircuit.03",func(){
			p.PrintlnMegenta("REST:22 手动熔断优先级验证")
			gc := p.GinkgoTest{}
			enable := gc.Get_Key_circuit_service_consumer_enable(p.ServerName)
			forceOpen := gc.Get_Key_circuit_service_manaul_forceOpen(p.ServerName)

			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = false
			data[p.Key_circuit_consumer_enable] = false
			data[p.Key_circuit_manaul_forceOpen] = false
			data[enable] = true
			data[forceOpen] = true
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.ToRestInvoke("1000","rest")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			log.Println(string(body))

			//log.Println("删除配置项")
			//key := []string{p.Key_circuit_consumer_enable,p.Key_circuit_manaul_forceOpen,enable,forceOpen}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(string(body)).To(Equal(""))
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.manualcircuit.01",func(){
			p.PrintlnMegenta("HIGHWAY:20 手动熔断全局验证")

			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = false
			data[p.Key_circuit_consumer_enable] = true
			data[p.Key_circuit_manaul_forceOpen] = true
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.ToHighwayInvoke("1000")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			log.Println(string(body))

			//log.Println("删除配置项")
			//key := []string{p.Key_circuit_consumer_enable,p.Key_circuit_manaul_forceOpen}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(string(body)).To(Equal("{\"Message\":\"circuit open\"}"))
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.manualcircuit.02",func(){
			p.PrintlnMegenta("HIGHWAY:21 手动熔断服务级别验证")
			gt := p.GinkgoTest{}
			enable := gt.Get_Key_circuit_service_consumer_enable(p.ServerName)
			forceOpen := gt.Get_Key_circuit_service_manaul_forceOpen(p.ServerName)

			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = false
			data[enable] = true
			data[forceOpen] = true
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.ToHighwayInvoke("1000")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			log.Println(string(body))

			//log.Println("删除配置项")
			//key := []string{enable,forceOpen}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(string(body)).To(Equal("{\"Message\":\"circuit open\"}"))
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.manualcircuit.03",func(){
			p.PrintlnMegenta("HIGHWAY:22 手动熔断优先级验证")
			gt := p.GinkgoTest{}
			enable := gt.Get_Key_circuit_service_consumer_enable(p.ServerName)
			forceOpen := gt.Get_Key_circuit_service_manaul_forceOpen(p.ServerName)

			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = false
			data[p.Key_circuit_consumer_enable] = false
			data[p.Key_circuit_manaul_forceOpen] = false
			data[enable] = true
			data[forceOpen] = true
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.ToHighwayInvoke("1000")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			log.Println(string(body))

			//log.Println("删除配置项")
			//key := []string{p.Key_circuit_consumer_enable,p.Key_circuit_manaul_forceOpen,enable,forceOpen}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(string(body)).To(Equal("{\"Message\":\"circuit open\"}"))
		})
	})
	Context("AutoCircuit",func(){
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.autocircuit.01",func(){
			p.PrintlnMegenta("REST:01 自动熔断全局验证")

			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = false
			data[p.Key_isolation_global_timeout_able] = true
			data[p.Key_isolation_global_timeoutInMilliseconds] = 300
			data[p.Key_circuit_consumer_enable] = true
			data[p.Key_circuit_consumer_requestVolumeThreshold] = 10
			data[p.Key_circuit_consumer_errorThresholdPercentage] = 50
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.ToInvokeAutoCircuit("rest")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body,_ := ioutil.ReadAll(resp.Body)
			istrue := p.CircuitHandle(string(body),"@",10)

			//log.Println("删除配置项")
			//key := []string{p.Key_isolation_global_timeout_able,p.Key_isolation_global_timeoutInMilliseconds,p.Key_circuit_consumer_enable,p.Key_circuit_consumer_requestVolumeThreshold,p.Key_circuit_consumer_errorThresholdPercentage}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(istrue).To(BeTrue())
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.autocircuit.02",func(){
			p.PrintlnMegenta("REST:02 自动熔断服务级别验证")
			gt :=p.GinkgoTest{}
			circuitEnable := gt.Get_Key_circuit_service_consumer_enable(p.ServerName)
			requestVolume := gt.Get_Key_circuit_service_consumer_requestVolumeThreshold(p.ServerName)
			errorThresholdPercentage := gt.Get_Key_circuit_service_consumer_errorThresholdPercentage(p.ServerName)

			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = false
			data[p.Key_isolation_global_timeout_able] = true
			data[p.Key_isolation_global_timeoutInMilliseconds] = 300
			data[circuitEnable] = true
			data[requestVolume] = 10
			data[errorThresholdPercentage] = 50
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.ToInvokeAutoCircuit("rest")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body,_ := ioutil.ReadAll(resp.Body)
			istrue := p.CircuitHandle(string(body),"@",10)

			//log.Println("删除配置项")
			//key := []string{p.Key_isolation_global_timeout_able,p.Key_isolation_global_timeoutInMilliseconds,circuitEnable,requestVolume,errorThresholdPercentage}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(istrue).To(BeTrue())
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.autocircuit.03",func(){
			p.PrintlnMegenta("REST:03 自动熔断优先级验证")
			gt :=p.GinkgoTest{}
			circuitEnable := gt.Get_Key_circuit_service_consumer_enable(p.ServerName)
			requestVolume := gt.Get_Key_circuit_service_consumer_requestVolumeThreshold(p.ServerName)
			errorThresholdPercentage := gt.Get_Key_circuit_service_consumer_errorThresholdPercentage(p.ServerName)

			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = false
			data[p.Key_isolation_global_timeout_able] = true
			data[p.Key_isolation_global_timeoutInMilliseconds] = 300
			data[p.Key_circuit_consumer_enable] = false
			data[p.Key_circuit_consumer_requestVolumeThreshold] = 10
			data[p.Key_circuit_consumer_errorThresholdPercentage] = 100
			data[circuitEnable] = true
			data[requestVolume] = 10
			data[errorThresholdPercentage] = 50
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.ToInvokeAutoCircuit("rest")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body,_ := ioutil.ReadAll(resp.Body)
			istrue := p.CircuitHandle(string(body),"@",10)

			//log.Println("删除配置项")
			//key := []string{p.Key_isolation_global_timeout_able,p.Key_isolation_global_timeoutInMilliseconds,p.Key_circuit_consumer_enable,p.Key_circuit_consumer_requestVolumeThreshold,p.Key_circuit_consumer_errorThresholdPercentage,circuitEnable,requestVolume,errorThresholdPercentage}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(istrue).To(BeTrue())
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.autocircuit.01",func(){
			p.PrintlnMegenta("HIGHWAY:01 自动熔断全局验证")

			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = false
			data[p.Key_isolation_global_timeout_able] = true
			data[p.Key_isolation_global_timeoutInMilliseconds] = 300
			data[p.Key_circuit_consumer_enable] = true
			data[p.Key_circuit_consumer_requestVolumeThreshold] = 10
			data[p.Key_circuit_consumer_errorThresholdPercentage] = 50
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.ToInvokeAutoCircuit("highway")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body,_ := ioutil.ReadAll(resp.Body)
			istrue := p.CircuitHandle(string(body),"@",10)

			//log.Println("删除配置项")
			//key := []string{p.Key_isolation_global_timeout_able,p.Key_isolation_global_timeoutInMilliseconds,p.Key_circuit_consumer_enable,p.Key_circuit_consumer_requestVolumeThreshold,p.Key_circuit_consumer_errorThresholdPercentage}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(istrue).To(BeTrue())
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.autocircuit.02",func(){
			p.PrintlnMegenta("HIGHWAY:02 自动熔断服务级别验证")
			gt :=p.GinkgoTest{}
			circuitEnable := gt.Get_Key_circuit_service_consumer_enable(p.ServerName)
			requestVolume := gt.Get_Key_circuit_service_consumer_requestVolumeThreshold(p.ServerName)
			errorThresholdPercentage := gt.Get_Key_circuit_service_consumer_errorThresholdPercentage(p.ServerName)

			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = false
			data[p.Key_isolation_global_timeout_able] = true
			data[p.Key_isolation_global_timeoutInMilliseconds] = 300
			data[circuitEnable] = true
			data[requestVolume] = 10
			data[errorThresholdPercentage] = 50
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.ToInvokeAutoCircuit("highway")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body,_ := ioutil.ReadAll(resp.Body)
			istrue := p.CircuitHandle(string(body),"@",10)

			//log.Println("删除配置项")
			//key := []string{p.Key_isolation_global_timeout_able,p.Key_isolation_global_timeoutInMilliseconds,circuitEnable,requestVolume,errorThresholdPercentage}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(istrue).To(BeTrue())
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.autocircuit.03",func(){
			p.PrintlnMegenta("HIGHWAY:03 自动熔断优先级验证")
			gt :=p.GinkgoTest{}
			circuitEnable := gt.Get_Key_circuit_service_consumer_enable(p.ServerName)
			requestVolume := gt.Get_Key_circuit_service_consumer_requestVolumeThreshold(p.ServerName)
			errorThresholdPercentage := gt.Get_Key_circuit_service_consumer_errorThresholdPercentage(p.ServerName)

			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = false
			data[p.Key_isolation_global_timeout_able] = true
			data[p.Key_isolation_global_timeoutInMilliseconds] = 300
			data[p.Key_circuit_consumer_enable] = false
			data[p.Key_circuit_consumer_requestVolumeThreshold] = 10
			data[p.Key_circuit_consumer_errorThresholdPercentage] = 100
			data[circuitEnable] = true
			data[requestVolume] = 10
			data[errorThresholdPercentage] = 50
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.ToInvokeAutoCircuit("highway")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body,_ := ioutil.ReadAll(resp.Body)
			istrue := p.CircuitHandle(string(body),"@",10)

			//log.Println("删除配置项")
			//key := []string{p.Key_isolation_global_timeout_able,p.Key_isolation_global_timeoutInMilliseconds,p.Key_circuit_consumer_enable,p.Key_circuit_consumer_requestVolumeThreshold,p.Key_circuit_consumer_errorThresholdPercentage,circuitEnable,requestVolume,errorThresholdPercentage}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(istrue).To(BeTrue())
		})
	})
	Context("Fail",func(){
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.fail.01",func(){
			p.PrintlnMegenta("REST:01 容错尝试同实例验证")

			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = false
			data[p.Key_retry_enable] = true
			data[p.Key_retry_onnext] = 0
			data[p.Key_retry_onsame] = 2
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始调用失败实例")
			url := p.ToInvokeFailInstance(p.REST,p.ServerName)
			_, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			arr := p.Read3(p.LogPath)
			clog := p.LogReadInLines(arr,3)

			//log.Println("删除配置项")
			//key := []string{p.Key_retry_enable,p.Key_retry_onnext,p.Key_retry_onsame}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(clog).To(BeTrue())
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.fail.02",func(){
			p.PrintlnMegenta("REST:02 容错尝试新实例验证")

			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = false
			data[p.Key_retry_enable] = true
			data[p.Key_retry_onnext] = 4
			data[p.Key_retry_onsame] = 0
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始调用失败实例")
			url := p.ToInvokeFailInstance(p.REST,p.ServerName)
			_, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			arr := p.Read3(p.LogPath)
			clog := p.LogReadInLinesOnnext(arr,5)

			//log.Println("删除配置项")
			//key := []string{p.Key_retry_enable,p.Key_retry_onnext,p.Key_retry_onsame}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(clog).To(BeTrue())
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.fail.01",func(){
			p.PrintlnMegenta("HIGHWAY:03 容错尝试同实例验证")

			//创建配置项
			log.Println("创建配置项")
			sc := p.ServiceCenter{}
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = false
			data[p.Key_retry_enable] = true
			data[p.Key_retry_onnext] = 0
			data[p.Key_retry_onsame] = 2
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))

			log.Println("创建空实例")
			serviceID := p.CreateService("FAILTEST","default","1.0.1")
			log.Println(serviceID)
			resp1 := p.CreateInstances(serviceID,"szx1234","highway:127.0.0.1:8080")
			log.Println(resp1)
			Expect(strings.Contains(resp1,"instanceId")).To(BeTrue())
			time.Sleep(4*time.Second)

			log.Println("开始调用")
			url := p.ToInvokeFailInstance(p.HIGHWAY,"FAILTEST")
			_, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			arr := p.Read3(p.LogPath)
			clog := p.LogReadInLinesOnnextByHighway(arr,9)

			log.Println("强制删除服务")
			log.Println(string(sc.DeleteServiceFromSC(serviceID)))
			//log.Println("删除配置项")
			//key := []string{p.Key_retry_enable,p.Key_retry_onnext,p.Key_retry_onsame}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(clog).To(Equal(3))
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.fail.02",func(){
			p.PrintlnMegenta("HIGHWAY:02 容错尝试新实例验证")

			//创建配置项
			log.Println("创建配置项")
			sc := p.ServiceCenter{}
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = false
			data[p.Key_retry_enable] = true
			data[p.Key_retry_onnext] = 3
			data[p.Key_retry_onsame] = 0
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))

			log.Println("创建空实例")
			serviceID := p.CreateService("FAILTESTA","default","1.0.1")
			log.Println(serviceID)
			resp1 := p.CreateInstances(serviceID,"szx1234","highway:127.0.0.1:8080")
			log.Println(resp1)
			Expect(strings.Contains(resp1,"instanceId")).To(BeTrue())
			time.Sleep(4*time.Second)

			log.Println("开始调用")
			url := p.ToInvokeFailInstance(p.HIGHWAY,"FAILTESTA")
			_, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			arr := p.Read3(p.LogPath)
			clog := p.LogReadInLinesOnnextByHighway(arr,9)

			log.Println("强制删除服务")
			log.Println(string(sc.DeleteServiceFromSC(serviceID)))
			//log.Println("删除配置项")
			//key := []string{p.Key_retry_enable,p.Key_retry_onnext,p.Key_retry_onsame}
			//ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			//Expect(string(ccbody)).To(Equal(p.Result_success))
			//log.Println("删除完毕")
			Expect(clog).To(Equal(3))
		})
	})
	Context("Fallback",func(){
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.fallback.01",func(){
			p.PrintlnMegenta("REST:01 容错降级验证")

			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = true
			data[p.Key_consumer_fallbackpolicy] = p.Value_fallbackpolicy_exception
			data[p.Key_isolation_global_timeout_able] = true
			data[p.Key_isolation_global_timeoutInMilliseconds] = 100
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.ToRestInvoke("1000","rest")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			log.Println(string(body))
			log.Println(resp.StatusCode)

			log.Println("删除配置项")
			key := []string{p.Key_fallback_consumer_enable,p.Key_consumer_fallbackpolicy,p.Key_isolation_global_timeout_able,p.Key_isolation_global_timeoutInMilliseconds}
			ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			log.Println("删除完毕")
			Expect(string(body)).To(Equal(""))
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.fallback.02",func(){
			p.PrintlnMegenta("REST:02 容错降级服务级别验证")
			gt := p.GinkgoTest{}
			fallbackEnable := gt.Get_Key_fallback_consumer_service_enable(p.ServerName)
			policy := gt.Get_Key_consumer_service_fallbackpolicy(p.ServerName)

			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[fallbackEnable] = true
			data[policy] = p.Value_fallbackpolicy_exception
			data[p.Key_isolation_global_timeout_able] = true
			data[p.Key_isolation_global_timeoutInMilliseconds] = 100
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.ToRestInvoke("1000","rest")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			log.Println(string(body))
			log.Println(resp.StatusCode)

			log.Println("删除配置项")
			key := []string{fallbackEnable,policy,p.Key_isolation_global_timeout_able,p.Key_isolation_global_timeoutInMilliseconds}
			ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			log.Println("删除完毕")
			Expect(string(body)).To(Equal(""))
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.REST.fallback.03",func(){
			p.PrintlnMegenta("REST:03 容错降级优先级验证")
			gt := p.GinkgoTest{}
			fallbackEnable := gt.Get_Key_fallback_consumer_service_enable(p.ServerName)
			policy := gt.Get_Key_consumer_service_fallbackpolicy(p.ServerName)

			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = true
			data[fallbackEnable] = false
			data[policy] = p.Value_fallbackpolicy_exception
			data[p.Key_isolation_global_timeout_able] = true
			data[p.Key_isolation_global_timeoutInMilliseconds] = 100
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.ToRestInvoke("1000","rest")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			log.Println(string(body))
			log.Println(resp.StatusCode)

			log.Println("删除配置项")
			key := []string{p.Key_fallback_consumer_enable,fallbackEnable,policy,p.Key_isolation_global_timeout_able,p.Key_isolation_global_timeoutInMilliseconds}
			ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			log.Println("删除完毕")
			Expect(string(body)).To(Equal(""))
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.fallback.01",func(){
			p.PrintlnMegenta("HIGHWAY:01 容错降级验证")

			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = true
			data[p.Key_consumer_fallbackpolicy] = p.Value_fallbackpolicy_exception
			data[p.Key_isolation_global_timeout_able] = true
			data[p.Key_isolation_global_timeoutInMilliseconds] = 100
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.ToHighwayInvoke("1000")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			log.Println(string(body))

			log.Println("删除配置项")
			key := []string{p.Key_fallback_consumer_enable,p.Key_consumer_fallbackpolicy,p.Key_isolation_global_timeout_able,p.Key_isolation_global_timeoutInMilliseconds}
			ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			log.Println("删除完毕")
			Expect(string(body)).To(Equal("{\"Message\":\"GoServerHoy is isolated because of error: timeout\"}"))
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.fallback.02",func(){
			p.PrintlnMegenta("HIGHWAY:02 容错降级服务级别验证")
			gt := p.GinkgoTest{}
			fallbackEnable := gt.Get_Key_fallback_consumer_service_enable(p.ServerName)
			policy := gt.Get_Key_consumer_service_fallbackpolicy(p.ServerName)

			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[fallbackEnable] = true
			data[policy] = p.Value_fallbackpolicy_exception
			data[p.Key_isolation_global_timeout_able] = true
			data[p.Key_isolation_global_timeoutInMilliseconds] = 100
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.ToHighwayInvoke("1000")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			log.Println(string(body))

			log.Println("删除配置项")
			key := []string{fallbackEnable,policy,p.Key_isolation_global_timeout_able,p.Key_isolation_global_timeoutInMilliseconds}
			ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			log.Println("删除完毕")
			Expect(string(body)).To(Equal("{\"Message\":\"GoServerHoy is isolated because of error: timeout\"}"))
		})
		It("PAAS.SIT.CSE.GOSDKThreshold.Highway.fallback.03",func(){
			p.PrintlnMegenta("HIGHWAY:03 容错降级优先级验证")
			gt := p.GinkgoTest{}
			fallbackEnable := gt.Get_Key_fallback_consumer_service_enable(p.ServerName)
			policy := gt.Get_Key_consumer_service_fallbackpolicy(p.ServerName)

			//创建配置项
			log.Println("创建配置项")
			data := make(map[string]interface{})
			data[p.Key_fallback_consumer_enable] = true
			data[fallbackEnable] = false
			data[policy] = p.Value_fallbackpolicy_exception
			data[p.Key_isolation_global_timeout_able] = true
			data[p.Key_isolation_global_timeoutInMilliseconds] = 100
			ccbody := p.CreateConfigToCC(p.ServiceName,data)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			time.Sleep(2*time.Second)

			log.Println("开始循环调用")
			url := p.ToHighwayInvoke("1000")
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			log.Println(string(body))

			log.Println("删除配置项")
			key := []string{p.Key_fallback_consumer_enable,fallbackEnable,policy,p.Key_isolation_global_timeout_able,p.Key_isolation_global_timeoutInMilliseconds}
			ccbody = p.DeleteConfigFromCC(p.ServiceName,key)
			Expect(string(ccbody)).To(Equal(p.Result_success))
			log.Println("删除完毕")
			Expect(string(body)).To(Equal("{\"Message\":\"timeout\"}"))
		})
	})
})
