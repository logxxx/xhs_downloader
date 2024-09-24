package config

import (
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/fileutil"
)

var (
	_globalCfg = &GlobalConfig{}
)

type GlobalConfig struct {
	DownloadPath string `json:"download_path" yaml:"download_path"`
	Port         int    `json:"port" yaml:"port"`
}

func init() {
	if !utils.HasFile("chore/config.yaml") {
		panic("config file not exist")
	}
	err := fileutil.ReadYamlFile("chore/config.yaml", _globalCfg)
	if err != nil {
		panic(err)
	}
	if _globalCfg.Port == 0 {
		_globalCfg.Port = 8080
	}
}

func GetConfig() *GlobalConfig {
	return _globalCfg
}
