package basicUtil

import (
	"log"
	"github.com/ServiceComb/go-chassis/core"
	//"context"
	//"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/metadata"
	"github.com/ServiceComb/go-chassis/client/rest"
	"time"
	"math/rand"
	"fmt"
	"strings"
	"io/ioutil"
	"encoding/json"
	"crypto/tls"
	"net/http"
	"golang.org/x/net/context"
	"strconv"
	"os"
)
var highwayInvoker = core.NewRPCInvoker()
var restInvoker = core.NewRestInvoker()
var ctx = context.Background()
//var ctx = metadata.NewContext(, map[string]string{"user-name":"Hoy",})
var stoken string
var cookies string
var cookieId []string
var res string

//error 异常打印
func Isok(err error){
	if err != nil {
		log.Println("err is",err)
	}
}

func CombinationURL(protocol, format string, v ...interface{}) string {
	return fmt.Sprintf("%s://%s", protocol, fmt.Sprintf(format, v...))
}

func HighwayInvoke(serviceName string,schemaId string,operationId string,arg interface{},reply interface{},options ...core.InvocationOption)(interface{}){
	err := highwayInvoker.Invoke(ctx,serviceName,schemaId,operationId,arg,reply,options...)
	if err !=nil {
		return err
	} else {
		return reply
	}
}

func RestfulInvoke(method string,url string,body []byte,options ...core.InvocationOption)string{
	//log.Println(url)
	restQuest, err := rest.NewRequest(method,url,body)
	Isok(err)
	result, err1 := restInvoker.ContextDo(ctx,restQuest,options...)
	Isok(err1)
	res = string(result.ReadBody())
	//log.Println(string(result.GetCookie("ServiceCombLB")))
	restQuest.Close()
	result.Close()
	return res
}

func RestfulInvokeSession(cookie,method,url string,body []byte,options ...core.InvocationOption)string{
	//log.Println(url)
	restQuest, err := rest.NewRequest(method,url,body)
	Isok(err)
	restQuest.SetCookie("ServiceCombLB",cookie)
	log.Println(restQuest.GetCookie("ServiceCombLB"))
	result, err1 := restInvoker.ContextDo(ctx,restQuest,options...)
	Isok(err1)
	res = string(result.ReadBody())
	//log.Println(result.GetStatusCode())
	restQuest.Close()
	result.Close()
	return res
}


//通道单并发访问访问
func currentTest(protocol string,ch chan string){
	var url string
	if protocol == "rest"{
		url = ToRestInvoke("1","rest")
	} else {
		url = ToHighwayInvoke("1")
	}
	resp, _ := http.Get(url)
	body, _ := ioutil.ReadAll(resp.Body)
	ch<-string(body)
}

//多并发请求下发与返回处理
func CurrentRequest(protocol string,maxNum,limitNum int)bool{
	flag := 0
	chs :=make([]chan string,maxNum)
	for i:=0;i<maxNum;i++{
		chs[i] =make(chan string)
		go currentTest(protocol,chs[i])
	}
	if protocol == "rest" {
		for _, v := range chs {
			res := <-v
			log.Println(res)
			if res == "a" || res == "b" {
				flag++
			}
		}
	}else {
		flag1 :=0
		for _, v := range chs {
			res := <-v
			log.Println(res)
			if res == "{\"Message\":\"max concurrency\"}" {
				flag1++
			}
		}
		flag = maxNum - flag1
	}
	PrintlnGreen(flag)
	if flag <= limitNum +2 {
		return true
	} else {
		return false
	}
}


//生成0-100的随机数字
//func Generate_Randnum() int{
//	rand.Seed(time.Now().Unix())
//	rnd := rand.Intn(100)
//	return rnd
//}

//生成num长度的随机字符串
func  Generate_RandString(num int) string {
	bytes := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < num; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	log.Println(string(result))
	return string(result)
}

//延时时间，秒（ms）
func DelayTime(delaytime int64){
	time.Sleep(time.Duration(delaytime)*time.Millisecond)
}

//延时时间，秒（ms）
func DelayTimeInt(delaytime int){
	time.Sleep(time.Duration(delaytime)*time.Millisecond)
}

//读取日志
func Read3(path string) []string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	result := string(fd)
	arr := strings.Split(result,"\n")
	log := arr[0:len(arr)-1]
	return log
}

func CreateConfigToCC(dimensionInfo string,item map[string]interface{})[]byte{
	type CreateConfigApi struct {
		DimensionInfo string                 `json:"dimensionsInfo"`
		Items         map[string]interface{} `json:"items"`
	}

	createConfig := new(CreateConfigApi)
	createConfig.DimensionInfo = dimensionInfo    //服务名
	createConfig.Items = item                     //配置项
	config, err := json.Marshal(createConfig)
	if err != nil {
		fmt.Println("failed to marshal body")
		return nil
	}

	url := CCUrl
	PrintlnGreen("body is:\n"+string(config))
	return requestDo("POST",url,string(config))
}

