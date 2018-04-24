package basicUtil

import "fmt"

const(
	AppId = "default"
	Version = "1.0.1"
	ServiceName = "ClientHoy"
	ServerName = "GoServerHoy"
	MidServerName = "GoMidServer"
	ClientName = "ClientHoy"
	OperationId_highwaySimpleInvoke = "LoadBalanceTest"
	AppId_highwaySimpleInvoke = "go"
	SchemaId_highwaySimpleInvoke = "hellworld"
	REST = "rest"
	HIGHWAY = "highway"

	//配置项key
	//负载均衡
	Key_strategy = "cse.loadbalance.strategy.name"
	Key_service_strategy = "cse.loadbalance.%s.strategy.name"
	//灰度发布（路由规则）
	Key_darklaunch_policy = "cse.darklaunch.policy.%s"
	//限流
	Key_qps_global_enable = "cse.flowcontrol.Consumer.qps.enabled"
	Key_qps_global_limit = "cse.flowcontrol.Consumer.qps.global.limit"
	Key_qps_service_limit = "cse.flowcontrol.Consumer.qps.limit.%s"
	Key_qps_provider_global_enable = "cse.flowcontrol.Provider.qps.enabled"
	Key_qps_provider_global_limit = "cse.flowcontrol.Provider.qps.global.limit"
	Key_qps_provider_service_limit = "cse.flowcontrol.Provider.qps.limit.%s"
	//隔离
	Key_isolation_global_timeout_able = "cse.isolation.Consumer.timeout.enabled"
	Key_isolation_global_timeoutInMilliseconds = "cse.isolation.Consumer.timeoutInMilliseconds"
	Key_isolation_service_timeout_able = "cse.isolation.Consumer.%s.timeout.enabled"
	Key_isolation_service_timeoutInMilliseconds = "cse.isolation.Consumer.%s.timeoutInMilliseconds"
	Key_isolation_global_maxConcurrentRequests = "cse.isolation.Consumer.maxConcurrentRequests"
	Key_isolation_service_maxConcurrentRequests = "cse.isolation.Consumer.%s.maxConcurrentRequests"
	//熔断
	Key_circuit_consumer_enable = "cse.circuitBreaker.Consumer.enabled"
	Key_circuit_manaul_forceOpen = "cse.circuitBreaker.Consumer.forceOpen"
	Key_circuit_manaul_forceClosed = "cse.circuitBreaker.Consumer.forceClosed"
	Key_circuit_service_consumer_enable = "cse.circuitBreaker.Consumer.%s.enabled"
	Key_circuit_service_manaul_forceOpen = "cse.circuitBreaker.Consumer.%s.forceOpen"
	Key_circuit_consumer_requestVolumeThreshold = "cse.circuitBreaker.Consumer.requestVolumeThreshold"
	Key_circuit_consumer_errorThresholdPercentage = "cse.circuitBreaker.Consumer.errorThresholdPercentage"
	Key_circuit_consumer_sleepWindowInMilliseconds = "cse.circuitBreaker.Consumer.sleepWindowInMilliseconds"
	Key_circuit_service_consumer_requestVolumeThreshold = "cse.circuitBreaker.Consumer.%s.requestVolumeThreshold"
	Key_circuit_service_consumer_errorThresholdPercentage = "cse.circuitBreaker.Consumer.%s.errorThresholdPercentage"
	Key_circuit_service_consumer_sleepWindowInMilliseconds = "cse.circuitBreaker.Consumer.%s.sleepWindowInMilliseconds"
	//容错
	Key_fallback_consumer_enable = "cse.fallback.Consumer.enabled"
	Key_consumer_fallbackpolicy = "cse.fallbackpolicy.Consumer.policy"
	Key_fallback_consumer_service_enable = "cse.fallback.Consumer.%s.enabled"
	Key_consumer_service_fallbackpolicy = "cse.fallbackpolicy.Consumer.%s.policy"
	//降级
	Key_retry_enable = "cse.loadbalance.retryEnabled"
	Key_retry_onsame = "cse.loadbalance.retryOnSame"
	Key_retry_onnext = "cse.loadbalance.retryOnNext"



	//配置项value
	Value_strategy_random = "Random"
	Value_strategy_round = "RoundRobin"
	Value_strategy_session = "SessionStickiness"
	Value_strategy_weight = "WeightedResponse"
	Value_fallbackpolicy_exception = "throwexception"
	Value_fallbackpolicy_returnnull = "returnnull"

	//中心接口请求响应
	Result_success = "{\"Result\": \"Success\"}"

	//logPath
	LogPath = "D:\\APItest\\gitSvn\\csetest\\cse-gosdk-testPiles\\src\\github.com\\cairixian\\log\\client\\chassis.log"

	//strCheck
	FallbackTimeOut = "\"error\":\"timeout\""
	MaxCurrentError = "\"error\":\"max concurrency\""
	CircuitError = "\"error\":\"circuit open\""
)
type GinkgoTest struct {
	ServerName string
}

