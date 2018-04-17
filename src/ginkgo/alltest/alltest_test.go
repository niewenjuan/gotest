package alltest_test

import (
	p "gotest/basicUtil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"io/ioutil"
	"net/http"
	"time"
)

var _ = Describe("Alltest", func() {
	AfterEach(func() {
		By("TestcaseFail")
	})
	BeforeEach(func(){
		p.PrintlnYellow("\n\n\n")
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

	Context("Loadbalance",func(){
		Context("LB_REST",func(){
			It("PAAS.SIT.CSE.GOSDK.REST.loadbalance.01",func(){
				p.PrintlnMegenta("round to random")
				//创建配置项
				key_strategy := p.Key_strategy
				strategy := p.Value_strategy_round
				log.Println("create config items")
				data := make(map[string]interface{})
				data[key_strategy] = strategy
				ccbody := p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				log.Println("begin to cycle invoke ......")
				url := p.CycleInvokeUrl("rest","40","0")
				resp, err := http.Get(url)
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				result := p.StrHandle(string(body),",",strategy)
				Expect(resp.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))

				log.Println("update config items")
				strategy = p.Value_strategy_random
				data[key_strategy] = strategy
				ccbody = p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				log.Println("begin to cycle invoke ......")
				resp1, err1 := http.Get(url)
				Expect(err1).NotTo(HaveOccurred())
				defer resp1.Body.Close()
				body, _ = ioutil.ReadAll(resp1.Body)
				result = p.StrHandle(string(body),",",strategy)
				Expect(resp1.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))
			})
			It("PAAS.SIT.CSE.GOSDK.REST.loadbalance.02",func(){
				p.PrintlnMegenta("round to weight")
				//创建配置项
				key_strategy := p.Key_strategy
				strategy := p.Value_strategy_round
				log.Println("create config items")
				data := make(map[string]interface{})
				data[key_strategy] = strategy
				ccbody := p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				log.Println("begin to cycle invoke ......")
				url := p.CycleInvokeUrl("rest","40","0")
				resp, err := http.Get(url)
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				result := p.StrHandle(string(body),",",strategy)
				Expect(resp.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))

				log.Println("update config items")
				strategy = p.Value_strategy_weight
				data[key_strategy] = strategy
				ccbody = p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				url = p.WeightInvokeUrl("rest","40")
				log.Println("begin to cycle invoke ......")
				resp1, err1 := http.Get(url)
				Expect(err1).NotTo(HaveOccurred())
				defer resp1.Body.Close()
				body, _ = ioutil.ReadAll(resp1.Body)
				result = p.StrHandle(string(body),",",strategy)
				Expect(resp1.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))
			})
			It("PAAS.SIT.CSE.GOSDK.REST.loadbalance.03",func(){
				p.PrintlnMegenta("round to session")
				//创建配置项
				key_strategy := p.Key_strategy
				strategy := p.Value_strategy_round
				log.Println("create config items")
				data := make(map[string]interface{})
				data[key_strategy] = strategy
				ccbody := p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				log.Println("begin to cycle invoke ......")
				url := p.CycleInvokeUrl("rest","40","0")
				resp, err := http.Get(url)
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				result := p.StrHandle(string(body),",",strategy)
				Expect(resp.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))

				log.Println("update config items")
				strategy = p.Value_strategy_session
				data[key_strategy] = strategy
				ccbody = p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(2*time.Second)
				url = p.SessionInvokeUrl("rest","40")
				log.Println("begin to cycle invoke ......")
				resp1, err1 := http.Get(url)
				Expect(err1).NotTo(HaveOccurred())
				defer resp1.Body.Close()
				body, _ = ioutil.ReadAll(resp1.Body)
				result = p.StrHandle(string(body),",",strategy)
				Expect(resp1.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))
			})
			It("PAAS.SIT.CSE.GOSDK.REST.loadbalance.04",func(){
				p.PrintlnMegenta("random to weight")
				//创建配置项
				key_strategy := p.Key_strategy
				strategy := p.Value_strategy_random
				log.Println("create config items")
				data := make(map[string]interface{})
				data[key_strategy] = strategy
				ccbody := p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				log.Println("begin to cycle invoke ......")
				url := p.CycleInvokeUrl("rest","40","0")
				resp, err := http.Get(url)
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				result := p.StrHandle(string(body),",",strategy)
				Expect(resp.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))

				log.Println("update config items")
				strategy = p.Value_strategy_weight
				data[key_strategy] = strategy
				ccbody = p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				url = p.WeightInvokeUrl("rest","40")
				log.Println("begin to cycle invoke ......")
				resp1, err1 := http.Get(url)
				Expect(err1).NotTo(HaveOccurred())
				defer resp1.Body.Close()
				body, _ = ioutil.ReadAll(resp1.Body)
				result = p.StrHandle(string(body),",",strategy)
				Expect(resp1.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))
			})
			It("PAAS.SIT.CSE.GOSDK.REST.loadbalance.05",func(){
				p.PrintlnMegenta("random to round")
				//创建配置项
				key_strategy := p.Key_strategy
				strategy := p.Value_strategy_random
				log.Println("create config items")
				data := make(map[string]interface{})
				data[key_strategy] = strategy
				ccbody := p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				log.Println("begin to cycle invoke ......")
				url := p.CycleInvokeUrl("rest","40","0")
				resp, err := http.Get(url)
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				result := p.StrHandle(string(body),",",strategy)
				Expect(resp.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))

				log.Println("update config items")
				strategy = p.Value_strategy_round
				data[key_strategy] = strategy
				ccbody = p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				log.Println("begin to cycle invoke ......")
				resp1, err1 := http.Get(url)
				Expect(err1).NotTo(HaveOccurred())
				defer resp1.Body.Close()
				body, _ = ioutil.ReadAll(resp1.Body)
				result = p.StrHandle(string(body),",",strategy)
				Expect(resp1.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))
			})
			It("PAAS.SIT.CSE.GOSDK.REST.loadbalance.06",func(){
				p.PrintlnMegenta("random to session")
				//创建配置项
				key_strategy := p.Key_strategy
				strategy := p.Value_strategy_random
				log.Println("create config items")
				data := make(map[string]interface{})
				data[key_strategy] = strategy
				ccbody := p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				log.Println("begin to cycle invoke ......")
				url := p.CycleInvokeUrl("rest","40","0")
				resp, err := http.Get(url)
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				result := p.StrHandle(string(body),",",strategy)
				Expect(resp.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))

				log.Println("update config items")
				strategy = p.Value_strategy_session
				data[key_strategy] = strategy
				ccbody = p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(2*time.Second)
				url = p.SessionInvokeUrl("rest","40")
				log.Println("begin to cycle invoke ......")
				resp1, err1 := http.Get(url)
				Expect(err1).NotTo(HaveOccurred())
				defer resp1.Body.Close()
				body, _ = ioutil.ReadAll(resp1.Body)
				result = p.StrHandle(string(body),",",strategy)
				Expect(resp1.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))
			})
			It("PAAS.SIT.CSE.GOSDK.REST.loadbalance.07",func(){
				p.PrintlnMegenta("weight to session")
				//创建配置项
				key_strategy := p.Key_strategy
				strategy := p.Value_strategy_weight
				log.Println("create config items")
				data := make(map[string]interface{})
				data[key_strategy] = strategy
				ccbody := p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				log.Println("begin to cycle invoke ......")
				url := p.WeightInvokeUrl("rest","40")
				resp, err := http.Get(url)
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				result := p.StrHandle(string(body),",",strategy)
				Expect(resp.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))

				log.Println("update config items")
				strategy = p.Value_strategy_session
				data[key_strategy] = strategy
				ccbody = p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(2*time.Second)
				url = p.SessionInvokeUrl("rest","40")
				log.Println("begin to cycle invoke ......")
				resp1, err1 := http.Get(url)
				Expect(err1).NotTo(HaveOccurred())
				defer resp1.Body.Close()
				body, _ = ioutil.ReadAll(resp1.Body)
				result = p.StrHandle(string(body),",",strategy)
				Expect(resp1.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))
			})
			It("PAAS.SIT.CSE.GOSDK.REST.loadbalance.08",func(){
				p.PrintlnMegenta("weight to random")
				//创建配置项
				key_strategy := p.Key_strategy
				strategy := p.Value_strategy_weight
				log.Println("create config items")
				data := make(map[string]interface{})
				data[key_strategy] = strategy
				ccbody := p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				log.Println("begin to cycle invoke ......")
				url := p.WeightInvokeUrl("rest","40")
				resp, err := http.Get(url)
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				result := p.StrHandle(string(body),",",strategy)
				Expect(resp.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))

				log.Println("update config items")
				strategy = p.Value_strategy_random
				data[key_strategy] = strategy
				ccbody = p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				log.Println("begin to cycle invoke ......")
				url = p.CycleInvokeUrl("rest","40","0")
				resp1, err1 := http.Get(url)
				Expect(err1).NotTo(HaveOccurred())
				defer resp1.Body.Close()
				body, _ = ioutil.ReadAll(resp1.Body)
				result = p.StrHandle(string(body),",",strategy)
				Expect(resp1.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))
			})
			It("PAAS.SIT.CSE.GOSDK.REST.loadbalance.09",func(){
				p.PrintlnMegenta("weight to round")
				//创建配置项
				key_strategy := p.Key_strategy
				strategy := p.Value_strategy_weight
				log.Println("create config items")
				data := make(map[string]interface{})
				data[key_strategy] = strategy
				ccbody := p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				log.Println("begin to cycle invoke ......")
				url := p.WeightInvokeUrl("rest","40")
				resp, err := http.Get(url)
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				result := p.StrHandle(string(body),",",strategy)
				Expect(resp.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))

				log.Println("update config items")
				strategy = p.Value_strategy_round
				data[key_strategy] = strategy
				ccbody = p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				log.Println("begin to cycle invoke ......")
				url = p.CycleInvokeUrl("rest","40","0")
				resp1, err1 := http.Get(url)
				Expect(err1).NotTo(HaveOccurred())
				defer resp1.Body.Close()
				body, _ = ioutil.ReadAll(resp1.Body)
				result = p.StrHandle(string(body),",",strategy)
				Expect(resp1.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))
			})
			It("PAAS.SIT.CSE.GOSDK.REST.loadbalance.10",func(){
				p.PrintlnMegenta("session to round")
				//创建配置项
				key_strategy := p.Key_strategy
				strategy := p.Value_strategy_session
				log.Println("create config items")
				data := make(map[string]interface{})
				data[key_strategy] = strategy
				ccbody := p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(2*time.Second)
				log.Println("begin to cycle invoke ......")
				url := p.SessionInvokeUrl("rest","40")
				resp, err := http.Get(url)
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				result := p.StrHandle(string(body),",",strategy)
				Expect(resp.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))

				log.Println("update config items")
				strategy = p.Value_strategy_round
				data[key_strategy] = strategy
				ccbody = p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				log.Println("begin to cycle invoke ......")
				url = p.CycleInvokeUrl("rest","40","0")
				resp1, err1 := http.Get(url)
				Expect(err1).NotTo(HaveOccurred())
				defer resp1.Body.Close()
				body, _ = ioutil.ReadAll(resp1.Body)
				result = p.StrHandle(string(body),",",strategy)
				Expect(resp1.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))
			})
			It("PAAS.SIT.CSE.GOSDK.REST.loadbalance.11",func(){
				p.PrintlnMegenta("session to weight")
				//创建配置项
				key_strategy := p.Key_strategy
				strategy := p.Value_strategy_session
				log.Println("create config items")
				data := make(map[string]interface{})
				data[key_strategy] = strategy
				ccbody := p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(2*time.Second)
				log.Println("begin to cycle invoke ......")
				url := p.SessionInvokeUrl("rest","40")
				resp, err := http.Get(url)
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				result := p.StrHandle(string(body),",",strategy)
				Expect(resp.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))

				log.Println("update config items")
				strategy = p.Value_strategy_weight
				data[key_strategy] = strategy
				ccbody = p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				url = p.WeightInvokeUrl("rest","40")
				log.Println("begin to cycle invoke ......")
				resp1, err1 := http.Get(url)
				Expect(err1).NotTo(HaveOccurred())
				defer resp1.Body.Close()
				body, _ = ioutil.ReadAll(resp1.Body)
				result = p.StrHandle(string(body),",",strategy)
				Expect(resp1.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))
			})
			It("PAAS.SIT.CSE.GOSDK.REST.loadbalance.12",func(){
				p.PrintlnMegenta("session to random")
				//创建配置项
				key_strategy := p.Key_strategy
				strategy := p.Value_strategy_session
				log.Println("create config items")
				data := make(map[string]interface{})
				data[key_strategy] = strategy
				ccbody := p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(2*time.Second)
				log.Println("begin to cycle invoke ......")
				url := p.SessionInvokeUrl("rest","40")
				resp, err := http.Get(url)
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				result := p.StrHandle(string(body),",",strategy)
				Expect(resp.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))

				log.Println("update config items")
				strategy = p.Value_strategy_random
				data[key_strategy] = strategy
				ccbody = p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				log.Println("begin to cycle invoke ......")
				url = p.CycleInvokeUrl("rest","40","0")
				resp1, err1 := http.Get(url)
				Expect(err1).NotTo(HaveOccurred())
				defer resp1.Body.Close()
				body, _ = ioutil.ReadAll(resp1.Body)
				result = p.StrHandle(string(body),",",strategy)
				Expect(resp1.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))
			})
		})
		Context("LB_HIGHWAY",func(){
			It("PAAS.SIT.CSE.GOSDK.HIGHWAY.loadbalance.01",func(){
				p.PrintlnMegenta("round to random")
				//创建配置项
				key_strategy := p.Key_strategy
				protocl := "highway"
				strategy := p.Value_strategy_round
				log.Println("create config items")
				data := make(map[string]interface{})
				data[key_strategy] = strategy
				ccbody := p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				log.Println("begin to cycle invoke ......")
				url := p.CycleInvokeUrl(protocl,"40","0")
				resp, err := http.Get(url)
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				result := p.StrHandle(string(body),",",strategy)
				Expect(resp.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))

				log.Println("update config items")
				strategy = p.Value_strategy_random
				data[key_strategy] = strategy
				ccbody = p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				log.Println("begin to cycle invoke ......")
				resp1, err1 := http.Get(url)
				Expect(err1).NotTo(HaveOccurred())
				defer resp1.Body.Close()
				body, _ = ioutil.ReadAll(resp1.Body)
				result = p.StrHandle(string(body),",",strategy)
				Expect(resp1.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))
			})
			It("PAAS.SIT.CSE.GOSDK.HIGHWAY.loadbalance.02",func(){
				p.PrintlnMegenta("round to weight")
				//创建配置项
				key_strategy := p.Key_strategy
				protocl := "highway"
				strategy := p.Value_strategy_round
				log.Println("create config items")
				data := make(map[string]interface{})
				data[key_strategy] = strategy
				ccbody := p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				log.Println("begin to cycle invoke ......")
				url := p.CycleInvokeUrl(protocl,"40","0")
				resp, err := http.Get(url)
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				result := p.StrHandle(string(body),",",strategy)
				Expect(resp.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))

				log.Println("update config items")
				strategy = p.Value_strategy_weight
				data[key_strategy] = strategy
				ccbody = p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				url = p.WeightInvokeUrl(protocl,"40")
				log.Println("begin to cycle invoke ......")
				resp1, err1 := http.Get(url)
				Expect(err1).NotTo(HaveOccurred())
				defer resp1.Body.Close()
				body, _ = ioutil.ReadAll(resp1.Body)
				result = p.StrHandle(string(body),",",strategy)
				Expect(resp1.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))
			})
			It("PAAS.SIT.CSE.GOSDK.HIGHWAY.loadbalance.03",func(){
				p.PrintlnMegenta("random to weight")
				//创建配置项
				key_strategy := p.Key_strategy
				protocl := "highway"
				strategy := p.Value_strategy_random
				log.Println("create config items")
				data := make(map[string]interface{})
				data[key_strategy] = strategy
				ccbody := p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				log.Println("begin to cycle invoke ......")
				url := p.CycleInvokeUrl(protocl,"40","0")
				resp, err := http.Get(url)
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				result := p.StrHandle(string(body),",",strategy)
				Expect(resp.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))

				log.Println("update config items")
				strategy = p.Value_strategy_weight
				data[key_strategy] = strategy
				ccbody = p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				url = p.WeightInvokeUrl(protocl,"40")
				log.Println("begin to cycle invoke ......")
				resp1, err1 := http.Get(url)
				Expect(err1).NotTo(HaveOccurred())
				defer resp1.Body.Close()
				body, _ = ioutil.ReadAll(resp1.Body)
				result = p.StrHandle(string(body),",",strategy)
				Expect(resp1.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))
			})
			It("PAAS.SIT.CSE.GOSDK.HIGHWAY.loadbalance.04",func(){
				p.PrintlnMegenta("random to round")
				//创建配置项
				key_strategy := p.Key_strategy
				protocl := "highway"
				strategy := p.Value_strategy_random
				log.Println("create config items")
				data := make(map[string]interface{})
				data[key_strategy] = strategy
				ccbody := p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				log.Println("begin to cycle invoke ......")
				url := p.CycleInvokeUrl(protocl,"40","0")
				resp, err := http.Get(url)
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				result := p.StrHandle(string(body),",",strategy)
				Expect(resp.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))

				log.Println("update config items")
				strategy = p.Value_strategy_round
				data[key_strategy] = strategy
				ccbody = p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				log.Println("begin to cycle invoke ......")
				resp1, err1 := http.Get(url)
				Expect(err1).NotTo(HaveOccurred())
				defer resp1.Body.Close()
				body, _ = ioutil.ReadAll(resp1.Body)
				result = p.StrHandle(string(body),",",strategy)
				Expect(resp1.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))
			})
			It("PAAS.SIT.CSE.GOSDK.HIGHWAY.loadbalance.05",func(){
				p.PrintlnMegenta("weight to random")
				//创建配置项
				key_strategy := p.Key_strategy
				protocl := "highway"
				strategy := p.Value_strategy_weight
				log.Println("create config items")
				data := make(map[string]interface{})
				data[key_strategy] = strategy
				ccbody := p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				log.Println("begin to cycle invoke ......")
				url := p.WeightInvokeUrl(protocl,"40")
				resp, err := http.Get(url)
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				result := p.StrHandle(string(body),",",strategy)
				Expect(resp.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))

				log.Println("update config items")
				strategy = p.Value_strategy_random
				data[key_strategy] = strategy
				ccbody = p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				log.Println("begin to cycle invoke ......")
				url = p.CycleInvokeUrl(protocl,"40","0")
				resp1, err1 := http.Get(url)
				Expect(err1).NotTo(HaveOccurred())
				defer resp1.Body.Close()
				body, _ = ioutil.ReadAll(resp1.Body)
				result = p.StrHandle(string(body),",",strategy)
				Expect(resp1.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))
			})
			It("PAAS.SIT.CSE.GOSDK.HIGHWAY.loadbalance.06",func(){
				p.PrintlnMegenta("weight to round")
				//创建配置项
				key_strategy := p.Key_strategy
				protocl := "highway"
				strategy := p.Value_strategy_weight
				log.Println("create config items")
				data := make(map[string]interface{})
				data[key_strategy] = strategy
				ccbody := p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(1*time.Second)
				log.Println("begin to cycle invoke ......")
				url := p.WeightInvokeUrl(protocl,"40")
				resp, err := http.Get(url)
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				result := p.StrHandle(string(body),",",strategy)
				Expect(resp.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))

				log.Println("update config items")
				strategy = p.Value_strategy_round
				data[key_strategy] = strategy
				ccbody = p.CreateConfigToCC(p.ServiceName,data)
				Expect(string(ccbody)).To(Equal(p.Result_success))
				time.Sleep(2*time.Second)
				log.Println("begin to cycle invoke ......")
				url = p.CycleInvokeUrl(protocl,"40","0")
				resp1, err1 := http.Get(url)
				Expect(err1).NotTo(HaveOccurred())
				defer resp1.Body.Close()
				body, _ = ioutil.ReadAll(resp1.Body)
				result = p.StrHandle(string(body),",",strategy)
				Expect(resp1.StatusCode).To(Equal(200))
				Expect(result).To(Equal(strategy))
			})
		})

	})
})