//创建空服务定义
func CreateService(scname,appid,version string)string{
	sc := ServiceCenter{}
	body :=make(map[string]interface{},0)
	body["appId"] = appid
	body["version"] = version
	body["serviceName"]= scname
	res := JsonExchange(sc.CreateServiceToSc(body))
	return res["serviceId"]
}

//创建k空实例
func CreateInstances(serviceID,hostName,endpoints string)string{
	sc := ServiceCenter{}
	body := make(map[string]interface{},0)
	endpoint := make([]string,1)
	endpoint[0] = endpoints
	body["hostName"] = hostName
	body["status"] = "UP"
	body["endpoints"] = endpoint
	res := sc.CreateInstanceIdFromSC(serviceID,body)
	return string(res)
}

//删除配置项
func DeleteConfigFromCC(dimensionInfo string,key []string)[]byte{
	type DeleteConfigApi struct {
		DimensionInfo string                 `json:"dimensionsInfo"`
		Keys []string  `json:"keys"`
	}

	deleteConfig := new(DeleteConfigApi)
	deleteConfig.DimensionInfo = dimensionInfo
	deleteConfig.Keys = key
	config, err := json.Marshal(deleteConfig)
	if err != nil {
		fmt.Println("failed to marshal body")
		return nil
	}

	url := CCUrl
	return requestDo("DELETE",url,string(config))
}

//查询配置项
func QueryConfigFromCC(dimensionInfo string)[]byte{
	url := CCUrl + "?dimensionsInfo=" + dimensionInfo
	return requestDo("GET",url,"")
}

//请求下发
func requestDo(method, url, requestbody string)[]byte{
	payload := strings.NewReader(requestbody)
	req, _ := http.NewRequest(method, url, payload)
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}
	client := &http.Client{Transport: transCfg}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-Auth-Token", autoToken())

	res, _ := client.Do(req)

	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	return body
}

//请求下发（header包含）
func requestDoA(method, url, requestbody, serviceId string)[]byte{
	payload := strings.NewReader(requestbody)
	req, _ := http.NewRequest(method, url, payload)
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}
	client := &http.Client{Transport: transCfg}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-Auth-Token", autoToken())
	req.Header.Add("X-ConsumerId",serviceId)

	res, _ := client.Do(req)

	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	return body
}

func JsonExchange(message []byte)(map[string]string){
	var dat map[string]string
	if err := json.Unmarshal(message, &dat); err == nil {
		return dat
	} else {
		return nil
	}
}
func JsonExchangeInter(message []byte)(map[string]interface{}){
	var dat map[string]interface{}
	if err := json.Unmarshal(message, &dat); err == nil {
		return dat
	} else {
		return nil
	}
}

func JsonMarshalExhange(v interface{})(map[string]string){
	res,err := json.Marshal(v)
	if err == nil {
		result := JsonExchange(res)
		return result
	} else {
		return nil
	}
}


//循环删除某个服务下所有配置项
func DeleteAllConfigForService(serviceName string){
	var dat map[string]map[string]interface{}
	arr := make([]string,0)
	if err := json.Unmarshal(QueryConfigFromCC(serviceName), &dat); err == nil {
		if v, ok := dat[serviceName]; ok {
			for k, _ := range v {
				arr = append(arr,k)
			}
			ccbody := DeleteConfigFromCC(serviceName,arr)
			log.Println(string(ccbody))
		}
	} else {
		PrintlnRed(err)
	}
}

//获取token
func getToken()string{

	type CreateConfigApi struct {
		Auth    map[string]interface{} `json:"auth"`
	}

	//auth
	domain := make(map[string]interface{}, 0)
	domain["name"] = Token_DomainName                 //修改租户
	user := make(map[string]interface{}, 0)
	user["domain"] = domain
	user["password"] = Token_PassWord                 //租户密码
	user["name"] = Token_UserName                   //修改用户
	password := make(map[string]interface{}, 0)
	password["user"] = user
	met := []string{"password"}
	identity := make(map[string]interface{}, 0)
	identity["methods"] = met
	identity["password"] = password

	//scope
	project := make(map[string]interface{}, 0)
	data := make(map[string]interface{}, 0)
	data["name"] = Token_ProjectName                    //修改project name
	project["project"]= data

	auth := make(map[string]interface{}, 0)
	auth["identity"] = identity
	auth["scope"] = project

	createConfig := new(CreateConfigApi)
	createConfig.Auth =auth

	config, err := json.Marshal(createConfig)
	//PrintlnGreen(string(config))
	if err != nil {
		fmt.Println("failed to marshal body",err)
	}

	url := GetTokenUrl
	payload := strings.NewReader(string(config))

	req, _ := http.NewRequest("POST", url, payload)

	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}
	client := &http.Client{Transport: transCfg}

	req.Header.Add("content-type", "application/json")
	res, _ := client.Do(req)

	header := res.Header
	token := header.Get("x-subject-token")
	return token
}

