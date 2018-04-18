package main


import (
	//"time"
	//"fmt"

	//"strings"
	//"gotest/basicUtil"
	//"fmt"
	"fmt"
	"net/http"
	"io"
	//"os"
	"strings"
	"io/ioutil"
	"encoding/json"
	"crypto/tls"
	"log"
	"github.com/cairixian/gotest/basicUtil"
	"strconv"
	"os"
)

var token string

func main(){
	////basicUtil.Generate_Randnum()
	//log.Print(basicUtil.ANSI_COLOR_LIGHT_BLUE)
	//log.Print("test")
	//log.Println(basicUtil.ANSI_COLOR_LIGHT_RESET)
	//basicUtil.PrintlnGreen("test one")
	//time.Sleep(10*time.Second)
	//basicUtil.PrintlnRed()
	//log.Println("test two")
	//basicUtil.DelayTime(10000)
	//log.Println("test tone")
	//log.Print()
	//basicUtil.PrintlnRed("sdsdsdsd")
	//temp :="instant_id:\"bl2gWB3\""
	////temp =strings.Replace(temp,"\"","",-1)
	//log.Print(temp)
	//defer log.Println("sdsdsd",10)
	//defer log.Print("sdsdsdsssss")
	//defer func() {     //必须要先声明defer，否则不能捕获到panic异常
	//	fmt.Println("c")
	//	if err := recover(); err != nil {
	//		fmt.Println(err)    //这里的err其实就是panic传入的内容，55
	//	}
	//	fmt.Println("d")
	//}()
	//panic("2323")
	//log.Println("test")

	//http.HandleFunc("/hello", HelloServer)
	//err := http.ListenAndServe("0.0.0.0:4099", nil)
	//if err != nil {
	//	log.Fatal("ListenAndServe: ", err)
	//}
	//basicUtil.PrintlnRed("test")
	//basicUtil.PrintlnRed("test")
	//log.Println(string(toUrl("cse.loadbalance.strategy.name", "RoundRobin", "Client@default#0.1")))
	//var result string
	//res := basicUtil.CreateConfigToCC("cse.loadbalance.strategy.name","Random","ClientHoy@default#1.0.1")
	//result = string(res)
	//log.Println(result)
	//var t1 time.Time
	//log.Println(t1)
	//time.Sleep(120*time.Second)
	//log.Println(time.Since(t1))
	//log.Println(time.Since(t1)>100)
	//myTest()
	//log.Println(basicUtil.Stoken == "")
	//key := make([]string,2)
	//key[0] = "sssssss"
	//key[1] = "sdddddd"
	//log.Println(key)
	//time := time.Now()
	//log.Println(time)
	//var dat map[string]map[string]string
	//arr := make([]string,0)
	////var date map[string]string
	//if err := json.Unmarshal(basicUtil.QueryConfigFromCC("ClientHoy"), &dat); err == nil {
	//	fmt.Println(dat)
	//	if v, ok := dat["ClientHoy"]; ok {
	//		for k, _ := range v {
	//			log.Println(k)
	//			arr = append(arr,k)
	//		}
	//		ccbody := basicUtil.DeleteConfigFromCC("ClientHoy",arr)
	//		log.Println(string(ccbody))
	//	}
	//} else {
	//	fmt.Println(err)
	//}
	gt :=basicUtil.ServiceCenter{}
	//dat := basicUtil.JsonExchange(gt.QueryExistenceFromSC("default","ClientHoy","1.0.0"))
	//fmt.Println(string(gt.QueryInstanceIdFromSC(dat["serviceId"])))
	//var res float64 = 1.0/2.0
	//ss := basicUtil.Read3(basicUtil.LogPath)
	//log.Println(res)
	//log.Println(ss[len(ss)-3])
	//log.Println(ss[len(ss)-2])
	//log.Println(ss[len(ss)-1])
	//min := time.Now().Minute()
	//sec := time.Now().Second()
	//log.Println(min,":",sec)
	//body :=make(map[string]interface{},0)
	//body["appId"] = "TestApphhf"
	//body["version"] = "1.0.0"
	//
	//for i :=0;i<460;i++ {
	//	str := strconv.Itoa(i)
	//	body["serviceName"] = "hhf"+str
	//	basicUtil.PrintlnCyan(string(basicUtil.CreateServiceToSc(body)))
	//}
	//sc := basicUtil.ServiceCenter{}
	//res := sc.QueryExistenceFromSC("default","GoServerHoy","1.0.1")
	//ss := basicUtil.JsonExchange(res)
	//log.Println(ss["serviceId"])
	//basicUtil.CurrentRequest()
	//ss := float64(55)/float64(100)
	//log.Println(float64(ss))
	//ss := basicUtil.JsonExchange()
	var dat map[string][]interface{}

	if err := json.Unmarshal(gt.QueryAllServicesFromSC(), &dat); err == nil {
		ss := dat["services"]
		for i :=0;i<len(ss);i++ {
			log.Println(ss[i])
			res,_ := json.Marshal(ss[i])
			aa := basicUtil.JsonExchange(res)
			log.Println(aa["serviceId"])
		}
		//log.Println(ss[0])
		//res,_ := json.Marshal((ss[0]))
		//log.Println(string(res))
		//aa := basicUtil.JsonExchange(res)
		//log.Println(aa["serviceId"])

	} else {
		log.Println(err)
	}
	//log.Println(ss[0])


}
func test(){
	token = "test"
}
func currentTest(ch chan string){
	ch<-"test"
}
func read3(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}

