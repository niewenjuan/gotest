package basicUtil

import (
	"fmt"
	"encoding/json"
)

type ServiceCenter struct {

}

//query service by appid,version and scname
func (sc *ServiceCenter) QueryExistenceFromSC(appid,servicename,version string)[]byte{
	tempUrl := fmt.Sprintf("?type=microservice&appId=%s&serviceName=%s&version=%s",appid,servicename,version)
	url := SCExistenceUrl + tempUrl
	//PrintlnCyan(url)
	return requestDo("GET",url,"")
}

//create Service
func (sc *ServiceCenter) CreateServiceToSc(body map[string]interface{})[]byte{
	type CreateService struct {
		Service         map[string]interface{} `json:"service"`
	}
	createservice := new(CreateService)
	createservice.Service = body    //服务名
	// 配置项
	config, err := json.Marshal(createservice)
	if err != nil {
		fmt.Println("failed to marshal body")
		return nil
	}

	url := SCUrl
	return requestDo("POST",url,string(config))
}

//query instances
func (sc *ServiceCenter) QueryInstanceIdFromSC(serviceId string)[]byte{
	tempUrl := fmt.Sprintf("/%s/instances",serviceId)
	url := SCUrl + tempUrl
	//PrintlnCyan(url)
	return requestDoA("GET",url,"",serviceId)
}

//create instances
func (sc *ServiceCenter) CreateInstanceIdFromSC(serviceId string,body map[string]interface{})[]byte{
	type CreateInstance struct {
		Instance         map[string]interface{} `json:"instance"`
	}
	createinstance := new(CreateInstance)
	createinstance.Instance = body    //服务名
	// 配置项
	config, err := json.Marshal(createinstance)
	if err != nil {
		fmt.Println("failed to marshal body")
		return nil
	}

	tempUrl := fmt.Sprintf("/%s/instances",serviceId)
	url := SCUrl + tempUrl


	//PrintlnCyan(url)
	return requestDo("POST",url,string(config))
}

//delete service
func (sc *ServiceCenter) DeleteServiceFromSC(serviceId string)[]byte{
	url := SCUrl + "/"+serviceId+"?force=1"
	return requestDo("DELETE",url,"")
}

//query all service
func (sc *ServiceCenter) QueryAllServicesFromSC()[]byte{
	//PrintlnCyan(url)
	return requestDo("GET",SCUrl,"")
}

