package configcenter

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/ServiceComb/go-chassis/bootstrap"
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/endpoint-discovery"
	"github.com/ServiceComb/go-chassis/core/lager"
	chassisTLS "github.com/ServiceComb/go-chassis/core/tls"

	// go-archaius package is imported for to initialize the config-center configurations
	"github.com/ServiceComb/go-archaius"
	"github.com/ServiceComb/go-archaius/core"
	"github.com/ServiceComb/go-archaius/sources/configcenter-source"
	"github.com/ServiceComb/go-cc-client"
	"github.com/ServiceComb/go-cc-client/member-discovery"
)

const (
	//Name is a variable of type string
	Name          = "configcenter"
	maxValue      = 256
	emptyDimeInfo = "Issue with regular expression or exceeded the max length"
)

// InitConfigCenter initialize config center
func InitConfigCenter() error {
	configCenterURL, err := isConfigCenter()
	if err != nil {
		return nil
	}

	var enableSSL bool
	tlsConfig, tlsError := getTLSForClient(configCenterURL)
	if tlsError != nil {
		lager.Logger.Errorf(tlsError, "Get %s.%s TLS config failed, err:", Name, common.Consumer)
		return tlsError
	}

	/*This condition added because member discovery can have multiple ip's with IsHTTPS
	having both true and false value.*/
	if tlsConfig != nil {
		enableSSL = true
	}

	dimensionInfo := getUniqueIDForDimInfo()

	if dimensionInfo == "" {
		err := errors.New("empty dimension info: " + emptyDimeInfo)
		lager.Logger.Error("empty dimension info", err)
		return err
	}

	if config.GlobalDefinition.Cse.Config.Client.TenantName == "" {
		config.GlobalDefinition.Cse.Config.Client.TenantName = common.DefaultTenant
	}

	if config.GlobalDefinition.Cse.Config.Client.RefreshInterval == 0 {
		config.GlobalDefinition.Cse.Config.Client.RefreshInterval = 30
	}

	err = initConfigCenter(configCenterURL,
		dimensionInfo, config.GlobalDefinition.Cse.Config.Client.TenantName,
		enableSSL, tlsConfig)
	if err != nil {
		lager.Logger.Error("failed to init config center", err)
		return err
	}

	lager.Logger.Warnf("config center init success")
	return nil
}

func isConfigCenter() (string, error) {
	configCenterURL := config.GlobalDefinition.Cse.Config.Client.ServerURI
	if configCenterURL == "" {
		ccURL, err := endpoint.GetEndpointFromServiceCenter("default", "CseConfigCenter", "latest")
		if err != nil {
			lager.Logger.Errorf(err, "empty config center endpoint, please provide the config center endpoint")
			return "", err
		}

		configCenterURL = ccURL
	}

	return configCenterURL, nil
}

func getTLSForClient(configCenterURL string) (*tls.Config, error) {
	if !strings.Contains(configCenterURL, "://") {
		return nil, nil
	}
	ccURL, err := url.Parse(configCenterURL)
	if err != nil {
		lager.Logger.Error("Error occurred while parsing config center Server Uri", err)
		return nil, err
	}
	if ccURL.Scheme == common.HTTP {
		return nil, nil
	}

	sslTag := Name + "." + common.Consumer
	tlsConfig, sslConfig, err := chassisTLS.GetTLSConfigByService(Name, "", common.Consumer)
	if err != nil {
		if chassisTLS.IsSSLConfigNotExist(err) {
			return nil, fmt.Errorf("%s TLS mode, but no ssl config", sslTag)
		}
		return nil, err
	}
	lager.Logger.Warnf("%s TLS mode, verify peer: %t, cipher plugin: %s.",
		sslTag, sslConfig.VerifyPeer, sslConfig.CipherPlugin)

	return tlsConfig, nil
}

