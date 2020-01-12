package plugin

import (
	"github.com/pcmid/mdns/core/common"
)

const (
	IN   = 0x01
	OUT  = 0x02
	BOTH = 0x03
)

var Plugins map[string]Plugin

func Register(p Plugin) {
	if Plugins == nil {
		Plugins = make(map[string]Plugin)
	}
	Plugins[p.Name()] = p
}

func Get(name string) Plugin {
	if Plugins == nil {
		return nil
	}
	return Plugins[name]
}

func init() {

}

type Plugin interface {
	Name() string
	Init(configDir string) error
	Where() uint8
	HandleDns(*common.Context)
}
