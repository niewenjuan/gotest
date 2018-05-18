package servicecenter

import (
	"fmt"

	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"

	client "github.com/ServiceComb/go-sc-client"
	"github.com/ServiceComb/go-sc-client/model"
)

const (
	// ServiceCenter constant string
	ServiceCenter = "servicecenter"
)

// Registrator to represent the object of service center to call the APIs of service center
type Registrator struct {
	Name           string
	registryClient *client.RegistryClient
	opts           client.Options
}

// RegisterService : 注册微服务
func (r *Registrator) RegisterService(ms *registry.MicroService) (string, error) {
	serviceKey := registry.Microservice2ServiceKeyStr(ms)
	microservice := ToSCService(ms)
	sid, err := r.registryClient.GetMicroServiceID(microservice.AppID, microservice.ServiceName, microservice.Version, microservice.Environment)
	if err != nil {
		lager.Logger.Warnf("Get microservice [%s] failed", serviceKey, err)
	}
	if sid == "" {
		lager.Logger.Warnf("Microservice [%s] not exists in registry, register it", serviceKey, err)
		sid, err = r.registryClient.RegisterService(microservice)
		if err != nil {
			lager.Logger.Errorf(err, "Register microservice [%s] failed", serviceKey)
			return "", err
		}
	} else {
		lager.Logger.Warnf("Microservice [%s] exists in registry, no need register", serviceKey, err)
	}

	return sid, nil
}

// RegisterServiceInstance : 注册微服务
func (r *Registrator) RegisterServiceInstance(sid string, cIns *registry.MicroServiceInstance) (string, error) {
	instance := ToSCInstance(cIns)
	instance.ServiceID = sid
	instanceID, err := r.registryClient.RegisterMicroServiceInstance(instance)
	if err != nil {
		lager.Logger.Errorf(err, "RegisterMicroServiceInstance failed.")
		return "", err
	}
	value, ok := registry.SelfInstancesCache.Get(instance.ServiceID)
	if !ok {
		lager.Logger.Warnf("RegisterMicroServiceInstance get SelfInstancesCache failed, Mid/Sid: %s/%s", instance.ServiceID, instanceID)
	}
	instanceIDs, ok := value.([]string)
	if !ok {
		lager.Logger.Warnf("RegisterMicroServiceInstance type asserts failed,  Mid/Sid: %s/%s", instance.ServiceID, instanceID)
	}
	var isRepeat bool
	for _, va := range instanceIDs {
		if va == instanceID {
			isRepeat = true
		}
	}
	if !isRepeat {
		instanceIDs = append(instanceIDs, instanceID)
	}
	registry.SelfInstancesCache.Set(instance.ServiceID, instanceIDs, 0)
	lager.Logger.Infof("RegisterMicroServiceInstance success, MicroServiceID: %s", instance.ServiceID)

	if instance.HealthCheck == nil ||
		instance.HealthCheck.Mode == model.CheckByHeartbeat {
		registry.HBService.AddTask(sid, instanceID)
	}
	lager.Logger.Infof("RegisterMicroServiceInstance success, microServiceID/instanceID: %s/%s.", sid, instanceID)
	return instanceID, nil
}

