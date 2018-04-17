package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/config/schema"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/util/fileutil"

	"gopkg.in/yaml.v2"
)

// GlobalDefinition is having the information about region, load balancing, service center, config center,
// protocols, and handlers for the micro service
var GlobalDefinition *model.GlobalCfg
var lbConfig *model.LBWrapper

// MicroserviceDefinition is having the info about application id, provider info, description of the service,
// and description of the instance
var MicroserviceDefinition *model.MicroserviceCfg

// PassLagerDefinition is having the information about loging
var PassLagerDefinition *model.PassLagerCfg

//HystrixConfig is having info about isolation, circuit breaker, fallback properities of the micro service
var HystrixConfig *model.HystrixConfigWrapper

// NodeIP gives the information of node ip
var NodeIP string

// SelfServiceID 单进程多微服务根本没法记录依赖关系，因为一个进程里有多个微服务，你在调用别的微服务时到底该怎么添加依赖关系？
//只能随意赋值个id
var SelfServiceID string

// SelfServiceName is self micro service name
var SelfServiceName string

// SelfMetadata is gives meta data of the self micro service
var SelfMetadata map[string]string

// SelfVersion gives version of the self micro service
var SelfVersion string

// ErrNoName is used to represent the service name missing error
var ErrNoName = errors.New("Microservice name is missing in description file")

// parse unmarshal configurations on respective structure
func parse() error {
	err := readGlobalConfigFile()
	if err != nil {
		return err
	}
	err = ReadLBFromArchaius()
	if err != nil {
		return err
	}

	err = ReadHystrixFromArchaius()
	if err != nil {
		return err
	}

	err = readMicroserviceConfigFiles()
	if err != nil {
		return err
	}

	populateConfigCenterAddress()
	populateServiceRegistryAddress()
	populateMonitorServerAddress()
	populateServiceEnvironment()
	populateServiceName()
	populateVersion()
	return nil
}

// populateServiceRegistryAddress populate service registry address
func populateServiceRegistryAddress() {
	//Registry Address , higher priority for environment variable
	registryAddrFromEnv := archaius.GetString(common.CseRegistryAddress, "")
	if registryAddrFromEnv != "" {
		GlobalDefinition.Cse.Service.Registry.Address = registryAddrFromEnv
	}
}

// populateConfigCenterAddress populate config center address
func populateConfigCenterAddress() {
	//Config Center Address , higher priority for environment variable
	configCenterAddrFromEnv := archaius.GetString(common.CseConfigCenterAddress, "")
	if configCenterAddrFromEnv != "" {
		GlobalDefinition.Cse.Config.Client.ServerURI = configCenterAddrFromEnv
	}
}

// populateMonitorServerAddress populate monitor server address
func populateMonitorServerAddress() {
	//Monitor Center Address , higher priority for environment variable
	monitorServerAddrFromEnv := archaius.GetString(common.CseMonitorServer, "")
	if monitorServerAddrFromEnv != "" {
		GlobalDefinition.Cse.Monitor.Client.ServerURI = monitorServerAddrFromEnv
	}
}

// populateServiceEnvironment populate service environment
func populateServiceEnvironment() {
	if e := archaius.GetString(common.Env, ""); e != "" {
		MicroserviceDefinition.ServiceDescription.Environment = e
	}
}

// populateServiceName populate service name
func populateServiceName() {
	if e := archaius.GetString(common.ServiceName, ""); e != "" {
		MicroserviceDefinition.ServiceDescription.Name = e
	}
}

// populateVersion populate version
func populateVersion() {
	if e := archaius.GetString(common.Version, ""); e != "" {
		MicroserviceDefinition.ServiceDescription.Version = e
	}
}

// readGlobalConfigFile for to unmarshal the global config file(chassis.yaml) information
func readGlobalConfigFile() error {
	globalDef := model.GlobalCfg{}
	err := archaius.UnmarshalConfig(&globalDef)
	if err != nil {
		return err
	}
	GlobalDefinition = &globalDef

	return nil
}

// ReadLBFromArchaius for to unmarshal the global config file(chassis.yaml) information
func ReadLBFromArchaius() error {
	lbMutex.Lock()
	defer lbMutex.Unlock()
	lbDef := model.LBWrapper{}
	err := archaius.UnmarshalConfig(&lbDef)
	if err != nil {
		return err
	}
	lbConfig = &lbDef

	return nil
}