func getUniqueIDForDimInfo() string {
	serviceName := config.MicroserviceDefinition.ServiceDescription.Name
	version := config.MicroserviceDefinition.ServiceDescription.Version
	appName := config.MicroserviceDefinition.AppID
	if appName == "" {
		appName = config.GlobalDefinition.AppID
	}

	if appName != "" {
		serviceName = serviceName + "@" + appName
	}

	if version != "" {
		serviceName = serviceName + "#" + version
	}

	if len(serviceName) > maxValue {
		lager.Logger.Errorf(nil, "exceeded max value %d for dimensionInfo %s with length %d", maxValue, serviceName,
			len(serviceName))
		return ""
	}

	dimeExp := `\A([^\$\%\&\+\(/)\[\]\" "\"])*\z`
	dimRegexVar, err := regexp.Compile(dimeExp)
	if err != nil {
		lager.Logger.Error("not a valid regular expression", err)
		return ""
	}

	if !dimRegexVar.Match([]byte(serviceName)) {
		lager.Logger.Errorf(nil, "invalid value for dimension info, doesnot setisfy the regular expression for dimInfo:%s",
			serviceName)
		return ""
	}

	return serviceName
}

func initConfigCenter(ccEndpoint, dimensionInfo, tenantName string, enableSSL bool, tlsConfig *tls.Config) error {

	refreshMode := archaius.GetInt("cse.config.client.refreshMode", common.DefaultRefreshMode)
	if refreshMode != 0 && refreshMode != 1 {
		err := errors.New(client.RefreshModeError)
		lager.Logger.Error(client.RefreshModeError, err)
		return err
	}
	memDiscovery := memberdiscovery.NewConfiCenterInit(tlsConfig, tenantName, enableSSL)

	configCenters := strings.Split(ccEndpoint, ",")
	cCenters := make([]string, 0)
	for _, value := range configCenters {
		value = strings.Replace(value, " ", "", -1)
		cCenters = append(cCenters, value)
	}

	memDiscovery.ConfigurationInit(cCenters)

	if enbledAutoDiscovery() {
		refreshError := memDiscovery.RefreshMembers()
		if refreshError != nil {
			lager.Logger.Error(client.ConfigServerMemRefreshError, refreshError)
			return errors.New(client.ConfigServerMemRefreshError)
		}
	}

	configCenterSource := configcentersource.NewConfigCenterSource(
		memDiscovery, dimensionInfo, tlsConfig, tenantName, refreshMode,
		config.GlobalDefinition.Cse.Config.Client.RefreshInterval, enableSSL)

	err := archaius.DefaultConf.ConfigFactory.AddSource(configCenterSource)
	if err != nil {
		lager.Logger.Error("failed to do add source operation!!", err)
		return err
	}

	// Get the whole configuration
	//config := factory.GetConfigurations()
	//logger.Info("init config center %+v", config)

	eventHandler := EventListener{
		Name:    "EventHandler",
		Factory: archaius.DefaultConf.ConfigFactory,
	}

	memberdiscovery.MemberDiscoveryService = memDiscovery
	archaius.DefaultConf.ConfigFactory.RegisterListener(eventHandler, "a*")

	if err := refreshGlobalConfig(); err != nil {
		lager.Logger.Error("failed to refresh global config for lb and cb", err)
		return err
	}

	return nil
}

func refreshGlobalConfig() error {
	err := config.ReadHystrixFromArchaius()
	if err != nil {
		return err
	}
	return config.ReadLBFromArchaius()
}

//EventListener is a struct
type EventListener struct {
	Name    string
	Factory goarchaius.ConfigurationFactory
}

//Event is a method
func (e EventListener) Event(event *core.Event) {
	value := e.Factory.GetConfigurationByKey(event.Key)
	lager.Logger.Infof("config value %s | %s", event.Key, value)
}

func enbledAutoDiscovery() bool {
	if config.GlobalDefinition.Cse.Config.Client.Autodiscovery {
		return true
	}

	return false
}
func init() {
	bootstrap.InstallPlugin("config_center",
		bootstrap.BootstrapFunc(InitConfigCenter))
}