// RegisterServiceAndInstance : 注册微服务
func (r *Registrator) RegisterServiceAndInstance(cMicroService *registry.MicroService, cInstance *registry.MicroServiceInstance) (string, string, error) {
	microService := ToSCService(cMicroService)
	instance := ToSCInstance(cInstance)
	microServiceID, err := r.registryClient.GetMicroServiceID(microService.AppID, microService.ServiceName, microService.Version, microService.Environment)
	if microServiceID == "" {
		microServiceID, err = r.registryClient.RegisterService(microService)
		if err != nil {
			lager.Logger.Errorf(err, "RegisterMicroService failed")
			return "", "", err
		}
		lager.Logger.Debugf("RegisterMicroService success, microServiceID: %s", microServiceID)
	}
	instance.ServiceID = microServiceID
	instanceID, err := r.registryClient.RegisterMicroServiceInstance(instance)
	if err != nil {
		lager.Logger.Errorf(err, "RegisterMicroServiceInstance failed.")
		return microServiceID, "", err
	}

	value, ok := registry.SelfInstancesCache.Get(instance.ServiceID)
	if !ok {
		lager.Logger.Warnf("RegisterMicroServiceInstance get SelfInstancesCache failed, Mid/Sid: %s/%s", instance.ServiceID, instanceID)
	}
	instanceIDs, ok := value.([]string)
	if !ok {
		lager.Logger.Warnf("RegisterMicroServiceInstance type asserts failed,  Mid/Sid: %s/%s", instance.ServiceID, instanceID)
	}
	var isRepeat bool
	for _, va := range instanceIDs {
		if va == instanceID {
			isRepeat = true
		}
	}
	if !isRepeat {
		instanceIDs = append(instanceIDs, instanceID)
	}
	registry.SelfInstancesCache.Set(instance.ServiceID, instanceIDs, 0)
	lager.Logger.Infof("RegisterMicroServiceInstance success, MicroServiceID: %s", instance.ServiceID)

	if instance.HealthCheck == nil ||
		instance.HealthCheck.Mode == model.CheckByHeartbeat {
		registry.HBService.AddTask(microServiceID, instanceID)
	}
	lager.Logger.Infof("RegisterMicroServiceInstance success, microServiceID/instanceID: %s/%s.", microServiceID, instanceID)
	return microServiceID, instanceID, nil
}

// UnRegisterMicroServiceInstance : 去注册微服务实例
func (r *Registrator) UnRegisterMicroServiceInstance(microServiceID, microServiceInstanceID string) error {
	isSuccess, err := r.registryClient.UnregisterMicroServiceInstance(microServiceID, microServiceInstanceID)
	if !isSuccess || err != nil {
		lager.Logger.Errorf(nil, "unregisterMicroServiceInstance failed, microServiceID/instanceID = %s/%s.", microServiceID, microServiceInstanceID)
		return err
	}

	value, ok := registry.SelfInstancesCache.Get(microServiceID)
	if !ok {
		lager.Logger.Warnf("UnregisterMicroServiceInstance get SelfInstancesCache failed, Mid/Sid: %s/%s", microServiceID, microServiceInstanceID)
	}
	instanceIDs, ok := value.([]string)
	if !ok {
		lager.Logger.Warnf("UnregisterMicroServiceInstance type asserts failed, Mid/Sid: %s/%s", microServiceID, microServiceInstanceID)
	}
	var newInstanceIDs = make([]string, 0)
	for _, v := range instanceIDs {
		if v != microServiceInstanceID {
			newInstanceIDs = append(newInstanceIDs, v)
		}
	}
	registry.SelfInstancesCache.Set(microServiceID, newInstanceIDs, 0)

	lager.Logger.Debugf("unregisterMicroServiceInstance success, microServiceID/instanceID = %s/%s.", microServiceID, microServiceInstanceID)
	return nil
}

// Heartbeat : Keep instance heartbeats.
func (r *Registrator) Heartbeat(microServiceID, microServiceInstanceID string) (bool, error) {
	bo, err := r.registryClient.Heartbeat(microServiceID, microServiceInstanceID)
	if err != nil {
		lager.Logger.Errorf(err, "Heartbeat failed, microServiceID/instanceID: %s/%s.", microServiceID, microServiceInstanceID)
		return false, err
	}
	if bo == false {
		lager.Logger.Errorf(err, "Heartbeat failed, microServiceID/instanceID: %s/%s.", microServiceID, microServiceInstanceID)
		return bo, err
	}
	lager.Logger.Debugf("Heartbeat success, microServiceID/instanceID: %s/%s.", microServiceID, microServiceInstanceID)
	return bo, nil
}

// AddDependencies ： 注册微服务的依赖关系
func (r *Registrator) AddDependencies(cDep *registry.MicroServiceDependency) error {
	request := ToSCDependency(cDep)
	err := r.registryClient.AddDependencies(request)
	if err != nil {
		lager.Logger.Errorf(err, "AddDependencies failed.")
		return err
	}
	lager.Logger.Debugf("AddDependencies success.")
	return nil
}