func myTest(){
	//Consumer_rest_IP := "127.0.0.1:2098"
	//data := make(map[string]interface{})
	//data["cse.loadbalance.strategy.name"] = "Random"
	//res := basicUtil.CreateConfigToCC("ClientHoy",data)
	delaytime := 3000
	u := basicUtil.CombinationURL("cse","%s%s","GoServerHoy",fmt.Sprintf("/demo/LoadBalanceTest/%s",strconv.Itoa(delaytime)))
	log.Println(u)
	log.Println(fmt.Sprintf("/sds/sdsd/%s","test"))

	//time.Sleep(1*time.Second)
	//url := "http://" + Consumer_rest_IP + "/consumer/simple/1"
	//resp, _ := http.Get(url)
	//defer resp.Body.Close()
	//body, _ := ioutil.ReadAll(resp.Body)
	//log.Println(string(body))
}
func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func Recover(){
	if err:=recover();err!=nil{
		fmt.Println(err) // 这里的err其实就是panic传入的内容，55
	}
}

func toUrl(key, value, dimensionInfo string)[]byte{
	type CreateConfigApi struct {
		DimensionInfo string                 `json:"dimensionsInfo"`
		Items         map[string]interface{} `json:"items"`
	}

	createConfig := new(CreateConfigApi)
	createConfig.DimensionInfo = dimensionInfo
	data := make(map[string]interface{}, 0)
	data[key] = value
	createConfig.Items = data
	config, err := json.Marshal(createConfig)
	if err != nil {
		fmt.Println("failed to marshal body")
		return nil
	}

	//url := os.Getenv("CSE_CONFIG_CENTER_ADDR") + "/v3/default/configuration/items"//
	url := "https://cse.cn-north-1.myhwclouds.com:443" + "/v3/default/configuration/items"
	payload := strings.NewReader(string(config))

	req, _ := http.NewRequest("POST", url, payload)

	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}
	client := &http.Client{Transport: transCfg}
	//client := &http.Client{}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-Auth-Token", "MIIPnQYJKoZIhvcNAQcCoIIPjjCCD4oCAQExDTALBglghkgBZQMEAgEwgg3rBgkqhkiG9w0BBwGggg3cBIIN2HsidG9rZW4iOnsiZXhwaXJlc19hdCI6IjIwMTgtMDEtMzFUMDI6NDM6NDMuNzgxMDAwWiIsIm1ldGhvZHMiOlsicGFzc3dvcmQiXSwiY2F0YWxvZyI6W10sInJvbGVzIjpbeyJuYW1lIjoidGVfYWRtaW4iLCJpZCI6IjE5OTJjMWRmOWFkNjQxMmU5YzgzMzAzMmNkNzBjYThmIn0seyJuYW1lIjoib3BfZ2F0ZWRfbW9sYXAiLCJpZCI6IjAifSx7Im5hbWUiOiJvcF9nYXRlZF9zY2Nfd2FmIiwiaWQiOiIwIn0seyJuYW1lIjoib3BfZ2F0ZWRfYWlzX29jcl92YXRfaW52b2ljZSIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX2Vjc19tMyIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX29wX2dhdGVkX3NjY19odmQiLCJpZCI6IjAifSx7Im5hbWUiOiJvcF9nYXRlZF9HQUNTIiwiaWQiOiIwIn0seyJuYW1lIjoib3BfZ2F0ZWRfYWlzX29jcl9oYW5kd3JpdGluZyIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX2Nlc19hZ3QiLCJpZCI6IjAifSx7Im5hbWUiOiJvcF9nYXRlZF9lY3NfaTMiLCJpZCI6IjAifSx7Im5hbWUiOiJvcF9nYXRlZF9kb21haW4iLCJpZCI6IjAifSx7Im5hbWUiOiJvcF9nYXRlZF9jb2xkIiwiaWQiOiIwIn0seyJuYW1lIjoib3BfZ2F0ZWRfcmRzX215b3B0IiwiaWQiOiIwIn0seyJuYW1lIjoib3BfZ2F0ZWRfbmV3ZGRzIiwiaWQiOiIwIn0seyJuYW1lIjoib3BfZ2F0ZWRfb3BfZ2F0ZWRfc2NjX3NjcyIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX2Rkc19yZXBsaWNhIiwiaWQiOiIwIn0seyJuYW1lIjoib3BfZ2F0ZWRfc2NjX2FycyIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX2Vjc19ub3JtYWxleGNsdXNpdmUiLCJpZCI6IjAifSx7Im5hbWUiOiJvcF9nYXRlZF9zY2Nfc2FzIiwiaWQiOiIwIn0seyJuYW1lIjoib3BfZ2F0ZWRfY2xvdWRjYyIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX2Vjc3F1aWNrZGVwbG95IiwiaWQiOiIwIn0seyJuYW1lIjoib3BfZ2F0ZWRfcmRzX2h3c3FsIiwiaWQiOiIwIn0seyJuYW1lIjoib3BfZ2F0ZWRfcnRzIiwiaWQiOiIwIn0seyJuYW1lIjoib3BfZ2F0ZWRfY3NicyIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX3dhZiIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX2Vjc19kZWgiLCJpZCI6IjAifSx7Im5hbWUiOiJvcF9nYXRlZF9haXNfaW1hZ2VfY2xhcml0eV9kZXRlY3QiLCJpZCI6IjAifSx7Im5hbWUiOiJvcF9nYXRlZF9jaGFuZ2hlbjEyMyIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX2t2bSIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX3NjY193dHAiLCJpZCI6IjAifSx7Im5hbWUiOiJvcF9nYXRlZF9sbGQiLCJpZCI6IjAifSx7Im5hbWUiOiJvcF9nYXRlZF9haXNfb2NyX2dlbmVyYWxfdGFibGUiLCJpZCI6IjAifSx7Im5hbWUiOiJvcF9nYXRlZF9JTSIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX2Vjc19yZWN5Y2xlYmluIiwiaWQiOiIwIn0seyJuYW1lIjoib3BfZ2F0ZWRfTFRTIiwiaWQiOiIwIn0seyJuYW1lIjoib3BfZ2F0ZWRfY2FkIiwiaWQiOiIwIn0seyJuYW1lIjoib3BfZ2F0ZWRfYWlzX29jcl92ZWhpY2xlX2xpY2Vuc2UiLCJpZCI6IjAifSx7Im5hbWUiOiJvcF9nYXRlZF9oc3MiLCJpZCI6IjAifSx7Im5hbWUiOiJvcF9nYXRlZF9mZ3MiLCJpZCI6IjAifSx7Im5hbWUiOiJvcF9nYXRlZF9kY3MiLCJpZCI6IjAifSx7Im5hbWUiOiJvcF9nYXRlZF9kd3NfZmVhdHVyZSIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX2Vjc19ldDIiLCJpZCI6IjAifSx7Im5hbWUiOiJvcF9sZWdhY3kiLCJpZCI6IjAifSx7Im5hbWUiOiJvcF9nYXRlZF9jbG91ZElNIiwiaWQiOiIwIn0seyJuYW1lIjoib3BfZ2F0ZWRfYWlzX29jcl9pZF9jYXJkIiwiaWQiOiIwIn0seyJuYW1lIjoib3BfZ2F0ZWRfc2NjX2h2ZCIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX29wX2dhdGVkX3Jkc19jdXN0b21lcmNsb3VkIiwiaWQiOiIwIn0seyJuYW1lIjoib3BfZ2F0ZWRfZWNzX25vcm1hbF9zMyIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX2dwdSIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX2dhdGVkX2ttcyIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX2Fpc19pbWFnZV9yZWNhcHR1cmVfZGV0ZWN0IiwiaWQiOiIwIn0seyJuYW1lIjoib3BfZ2F0ZWRfZWNzX2hwYyIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX2Fpc19vY3JfZHJpdmVyX2xpY2Vuc2UiLCJpZCI6IjAifSx7Im5hbWUiOiJvcF9nYXRlZF9vcF9nYXRlZF9zY2Nfd2FmIiwiaWQiOiIwIn0seyJuYW1lIjoib3BfZ2F0ZWRfbWxzIiwiaWQiOiIwIn0seyJuYW1lIjoib3BfZ2F0ZWRfYnV6aGkyMyIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX0Z1bmN0aW9uR3JhcGgiLCJpZCI6IjAifSx7Im5hbWUiOiJvcF9nYXRlZF90YWciLCJpZCI6IjAifSx7Im5hbWUiOiJvcF9nYXRlZF90bXMiLCJpZCI6IjAifSx7Im5hbWUiOiJvcF9nYXRlZF9kY3NfY2x1c3RlciIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX2Vjc19obWVtIiwiaWQiOiIwIn0seyJuYW1lIjoib3BfZ2F0ZWRfc2NjX3NzYSIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX2RycyIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX2Vjc19oMyIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX2Fpc19pbWFnZV9hbnRpcG9ybiIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX3dlYnNjYW4iLCJpZCI6IjAifSx7Im5hbWUiOiJvcF9nYXRlZF9haXNfb2NyX2N1c3RvbV9mb3JtIiwiaWQiOiIwIn0seyJuYW1lIjoib3BfZ2F0ZWRfbmF0IiwiaWQiOiIwIn0seyJuYW1lIjoib3BfZ2F0ZWRfb3BfZ2F0ZWRfc2NjX3B0cyIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX2Vjc19kMyIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX2VsYXN0aWNzZWFyY2giLCJpZCI6IjAifSx7Im5hbWUiOiJvcF9nYXRlZF9lY3NfZGlza2ludGVuc2l2ZSIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX25hdGd3IiwiaWQiOiIwIn0seyJuYW1lIjoib3BfZ2F0ZWRfYWlzX2ltYWdlX3RhZ2dpbmciLCJpZCI6IjAifSx7Im5hbWUiOiJvcF9nYXRlZF9hcGlnIiwiaWQiOiIwIn0seyJuYW1lIjoib3BfZ2F0ZWRfZGNzX21lbWNhY2hlZCIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX3Jkcy10cmFuc2ZlciIsImlkIjoiMCJ9LHsibmFtZSI6Im9wX2dhdGVkX2Fpc19hc3Jfc2VudGVuY2UiLCJpZCI6IjAifV0sInByb2plY3QiOnsiZG9tYWluIjp7Im5hbWUiOiJjYW95b25nc2hlbmciLCJpZCI6ImIxZmYwOTJhODMxMTQxYTk4NWRlMDFlNjNjNjc5MzQ5In0sIm5hbWUiOiJzb3V0aGNoaW5hIiwiaWQiOiJhMzgzMmQ5NDlkNjA0MmMxYjBhYTkzYmVlYjE0ODI2MSJ9LCJpc3N1ZWRfYXQiOiIyMDE4LTAxLTMwVDAyOjQzOjQzLjc4MTAwMFoiLCJ1c2VyIjp7ImRvbWFpbiI6eyJuYW1lIjoiY2FveW9uZ3NoZW5nIiwiaWQiOiJiMWZmMDkyYTgzMTE0MWE5ODVkZTAxZTYzYzY3OTM0OSJ9LCJuYW1lIjoiY2FveW9uZ3NoZW5nIiwiaWQiOiIxNDQ5ZTZiMDY0YzM0ZTYyOWVlOTE2Y2UxZTYzYzFmZSJ9fX0xggGFMIIBgQIBATBcMFcxCzAJBgNVBAYTAlVTMQ4wDAYDVQQIDAVVbnNldDEOMAwGA1UEBwwFVW5zZXQxDjAMBgNVBAoMBVVuc2V0MRgwFgYDVQQDDA93d3cuZXhhbXBsZS5jb20CAQEwCwYJYIZIAWUDBAIBMA0GCSqGSIb3DQEBAQUABIIBACn1bxLHPx4nGjHC42MxbnuFJUWBTERsHJxDH+ZJDpOn7j9tiGeHsRnQdJJLPL4RUQm3gnO7A0r5Gm4wgo7GxNSGW99zL+79PrAgmFRU+CwEMIBnnAsakWco4+-vRQ6Ge6-0nfz0tVhi7bSxHNS-oSx5tI2ZMdidNY89gyGZkfcZHch66VMbvFO5EpaQf9eWb2BrgLSaNwS5IuNu-PtIdp6A8cZnRIkChs7lEyRN5l+dHXMl7RUoUvFRiiTGEIbADF0+tuu3cdLObTgYaCz6R-KRBs-wGyB5d7OR7WUDeKiJ9KP4FUycbMDbUQImcSuZLNdCt8VGzJxpG5j3rXs6Zk8=")

	res, _ := client.Do(req)

	body, _ := ioutil.ReadAll(res.Body)

	defer res.Body.Close()
	return body
}

