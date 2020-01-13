package config

import (
	"encoding/json"
	"github.com/pcmid/mdns/core/common"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	Addr     string `json:"Addr"`
	Upstream *common.DNSUpstream

	PluginConfDir string
	Plugins       []string
}

func NewConfig(configFile string) *Config {
	f, err := os.Open(configFile)
	if err != nil {
		log.Fatal("Open config file failed: ", err)
		os.Exit(1)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal("Read config file failed: ", err)
		os.Exit(1)
	}

	j := new(Config)
	err = json.Unmarshal(b, j)
	if err != nil {
		log.Fatal("Json syntax error: ", err)
		os.Exit(1)
	}

	if !strings.HasSuffix(j.PluginConfDir,"/") {
		j.PluginConfDir += "/"
	}

	return j
}
