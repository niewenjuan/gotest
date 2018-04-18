package bootstrap

import (
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/metrics"
)

var bootstrapPlugins map[string]BootstrapPlugin

//BootstrapPlugin is a interface which declares Init method
type BootstrapPlugin interface {
	Init() error
}

// The BootstrapFunc type is an adapter to allow the use of ordinary functions as bootstrapPlugin.
type BootstrapFunc func() error

//Init is a method
func (b BootstrapFunc) Init() error {
	return b()
}

//InstallPlugin is a function which installs plugin
func InstallPlugin(name string, plugin BootstrapPlugin) {
	bootstrapPlugins[name] = plugin
}

//Bootstrap is a function which logs message
func Bootstrap() {
	if _, ok := bootstrapPlugins["EE"]; ok {

	} /*else if _, ok := bootstrapPlugins["CE"]; ok {
		lager.Logger.Info("Bootstrap Huawei Community Edition.")
	}*/

	for name, p := range bootstrapPlugins {
		lager.Logger.Info("Bootstrap " + name)
		if err := p.Init(); err != nil {
			lager.Logger.Errorf(err, "Failed to init %s.", name)
		}

	}
}

func init() {
	bootstrapPlugins = make(map[string]BootstrapPlugin)

	InstallPlugin("metric", BootstrapFunc(metrics.Init))
}
