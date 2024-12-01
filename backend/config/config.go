package config

import (
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/fileutil"
	"os"
)

var (
	_globalCfg = &GlobalConfig{}
)

type GlobalConfig struct {
	DownloadPath string `json:"download_path" yaml:"download_path"`
	Port         int    `json:"port" yaml:"port"`
}

func init() {

	os.Chdir("D:\\mytest\\mywork\\xhs_downloader\\backend\\cmd")

	if !utils.HasFile("chore/config.yaml") {
		panic("config file not exist")
	}
	err := fileutil.ReadYamlFile("chore/config.yaml", _globalCfg)
	if err != nil {
		panic(err)
	}
}

func GetConfig() *GlobalConfig {
	return _globalCfg
}

func GetDownloadPath() string {
	return "chore/download"
}
