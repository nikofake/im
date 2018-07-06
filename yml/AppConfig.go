package yml

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/gogap/logrus"
)

var (
	PATH             = "./config/config.yml"
	Config  = AppConfig{}
)

type AppConfig struct {
	Monitor struct {
		Port int32
	}
	Connection struct {
		Port int32
	}

	Zk struct{
		Hosts []string `yaml:",flow"`
		Basepath string
	}
}

func init() {
	if content, err := ioutil.ReadFile(PATH); err != nil {
		panic(err)
	} else {
		if err = yaml.Unmarshal(content, &Config); err != nil {
			panic(err)
		}
	}
	logrus.Info("init config", Config)
}