// readPassLagerConfigFile is unmarshal the paas lager configuration file(lager.yaml)
func readPassLagerConfigFile(lagerFile string) error {
	passLagerDef := model.PassLagerCfg{}
	yamlFile, err := ioutil.ReadFile(lagerFile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &passLagerDef)
	if err != nil {
		log.Printf("Unmarshal: %v", err)
	}
	PassLagerDefinition = &passLagerDef

	return nil
}

// ReadHystrixFromArchaius is unmarshal hystrix configuration file(circuit_breaker.yaml)
func ReadHystrixFromArchaius() error {
	cbMutex.RLock()
	defer cbMutex.RUnlock()
	hystrixCnf := model.HystrixConfigWrapper{}
	err := archaius.UnmarshalConfig(&hystrixCnf)
	if err != nil {
		return err
	}
	HystrixConfig = &hystrixCnf
	return nil
}

// readMicroserviceConfigFiles read micro service configuration file
func readMicroserviceConfigFiles() error {
	MicroserviceDefinition = &model.MicroserviceCfg{}
	//find only one microservice yaml
	microserviceNames := schema.GetMicroserviceNames()
	defPath := fileutil.GetMicroserviceDesc()
	data, err := ioutil.ReadFile(defPath)
	if err != nil {
		lager.Logger.Errorf(err, fmt.Sprintf("WARN: Missing microservice description file: %s", err.Error()))
		if len(microserviceNames) == 0 {
			return errors.New("Missing microservice description file")
		}
		msName := microserviceNames[0]
		msDefPath := fileutil.MicroserviceDefinition(msName)
		lager.Logger.Warnf(nil, fmt.Sprintf("Try to find microservice description file in [%s]", msDefPath))
		data, err := ioutil.ReadFile(msDefPath)
		if err != nil {
			return fmt.Errorf("Missing microservice description file: %s", err.Error())
		}
		ReadMicroserviceConfigFromBytes(data)
		return nil
	}
	return ReadMicroserviceConfigFromBytes(data)
}

// ReadMicroserviceConfigFromBytes read micro service configurations from bytes
func ReadMicroserviceConfigFromBytes(data []byte) error {
	microserviceDef := model.MicroserviceCfg{}
	err := yaml.Unmarshal([]byte(data), &microserviceDef)
	if err != nil {
		return err
	}
	if microserviceDef.ServiceDescription.Name == "" {
		return ErrNoName
	}
	if microserviceDef.ServiceDescription.Version == "" {
		microserviceDef.ServiceDescription.Version = common.DefaultVersion
	}

	MicroserviceDefinition = &microserviceDef
	return nil
}

//GetLoadBalancing return lb config
func GetLoadBalancing() *model.LoadBalancing {
	if lbConfig != nil {
		return lbConfig.Prefix.LBConfig
	}
	return nil
}

//GetHystrixConfig return cb config
func GetHystrixConfig() *model.HystrixConfig {
	return HystrixConfig.HystrixConfig
}

// Init is initialize the configuration directory, lager, archaius, route rule, and schema
func Init() error {

	lagerFile := fileutil.PassLagerDefinition()

	err := readPassLagerConfigFile(lagerFile)
	if err != nil {
		log.Println("WARN:lager.yaml does not exist,use the default configuration")
	}

	lager.Initialize(PassLagerDefinition.Writers, PassLagerDefinition.LoggerLevel,
		PassLagerDefinition.LoggerFile, PassLagerDefinition.RollingPolicy,
		PassLagerDefinition.LogFormatText, PassLagerDefinition.LogRotateDate,
		PassLagerDefinition.LogRotateSize, PassLagerDefinition.LogBackupCount)

	err = archaius.Init()
	if err != nil {
		return err
	}
	lager.Logger.Infof("archaius init success")

	var schemaError error

	//Upload schemas using environment variable SCHEMA_ROOT
	schemaEnv := archaius.GetString(common.EnvSchemaRoot, "")
	if schemaEnv != "" {
		schemaError = schema.LoadSchema(schemaEnv, true)
	} else {
		schemaError = schema.LoadSchema(fileutil.GetConfDir(), false)
	}

	if schemaError != nil {
		return schemaError
	}

	//set microservice names
	msError := schema.SetMicroServiceNames(fileutil.GetConfDir())
	if msError != nil {
		return msError
	}

	NodeIP = archaius.GetString(common.EnvNodeIP, "")
	err = parse()
	if err != nil {
		return err
	}

	SelfServiceName = MicroserviceDefinition.ServiceDescription.Name
	SelfVersion = MicroserviceDefinition.ServiceDescription.Version
	SelfMetadata = MicroserviceDefinition.ServiceDescription.Properties
	if GlobalDefinition.AppID == "" {
		GlobalDefinition.AppID = common.DefaultApp
	}

	return err
}