func autoToken()string{
	if stoken == "" {
		return getToken()
	} else {
		return stoken
	}
}

func logicForStatry(array []string,stratycheck string)string{
	a :=0
	b :=0
	same := 0
	if array[0] == "a"||array[0] == "b" {
		if stratycheck == "Random" {
			if array[0] == array[1]  {
				return "Random"
			}else {
				for i:=0;i<len(array)-1;i++ {
					if array[i] == array[i+1] {
						same++
					}
				}
				if same >0 {
					return "Random"
				} else {
					return "NoMatch in 'same'"
				}
			}
		} else if stratycheck == "RoundRobin" {
			if array[0] == array[1] {
				return "RoundRobin should notbe tested in array[0] == array[1] "
			} else {
				for i:=0;i<len(array)-1;i++ {
					if array[i] == array[i+1] {
						same++
					}
				}
				if same ==0 {
					return "RoundRobin"
				} else {
					return "RoundRobin should notbe tested in array[i] == array[i+1] "
				}
			}
		} else if stratycheck == "SessionStickiness" {
			if array[0] != array[1] {
				return "SessionStickiness is based on array[0] = array[1]"
			} else {
				for i:=0;i<(len(array)-1);i++ {
					if array[i] == array[i+1] {
						same++
					}
				}
				if same == (len(array)-1) {
					return "SessionStickiness"
				} else {
					log.Println("same:",same,"len(array)-1:",len(array)-1)
					return "SessionStickiness should be same everytime"
				}
			}
		} else if stratycheck == "WeightedResponse" {
			for i:=0;i<len(array);i++ {
				if array[i] != "" {
					if array[i] == "a" {
						a++
					}
					if array[i] == "b" {
						b++
					}
				}
			}
			if a>b && b !=0 {
				return "WeightedResponse"
			} else {
				log.Println("a:",a,"b:",b)
				return "a>b and b !=0 should be in WeightedResponse"
			}
		} else {
			return "stratycheck must be Random|RoundRobin|SessionStickiness|WeightedResponse"
		}
	} else {
		return "result is nil"
	}
}

func logForDarklanuch(array []string)(float64,float64){
	counta := 0
	countb := 0
	if array[0] != "a"&&array[0] != "b"&&array[0] !="c" {
		return 1.0,1.0
	}else {
		for _,v := range array {
			if v == "a"||v =="b" {
				counta++
			} else if v == "c"{
				countb++
			} else {
				PrintlnRed("result is an error!\n")
				break
			}
		}
		log.Println(counta,countb)
		var res float64 = float64(counta) / float64(len(array))
		var res1 float64 = float64(countb) / float64(len(array))
		return res,res1
	}
}
//负载策略判断逻辑
func StrHandle(bodystr ,seq ,stratycheck string) string {
	PrintlnCyan(bodystr)
	array1 := strings.Split(bodystr,seq)
	array := array1[:len(array1)-1]
	return  logicForStatry(array,stratycheck)
}

//流控返回body处理
func QpsHandle(bodystr ,seq string) int {
	PrintlnCyan(bodystr)
	flag := 0
	flag1 := 0
	minarr := make([]int,0)
	secarr := make([]int,0)
	array1 := strings.Split(bodystr,seq)
	for i :=1;i<len(array1)-1;i++ {
		arr := strings.Split(array1[i],":")
		min,_ := strconv.Atoi(arr[0])
		sec,_ := strconv.Atoi(arr[1])
		//log.Println(min,sec)
		minarr = append(minarr,min)
		secarr = append(secarr,sec)
	}
	for k :=0;k<len(secarr);k++ {
		if k >0 {
			if secarr[k]>secarr[k-1] {
				if secarr[k] - secarr[k - 1] == 2 ||secarr[k] - secarr[k - 1] == 1 {
				} else {
					PrintlnRed(secarr[k],secarr[k - 1])
					log.Println("mode1")
					flag++
				}
			}else if secarr[k] == secarr[k-1] {
				log.Println("mode2")
				PrintlnRed(secarr[k],secarr[k - 1])
				flag1++
			}else {
				if secarr[k] + 60  - secarr[k - 1] ==2 || secarr[k] + 60 - secarr[k - 1] == 1 {
				} else {
					PrintlnRed(secarr[k] + 60, secarr[k - 1])
					log.Println("mode3")
					flag++
				}
			}
		}
	}
	log.Println("flag is:",flag,"flag1 is:",flag1)
	if flag1 < len(secarr)/2 {
		return  flag
	} else {
		return flag1
	}
}

