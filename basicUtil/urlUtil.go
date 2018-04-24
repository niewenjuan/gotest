package basicUtil

const(

	//服务端
	SimpleInvoke_url = "/simpleInvoke/{user}/{age}"
	SimpleInvoke_operationId = "SimpleInvoke"

	CircuitFailInvoke_url = "/failInvoke"
	CircuitFailInvoke_operationId = "CircuitFailTest"

	Qps_url = "/demo/QpsTest"
	Qps_operationId = "QpsTest"

	LoadBalanceTest_url = "/demo/LoadBalanceTest/{delaytime}"
	LoadBalanceTest_operationId = "LoadBalanceTest"

	LoadBalanceTest_Weight_url = "/demo/LoadBalanceTestWeight"
	LoadBalanceTest_Weight_operationId = "LoadBalanceTest_Weight"

	ConcurrentTest_url = "/demo/ConcurrentTest/delaytime/{delaytime}/number/{number}"
	ConcurrentTest_operationId = "ConcurrentTest"

	//basicUtil
	ANSI_COLOR_LIGHT_RED = "\x1b[1;31m"
	ANSI_COLOR_LIGHT_GREEN = "\x1b[1;32m"
	ANSI_COLOR_LIGHT_YELLOW = "\x1b[1;33m"
	//ANSI_COLOR_LIGHT_BLUE = "\x1b[1;34m"
	ANSI_COLOR_LIGHT_MAGENTA = "\x1b[1;35m"
	ANSI_COLOR_LIGHT_CYAN = "\x1b[1;36m"
	ANSI_COLOR_LIGHT_RESET = "\x1b[0m"

	//获取token信息（XJ租户）
	//Token_DomainName = "CSE_xwx282889"
	//Token_UserName = "CSE_xwx282889"
	//Token_PassWord = "Huawei@123"
	//Token_ProjectName = "southchina"

	//获取token信息（香港租户）
	Token_DomainName = "CSE_xwx282889"
	Token_UserName = "CSE_xwx282889"
	Token_PassWord = "Huawei@123"
	Token_ProjectName = "southchina_01"

	//类生产
	hcUrl = "https://cse.cn-north-1.myhuaweicloud.com:443"
	//香港region
	//hcUrl = "https://cse.cn-hk1.myhwclouds.com:443"

	CCUrl = hcUrl + "/v3/default/configuration/items"
	SCUrl = hcUrl + "/v4/default/registry/microservices"
	SCExistenceUrl = hcUrl + "/v4/default/registry/existence"
	GetTokenUrl = "https://192.144.1.37:31943/v3/auth/tokens"

	//autotest
	CreateConfig = "/createConfig"
	CreateConfig_oprerationId = "CreateConfigItem"

	//客户端
	C_SimpleInvoke_Url = "/consumer/highway/simple/{delaytime}"
	C_SimpleInvoke_OperationId = "SimpleInvokeServer"

	C_CrossAppHighwaySimpleInvoke_Url = "/consumer/highway/crossapp/{serviceName}/{appId}/{schemaId}/{operationId}/{delayTime}"
	C_CrossAppHighwaySimpleInvoke_OperationId = "SimpleInvokeServerCrossApp"

	C_RestSimpleInvoke_Url = "/consumer/rest/simple/{delayTime}"
	C_RestSimpleInvoke_OperationId = "RestSimpleInvokeServer"

	C_CrossAppRestSimpleInvoke_Url = "/consumer/rest/crossapp/{serviceName}/{appId}/{delayTime}/{option}"
	C_CrossAppRestSimpleInvoke_OperationId = "RestSimpleInvokeServerCrossApp"

	C_CrossAppQpsInvoke_Url = "/consumer/qps/crossapp/{protocol}/{num}"
	C_CrossAppQpsInvoke_OperationId = "CycleInvokeServerForQpsCrossApp"

	C_CycleInvokeQps_Url = "/consumer/qps/{protocol}/{num}/{delaytime}"
	C_CycleInvokeQps_OperationId = "CycleInvokeServerForQps"

	C_CycleInvoke_Url = "/consumer/cycle/{protocol}/{num}/{delaytime}"
	C_CycleInvoke_OperationId = "CycleInvokeServer"

	C_CycleInvoke_Weight_Url = "/consumer/weight/{protocol}/{num}"
	C_CycleInvoke_Weight_OperationId = "CycleInvokeServer_weight"

	C_CycleInvoke_Session_Url = "/consumer/session/{protocol}/{num}/{delaytime}"
	C_CycleInvoke_Session_OperationId = "CycleInvokeServer_session"

	C_QueryToken_Url = "/consumer/queryToken"
	C_QueryToken_OperationId = "QueryToken"

	C_AutoCircuit_Url = "/consumer/autocircuit/{procotol}"
	C_AutoCircuit_OperationId = "AutoCircuitRestFulTest"

	C_CircuitFailInvoke_Url = "/consumer/circuitfail/{protocol}/{serviceName}"
	C_CircuitFailInvoke_OperationId = "CircuitFailInvoke"
)