// AddSchemas to service center
func (r *Registrator) AddSchemas(microServiceID, schemaName, schemaInfo string) error {
	if err := r.registryClient.AddSchemas(microServiceID, schemaName, schemaInfo); err != nil {
		lager.Logger.Errorf(err, "AddSchemas failed.")
		return err
	}
	lager.Logger.Debugf("AddSchemas success.")
	return nil
}

// UpdateMicroServiceInstanceStatus : 更新微服务实例状态信息
func (r *Registrator) UpdateMicroServiceInstanceStatus(microServiceID, microServiceInstanceID, status string) error {
	isSuccess, err := r.registryClient.UpdateMicroServiceInstanceStatus(microServiceID, microServiceInstanceID, status)
	if !isSuccess {
		lager.Logger.Errorf(nil, "UpdateMicroServiceInstanceStatus failed, microServiceID/instanceID = %s/%s.", microServiceID, microServiceInstanceID)
		return err
	}
	lager.Logger.Debugf("UpdateMicroServiceInstanceStatus success, microServiceID/instanceID = %s/%s.", microServiceID, microServiceInstanceID)
	return nil
}

// UpdateMicroServiceProperties 更新微服务properties信息
func (r *Registrator) UpdateMicroServiceProperties(microServiceID string, properties map[string]string) error {
	microService := &model.MicroService{
		Properties: properties,
	}
	isSuccess, err := r.registryClient.UpdateMicroServiceProperties(microServiceID, microService)
	if !isSuccess {
		lager.Logger.Errorf(nil, "UpdateMicroService Properties failed, microServiceID/instanceID = %s.", microServiceID)
		return err
	}
	lager.Logger.Debugf("UpdateMicroService Properties success, microServiceID/instanceID = %s.", microServiceID)
	return nil
}

// UpdateMicroServiceInstanceProperties : 更新微服务实例properties信息
func (r *Registrator) UpdateMicroServiceInstanceProperties(microServiceID, microServiceInstanceID string, properties map[string]string) error {
	microServiceInstance := &model.MicroServiceInstance{
		Properties: properties,
	}
	isSuccess, err := r.registryClient.UpdateMicroServiceInstanceProperties(microServiceID, microServiceInstanceID, microServiceInstance)
	if !isSuccess {
		lager.Logger.Errorf(nil, "UpdateMicroServiceInstanceProperties failed, microServiceID/instanceID = %s/%s.", microServiceID, microServiceInstanceID)
		return err
	}
	lager.Logger.Debugf("UpdateMicroServiceInstanceProperties success, microServiceID/instanceID = %s/%s.", microServiceID, microServiceInstanceID)
	return nil
}

// Close : Close all client connection.
func (r *Registrator) Close() error {
	return closeClient(r.registryClient)
}

// ServiceDiscovery to represent the object of service center to call the APIs of service center
type ServiceDiscovery struct {
	Name           string
	registryClient *client.RegistryClient
	opts           client.Options
}

// GetMicroServiceID : 获取指定微服务的MicroServiceID
func (r *ServiceDiscovery) GetMicroServiceID(appID, microServiceName, version, env string) (string, error) {
	microServiceID, err := r.registryClient.GetMicroServiceID(appID, microServiceName, version, env)
	if err != nil {
		lager.Logger.Errorf(err, "GetMicroServiceID failed")
		return "", err
	}
	lager.Logger.Debugf("GetMicroServiceID success")
	return microServiceID, nil
}

// GetAllMicroServices : Get all MicroService information.
func (r *ServiceDiscovery) GetAllMicroServices() ([]*registry.MicroService, error) {
	microServices, err := r.registryClient.GetAllMicroServices()
	if err != nil {
		lager.Logger.Errorf(err, "GetAllMicroServices failed")
		return nil, err
	}
	mss := []*registry.MicroService{}
	for _, s := range microServices {
		mss = append(mss, ToMicroService(s))
	}
	lager.Logger.Debugf("GetAllMicroServices success, MicroService: %s", microServices)
	return mss, nil
}

// GetAllApplications : Get all Applications information.
func (r *ServiceDiscovery) GetAllApplications() ([]string, error) {
	apps, err := r.registryClient.GetAllApplications()
	if err != nil {
		lager.Logger.Errorf(err, "GetAllApplications failed")
		return nil, err
	}
	appArray := []string{}
	for _, s := range apps {
		appArray = append(appArray, s)
	}
	lager.Logger.Debugf("GetAllApplications success, Applications: %s", apps)
	return appArray, nil
}