//流控不生效返回body处理
func NoQpsHandle(bodystr ,seq string) int {
	PrintlnCyan(bodystr)
	flag := 0
	minarr := make([]int,0)
	secarr := make([]int,0)
	array1 := strings.Split(bodystr,seq)
	for i :=1;i<len(array1)-1;i++ {
		arr := strings.Split(array1[i],":")
		min,_ := strconv.Atoi(arr[0])
		sec,_ := strconv.Atoi(arr[1])
		//log.Println(min,sec)
		minarr = append(minarr,min)
		secarr = append(secarr,sec)
	}
	for k :=0;k<len(secarr);k++ {
		if k >0 {
			if secarr[k] == secarr[k-1] {

			}else {
				flag++
			}
		}
	}
	log.Println("flag is:",flag)
	//array := array1[:len(array1)-1]
	return  flag
}

//灰度规则判断逻辑
func DarklanuchHandle(bodystr,seq string)bool{
	PrintlnCyan(bodystr)
	array1 := strings.Split(bodystr,seq)
	array := array1[:len(array1)-1]
	res,res1 := logForDarklanuch(array)
	PrintlnCyan(res,res1)
	if res>0.45&&res<0.5&&res1>0.5&&res1<0.55 {
		log.Println("mode1")
		return true
	} else if res1>0.45&&res1<0.5&&res>0.5&&res<0.55{
		log.Println("mode2")
		return true
	}else if res ==0.5&&res1 ==0.5 {
		log.Println("mode3")
		return true
	} else {
		return false
	}
	//return logForDarklanuch(array)
}

//在超时熔断模式下rest返回body判断
func CircuitHandle(bodystr,seq string,index int)bool{
	PrintlnCyan(bodystr)
	array1 := strings.Split(bodystr,seq)
	array := array1[1:len(array1)]
	log.Println(len(array))
	if array[index] == "" {
		return true
	}else{
		return false
	}

}

//日志固定读取某些行,尝试同实例(REST)
func LogReadInLines(arr []string,lineIndex int)bool{
	count := 0
	count1 :=0
	array := arr[len(arr)-lineIndex:]
	for _,v := range array {
		PrintlnCyan(v)
		if strings.Contains(v,"ERROR")&&strings.Contains(v,"500") {
			if strings.Contains(v, "FAILA") {
				count++
			}
			if strings.Contains(v, "FAILB") {
				count1++
			}
		}
	}
	if (count==lineIndex&&count1==0)||(count==0&&count1==lineIndex){
		return true
	}else{
		return false
	}
}

//日志固定读取某些行，进行字符串校验
func LogReadInLinesOperation(arr []string,lineIndex int,strCheck string)bool{
	flag := 0
	array := arr[len(arr)-lineIndex:]
	for _,v := range array {
		PrintlnCyan(v)
		if strings.Contains(v,strCheck) {
			flag++
		}
	}
	if flag !=0 {
		return true
	} else {
		return false
	}
}

//日志固定读取某些行，进行次数统计
func LogReadInLinesOperationMaxCurrent(arr []string,invokeNum,targetNum int,strCheck string)bool{
	flag := 0
	array := arr[len(arr)-invokeNum:]
	for _,v := range array {
		PrintlnCyan(v)
		if strings.Contains(v,strCheck) {
			flag++
		}
	}
	log.Println("flag is:",flag)
	if invokeNum-flag-targetNum<= 2 {
		return true
	} else {
		return false
	}
}

//日志固定读取某些行,尝试同实例(REST)
func LogReadInLinesOnnext(arr []string,lineIndex int)bool{
	count := 0
	count1 :=0
	array := arr[len(arr)-lineIndex:]
	for _,v := range array {
		PrintlnCyan(v)
		if strings.Contains(v,"ERROR")&&strings.Contains(v,"500") {
			if strings.Contains(v, "FAILA") {
				count++
			}
			if strings.Contains(v, "FAILB") {
				count1++
			}
		}
	}
	if (count<lineIndex&&count>0&&count1<lineIndex&&count1>0&&count+count1==lineIndex){
		return true
	}else{
		return false
	}
}

//日志固定读取某些行,尝试同实例(Highway)
func LogReadInLinesOnnextByHighway(arr []string,lineIndex int)int{
	count := 0
	array := arr[len(arr)-lineIndex:]
	for _,v := range array {
		PrintlnCyan(v)
		if strings.Contains(v,"ERROR")&&strings.Contains(v,"Call got Error") {
			if strings.Contains(v, "transport_handler") {
				count++
			}
		}
	}
	return count
}

//INFO handle
func LogHandle(){

}



