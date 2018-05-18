package server

import (
	"fmt"
	"log"

	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"
	chassisTLS "github.com/ServiceComb/go-chassis/core/tls"
	"github.com/ServiceComb/go-chassis/util/iputil"
)

//NewFunc returns a ProtocolServer
type NewFunc func(Options) ProtocolServer

var serverPlugins = make(map[string]NewFunc)
var servers = make(map[string]ProtocolServer)

//InstallPlugin For developer
func InstallPlugin(protocol string, newFunc NewFunc) {
	serverPlugins[protocol] = newFunc
	log.Printf("Installed Server Plugin, protocol=%s", protocol)
}

//GetServerFunc returns the server function
func GetServerFunc(protocol string) (NewFunc, error) {
	f, ok := serverPlugins[protocol]
	if !ok {
		return nil, fmt.Errorf("Don't support protocol [%s]", protocol)
	}
	return f, nil
}

//GetServer return the server based on protocol
func GetServer(protocol string) (ProtocolServer, error) {
	s, ok := servers[protocol]
	if !ok {
		return nil, fmt.Errorf("[%s] server isn't running ", protocol)
	}
	return s, nil
}

//GetServers returns the map of servers
func GetServers() map[string]ProtocolServer {
	return servers
}

//ServerErr server error
var ServerErr = make(chan error)

//StartServer starting the server
func StartServer() error {

	for name, server := range servers {
		lager.Logger.Info("starting server " + name + "...")
		err := server.Start()
		if err != nil {
			lager.Logger.Errorf(err, "servers failed to start")
			return fmt.Errorf("Can not start [%s] server,%s", name, err.Error())
		}
		lager.Logger.Info(name + " server start success")
	}
	lager.Logger.Info("All server Start Completed")

	return nil
}

//UnRegistrySelfInstances this function removes the self instance
func UnRegistrySelfInstances() error {
	microserviceIDs := registry.SelfInstancesCache.Items()
	for mid := range microserviceIDs {
		value, ok := registry.SelfInstancesCache.Get(mid)
		if !ok {
			lager.Logger.Warnf("StartServer() get SelfInstancesCache failed, mid: %s", mid)
		}
		instanceIDs, ok := value.([]string)
		if !ok {
			lager.Logger.Warnf("StartServer() type asserts failed, mid: %s", mid)
		}
		for _, iid := range instanceIDs {
			err := registry.DefaultRegistrator.UnRegisterMicroServiceInstance(mid, iid)
			if err != nil {
				lager.Logger.Errorf(err, "StartServer() UnregisterMicroServiceInstance failed, mid/iid: %s/%s", mid, iid)
				return err
			}
		}
	}
	return nil
}

//Init initializes
func Init() error {
	var err error
	err = initialGlobal()
	if err != nil {
		return err
	}

	return nil
}

func initialGlobal() error {
	var err error
	for k, v := range config.GlobalDefinition.Cse.Protocols {

		if err = initialSingle(config.GlobalDefinition.Cse.Handler.Chain.Provider, v, k); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func initialSingle(providerMap map[string]string, p model.Protocol, name string) error {
	lager.Logger.Debugf("Init server of protocol [%s]", name)
	f, err := GetServerFunc(name)
	if err != nil {
		return fmt.Errorf("Do not support [%s] server", name)
	}

	sslTag := name + "." + common.Provider
	tlsConfig, sslConfig, err := chassisTLS.GetTLSConfigByService("", name, common.Provider)
	if err != nil {
		if !chassisTLS.IsSSLConfigNotExist(err) {
			return err
		}
	} else {
		lager.Logger.Warnf("%s TLS mode, verify peer: %t, cipher plugin: %s.",
			sslTag, sslConfig.VerifyPeer, sslConfig.CipherPlugin)
	}

	if p.Listen == "" {
		if p.Advertise != "" {
			p.Listen = p.Advertise
		} else {
			p.Listen = iputil.DefaultEndpoint4Protocol(name)
		}
	}
	chainName := common.DefaultChainName
	for name := range providerMap {
		if name != common.DefaultApp {
			chainName = name
			break
		}
	}

	var s ProtocolServer
	o := Options{
		Address:   p.Listen,
		ChainName: chainName,
		TLSConfig: tlsConfig,
	}
	s = f(o)
	servers[name] = s
	return nil
}