// GetMicroService : 根据microServiceID获取对应的微服务信息
func (r *ServiceDiscovery) GetMicroService(microServiceID string) (*registry.MicroService, error) {
	microService, err := r.registryClient.GetMicroService(microServiceID)
	if err != nil {
		lager.Logger.Errorf(err, "GetMicroService failed")
		return nil, err
	}
	lager.Logger.Debugf("GetMicroServices success, MicroService: %s", microService)
	return ToMicroService(microService), nil
}

// GetMicroServiceInstances : 获取指定微服务的所有实例
func (r *ServiceDiscovery) GetMicroServiceInstances(consumerID, providerID string) ([]*registry.MicroServiceInstance, error) {
	providerInstances, err := r.registryClient.GetMicroServiceInstances(consumerID, providerID)
	if err != nil {
		lager.Logger.Errorf(err, "GetMicroServiceInstances failed.")
		return nil, err
	}
	instances := filterInstances(providerInstances)
	lager.Logger.Debugf("GetMicroServiceInstances success, consumerID/providerID: %s/%s", consumerID, providerID)
	return instances, nil
}

// FindMicroServiceInstances find micro-service instances
func (r *ServiceDiscovery) FindMicroServiceInstances(consumerID, appID, microServiceName, version, env string) ([]*registry.MicroServiceInstance, error) {
	key := microServiceName + ":" + version + ":" + appID

	value, boo := registry.MicroserviceInstanceCache.Get(key)
	if !boo || value == nil {
		lager.Logger.Warnf("%s Get instances from remote, key: %s", consumerID, key)
		providerInstances, err := r.registryClient.FindMicroServiceInstances(consumerID, appID, microServiceName,
			findVersionRule(microServiceName))
		if err != nil {
			return nil, fmt.Errorf("FindMicroServiceInstances failed, ProviderID: %s, err: %s", key, err)
		}

		filterRestore(providerInstances, microServiceName, appID)
		value, boo = registry.MicroserviceInstanceCache.Get(key)
		if !boo || value == nil {
			lager.Logger.Debugf("Find no microservice instances for %s from cache", key)
			return nil, nil
		}
	}
	microServiceInstance, ok := value.([]*registry.MicroServiceInstance)
	if !ok {
		lager.Logger.Errorf(nil, "FindMicroServiceInstances failed, Type asserts failed.consumerIDL: %s", consumerID)
	}
	return microServiceInstance, nil
}

// GetDependentMicroServiceInstances : 获取指定微服务所依赖的所有实例
func (r *ServiceDiscovery) GetDependentMicroServiceInstances(appID, consumerMicroServiceName, version, env string) ([]*model.MicroServiceInstance, error) {
	var instancesAll []*model.MicroServiceInstance
	microServiceConsumerID, err := r.GetMicroServiceID(appID, consumerMicroServiceName, version, env)
	if err != nil {
		lager.Logger.Errorf(err, "GetMicroServiceID failed.")
		return nil, err
	}
	providers, err := r.registryClient.GetProviders(microServiceConsumerID)
	if err != nil {
		lager.Logger.Errorf(err, "Get Provider failed.")
		return nil, err
	}
	for _, provider := range providers.Services {
		microServiceProviderID, err := r.GetMicroServiceID(provider.AppID, provider.ServiceName, provider.Version, env)
		if err != nil {
			lager.Logger.Errorf(err, "GetMicroServiceID failed.")
			return nil, err
		}
		instances, err := r.GetMicroServiceInstances(microServiceConsumerID, microServiceProviderID)
		if err != nil {
			lager.Logger.Errorf(err, "GetMicroServiceInstances failed.")
			return nil, err
		}
		for _, value := range instances {
			instancesAll = append(instancesAll, ToSCInstance(value))
		}
	}
	lager.Logger.Debugf("GetDependentMicroServiceInstances success, appID/microServiceName/version: %s/%s/%s", appID, consumerMicroServiceName, version)
	return instancesAll, nil
}

