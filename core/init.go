package core

import (
	"github.com/pcmid/mdns/core/config"
	"github.com/pcmid/mdns/core/server"
	"github.com/pcmid/mdns/plugin"
	log "github.com/sirupsen/logrus"
)

func InitServer(configFilePath string) *server.Server {

	conf := config.NewConfig(configFilePath)

	var plugins []plugin.Plugin

	for i := range conf.Plugins {

		if conf.Plugins[i].Name == "" {
			continue
		}

		p := plugin.Get(conf.Plugins[i].Name)

		if p == nil {
			log.Errorf("unknown plugin: %s", conf.Plugins[i].Name)
			continue
		}

		if err := p.Init(conf.Plugins[i].Config); err == nil {
			plugins = append(plugins, p)
		} else {
			log.Errorf("Failed to init plugin: %s", conf.Plugins[i].Name)
		}
	}

	s := server.NewServer(conf.Addr, conf.Upstream, plugins)

	return s
}
