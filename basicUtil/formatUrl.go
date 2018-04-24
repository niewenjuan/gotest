package basicUtil

import "fmt"
const (
	Consumer_rest_IP = "http://127.0.0.1:2098"

)

//循环调用url
func CycleInvokeUrl(protocol,num,delaytime string) string {
	return fmt.Sprintf("%s/consumer/cycle/%s/%s/%s",Consumer_rest_IP,protocol,num,delaytime)
}

//循环调用url
func CycleInvokeQpsUrl(protocol,num,delaytime string) string {
	return fmt.Sprintf("%s/consumer/qps/%s/%s/%s",Consumer_rest_IP,protocol,num,delaytime)
}

//循环跨应用调用qps
func CycleInvokeCrossAppQpsUrl(protocol,num string) string {
	return fmt.Sprintf("%s/consumer/qps/crossapp/%s/%s",Consumer_rest_IP,protocol,num)
}

//权值
func WeightInvokeUrl(protocol,num string) string {
	return fmt.Sprintf("%s/consumer/weight/%s/%s",Consumer_rest_IP,protocol,num)
}

//会话
func SessionInvokeUrl(protocol,num string) string {
	return fmt.Sprintf("%s/consumer/session/%s/%s/%s",Consumer_rest_IP,protocol,num,"0")
}

//rest服务简单调用
func ToRestInvoke(delaytime,protocol string) string{
	return fmt.Sprintf("%s/consumer/rest/simple/%s",Consumer_rest_IP,delaytime)
}

//Highway服务简单调用
func ToHighwayInvoke(delaytime string) string{
	return fmt.Sprintf("%s/consumer/highway/simple/%s",Consumer_rest_IP,delaytime)
}

//跨APP服务rest简单调用
func ToRestInvokeCrossApp(serviceName,appId,delaytime,protocol string) string{
	return fmt.Sprintf("%s/consumer/rest/crossapp/%s/%s/%s/%s",Consumer_rest_IP,serviceName,appId,delaytime,protocol)
}

//跨APP服务rest简单调用
func ToHighwayInvokeCrossApp(serviceName,appId,schemaid,operationid,delaytime string) string{
	return fmt.Sprintf("%s/consumer/highway/crossapp/%s/%s/%s/%s/%s",Consumer_rest_IP,serviceName,appId,schemaid,operationid,delaytime)
}

//自动熔断验证
func ToInvokeAutoCircuit(protocol string)string{
	return fmt.Sprintf("%s/consumer/autocircuit/%s",Consumer_rest_IP,protocol)
}

//容错调用失败实例
func ToInvokeFailInstance(protocol,servicename string)string{
	return fmt.Sprintf("%s/consumer/circuitfail/%s/%s",Consumer_rest_IP,protocol,servicename)
}