// WatchMicroService : 支持用户自调用主动监听实例变化功能
func (r *ServiceDiscovery) WatchMicroService(selfMicroServiceID string, callback func(*model.MicroServiceInstanceChangedEvent)) {
	r.registryClient.WatchMicroService(selfMicroServiceID, callback)
}

// AutoSync updating the cache manager
func (r *ServiceDiscovery) AutoSync() {
	c := &CacheManager{
		registryClient: r.registryClient,
	}
	c.AutoSync()
}

// Close : Close all websocket connection.
func (r *ServiceDiscovery) Close() error {
	return closeClient(r.registryClient)
}

// ContractDiscovery to represent the object of service center to call the APIs of service center
type ContractDiscovery struct {
	Name           string
	registryClient *client.RegistryClient
	opts           client.Options
}

// GetMicroServicesByInterface get micro-services by interface
func (r *ContractDiscovery) GetMicroServicesByInterface(interfaceName string) (microService []*registry.MicroService) {
	value, ok := registry.SchemaInterfaceIndexedCache.Get(interfaceName)
	if !ok || value == nil {
		r.fillSchemaInterfaceIndexCache(nil, interfaceName)
		value, _ = registry.SchemaInterfaceIndexedCache.Get(interfaceName)
	}

	microServiceModel, ok := value.([]*model.MicroService)

	if !ok {
		lager.Logger.Errorf(nil, "GetMicroServicesByInterface failed, Type asserts failed")
	}

	for _, v := range microServiceModel {
		ms := ToMicroService(v)
		microService = append(microService, ms)
	}
	return microService
}

// GetSchemaContentByInterface get schema content by interface
func (r *ContractDiscovery) GetSchemaContentByInterface(interfaceName string) (schemas registry.SchemaContent) {
	value, ok := registry.SchemaInterfaceIndexedCache.Get(interfaceName)
	if !ok || value == nil {
		return r.fillSchemaInterfaceIndexCache(nil, interfaceName)
	}

	val, ok := value.([]*model.MicroService)

	if !ok {
		return schemas
	}

	return r.fillSchemaInterfaceIndexCache(val, interfaceName)

}

// GetSchemaContentByServiceName get schema content by service name
func (r *ContractDiscovery) GetSchemaContentByServiceName(svcName, version, appID, env string) (schemas []*registry.SchemaContent) {
	serviceID, err := r.registryClient.GetMicroServiceID(appID, svcName, version, env)
	if err != nil {
		return schemas
	}
	value, ok := registry.SchemaServiceIndexedCache.Get(serviceID)
	if !ok || value == nil {
		return r.fillSchemaServiceIndexCache(nil, serviceID)
	}

	val, ok := value.([]*model.MicroService)

	if !ok {
		return schemas
	}

	return r.fillSchemaServiceIndexCache(val, serviceID)

}

// fillSchemaServiceIndexCache fill schema service index cache
func (r *ContractDiscovery) fillSchemaServiceIndexCache(ms []*model.MicroService, serviceID string) (content []*registry.SchemaContent) {
	if ms == nil {
		microServiceList, err := r.registryClient.GetAllMicroServices()
		if err != nil {
			lager.Logger.Errorf(err, "Get instances failed")
			return content
		}

		return r.fillCacheAndGetServiceSchemaContent(microServiceList, serviceID)
	}

	return r.fillCacheAndGetServiceSchemaContent(ms, serviceID)
}

// fillCacheAndGetServiceSchemaContent fill cache and get services schema content
func (r *ContractDiscovery) fillCacheAndGetServiceSchemaContent(microServiceList []*model.MicroService, serviceID string) (schemaContent []*registry.SchemaContent) {

	for _, ms := range microServiceList {

		if ms.ServiceID == serviceID {

			for _, schemaName := range ms.Schemas {
				content, err := r.registryClient.GetSchema(ms.ServiceID, schemaName)
				if err != nil {
					continue
				}

				schema, err := unmarshalSchemaContent(content)
				if err != nil {
					continue
				}
				_, ok := registry.SchemaServiceIndexedCache.Get(serviceID)
				if !ok {
					var allServices []*model.MicroService
					allServices = append(allServices, ms)
					registry.SchemaServiceIndexedCache.Set(serviceID, allServices, 0)
				}
				schemaContent = append(schemaContent, schema)
			}
		}
	}

	return
}

