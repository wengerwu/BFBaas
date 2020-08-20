package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/paybf/baasmanager/baas-core/common/log"
	"github.com/spf13/viper"
	"os"
)

var Config *viper.Viper

var logger = log.GetLogger("gateway.config", log.INFO)

func init() {
	go watchConfig()
	loadConfig()
}

func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logger.Info("Config file changed:", e.Name)
		loadConfig()
	})
}

func loadConfig() {
	viper.SetConfigName("gwconfig")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/baas")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Errorf("Fatal error config file: %s \n", err)
		os.Exit(-1)
	}

	Config = viper.GetViper()
	logger.Infof("%v", Config.AllSettings())
}
