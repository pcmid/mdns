package core

import (
	"github.com/pcmid/mdns/core/config"
	"github.com/pcmid/mdns/core/server"
	"github.com/pcmid/mdns/plugin"
	log "github.com/sirupsen/logrus"
)

func InitServer(configFilePath string) {

	conf := config.NewConfig(configFilePath)

	var plugins []plugin.Plugin

	for _, name := range conf.Plugins {
		p := plugin.Get(name)

		if p == nil {
			log.Errorf("unknown plugin: %s", name)
			continue
		}

		p.Init(conf.PluginConfDir)
		plugins = append(plugins, p)
	}

	s := server.NewServer(conf.Addr, conf.Upstream, plugins)

	s.Run()
}