// fillSchemaInterfaceIndexCache fill schema interface index cache
func (r *ContractDiscovery) fillSchemaInterfaceIndexCache(ms []*model.MicroService, interfaceName string) (content registry.SchemaContent) {
	if ms == nil {
		microServiceList, err := r.registryClient.GetAllMicroServices()
		if err != nil {
			lager.Logger.Errorf(err, "Get instances failed")
			return content
		}

		return r.fillCacheAndGetInterfaceSchemaContent(microServiceList, interfaceName)
	}

	return r.fillCacheAndGetInterfaceSchemaContent(ms, interfaceName)
}

// fillCacheAndGetInterfaceSchemaContent fill cache and get interface schema content
func (r *ContractDiscovery) fillCacheAndGetInterfaceSchemaContent(microServiceList []*model.MicroService, interfaceName string) (schemaContent registry.SchemaContent) {

	for _, ms := range microServiceList {
		serviceID, err := r.registryClient.GetMicroServiceID(ms.AppID, ms.ServiceName, ms.Version, ms.Environment)
		if err != nil {
			continue
		}

		for _, schemaName := range ms.Schemas {
			content, err := r.registryClient.GetSchema(serviceID, schemaName)
			if err != nil {
				continue
			}

			schemaContent, err = parseSchemaContent(content)
			if err != nil {
				continue
			}

			interfaceValue := schemaContent.Info["x-java-interface"]
			if interfaceValue == "" {
				continue
			}

			value, ok := registry.SchemaInterfaceIndexedCache.Get(interfaceName)
			if !ok {
				var allServices []*model.MicroService
				allServices = append(allServices, ms)
				registry.SchemaInterfaceIndexedCache.Set(interfaceValue, allServices, 0)
			} else {
				val, _ := value.([]*model.MicroService)
				val = append(val, ms)
				registry.SchemaInterfaceIndexedCache.Set(interfaceValue, val, 0)

			}

			if interfaceName == interfaceValue {
				return schemaContent
			}
		}
	}

	return
}

// GetSchema from service center
func (r *ContractDiscovery) GetSchema(microServiceID, schemaName string) ([]byte, error) {
	var schemaContent []byte
	var err error
	if schemaContent, err = r.registryClient.GetSchema(microServiceID, schemaName); err != nil {
		lager.Logger.Errorf(err, "GetSchema failed.")
		return []byte(""), err
	}
	lager.Logger.Debugf("GetSchema success.")
	return schemaContent, nil

}

//Close close client connection
func (r *ContractDiscovery) Close() error {
	return closeClient(r.registryClient)
}
func newRegistrator(options registry.Options) registry.Registrator {
	sco := ToSCOptions(options)
	r := &client.RegistryClient{}
	if err := r.Initialize(sco); err != nil {
		lager.Logger.Errorf(err, "RegistryClient initialization failed.")
	}

	return &Registrator{
		Name:           ServiceCenter,
		registryClient: r,
		opts:           sco,
	}
}
func newServiceDiscovery(options registry.Options) registry.ServiceDiscovery {
	sco := ToSCOptions(options)
	r := &client.RegistryClient{}
	if err := r.Initialize(sco); err != nil {
		lager.Logger.Errorf(err, "RegistryClient initialization failed.")
	}

	return &ServiceDiscovery{
		Name:           ServiceCenter,
		registryClient: r,
		opts:           sco,
	}
}
func newContractDiscovery(options registry.Options) registry.ContractDiscovery {
	sco := ToSCOptions(options)
	r := &client.RegistryClient{}
	if err := r.Initialize(sco); err != nil {
		lager.Logger.Errorf(err, "RegistryClient initialization failed.")
	}

	return &ContractDiscovery{
		Name:           ServiceCenter,
		registryClient: r,
		opts:           sco,
	}
}

// init initialize the plugin of service center registry
func init() {
	registry.InstallRegistrator(ServiceCenter, newRegistrator)
	registry.InstallServiceDiscovery(ServiceCenter, newServiceDiscovery)
	registry.InstallContractDiscovery(ServiceCenter, newContractDiscovery)

}