func (gt *GinkgoTest)Get_Key_service_strategy(serviceName string)string{
	return fmt.Sprintf(Key_service_strategy,serviceName)
}

func (gt *GinkgoTest)Get_Key_server_strategy()string{
	return fmt.Sprintf(Key_service_strategy,ServerName)
}

func (gt *GinkgoTest)Get_Key_darklaunch_policy(scname string)string{
	return fmt.Sprintf(Key_darklaunch_policy,scname)
}

func (gt *GinkgoTest)Get_Key_isolation_service_timeout_able(scname string)string{
	return fmt.Sprintf(Key_isolation_service_timeout_able,scname)
}

func (gt *GinkgoTest)Get_Key_isolation_service_timeoutInMilliseconds(scname string)string{
	return fmt.Sprintf(Key_isolation_service_timeoutInMilliseconds,scname)
}

func (gt *GinkgoTest)Get_Key_circuit_service_consumer_enable(scname string)string{
	return fmt.Sprintf(Key_circuit_service_consumer_enable,scname)
}

func (gt *GinkgoTest)Get_Key_circuit_service_manaul_forceOpen(scname string)string{
	return fmt.Sprintf(Key_circuit_service_manaul_forceOpen,scname)
}

func (gt *GinkgoTest)Get_Key_fallback_consumer_service_enable(scname string)string{
	return fmt.Sprintf(Key_fallback_consumer_service_enable,scname)
}

func (gt *GinkgoTest)Get_Key_consumer_service_fallbackpolicy(scname string)string{
	return fmt.Sprintf(Key_consumer_service_fallbackpolicy,scname)
}

func (gt *GinkgoTest)Get_Key_isolation_service_maxConcurrentRequests(scname string)string{
	return fmt.Sprintf(Key_isolation_service_maxConcurrentRequests,scname)
}

func (gt *GinkgoTest)Get_Key_circuit_service_consumer_requestVolumeThreshold(scname string)string{
	return fmt.Sprintf(Key_circuit_service_consumer_requestVolumeThreshold,scname)
}

func (gt *GinkgoTest)Get_Key_circuit_service_consumer_errorThresholdPercentage(scname string)string{
	return fmt.Sprintf(Key_circuit_service_consumer_errorThresholdPercentage,scname)
}

func (gt *GinkgoTest)Get_Key_circuit_service_consumer_sleepWindowInMilliseconds(scname string)string{
	return fmt.Sprintf(Key_circuit_service_consumer_sleepWindowInMilliseconds,scname)
}

func (gt *GinkgoTest)Get_Key_qps_service_limit(scname string)string{
	return fmt.Sprintf(Key_qps_service_limit,scname)
}

func (gt *GinkgoTest)Get_Key_qps_provider_limit(scname string)string{
	return fmt.Sprintf(Key_qps_provider_service_limit,scname)
}