func getToken()string{

	type CreateConfigApi struct {
		Auth    map[string]interface{} `json:"auth"`
	}

	//auth
	domain := make(map[string]interface{}, 0)
	domain["name"] = "caoyongsheng"
	user := make(map[string]interface{}, 0)
	user["domain"] = domain
	user["password"] = "Huawei@123"
	user["name"] = "caoyongsheng"
	password := make(map[string]interface{}, 0)
	password["user"] = user

	met := []string{"password"}
	//methods := make(map[string][]string, 0)
	//methods["methods"] = met
	identity := make(map[string]interface{}, 0)
	identity["methods"] = met
	identity["password"] = password

	//scope
	project := make(map[string]interface{}, 0)
	data := make(map[string]interface{}, 0)
	data["name"] = "southchina"
	project["project"]= data
	//scope := make(map[string]interface{}, 0)
	//scope["scope"] = project

	auth := make(map[string]interface{}, 0)
	auth["identity"] = identity
	auth["scope"] = project

	createConfig := new(CreateConfigApi)
	createConfig.Auth =auth

	config, err := json.Marshal(createConfig)
	basicUtil.PrintlnGreen(string(config))
	if err != nil {
		fmt.Println("failed to marshal body",err)
	}

	url := "https://192.144.1.37:31943/v3/auth/tokens"
	payload := strings.NewReader(string(config))

	req, _ := http.NewRequest("POST", url, payload)

	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}
	client := &http.Client{Transport: transCfg}

	req.Header.Add("content-type", "application/json")
	res, _ := client.Do(req)

	header := res.Header
	token = header.Get("x-subject-token")
	return token
}


