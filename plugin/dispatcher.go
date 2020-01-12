package plugin

import (
	"encoding/json"
	"github.com/pcmid/mdns/core/common"
	"github.com/pcmid/mdns/plugin/lib/domain"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

func init() {
	Register(&Dispatcher{})
}

type DispatcherConf struct {
	Areas map[string]struct {
		DomainFile string
		Upstream   *common.DNSUpstream
	}
}

type Dispatcher struct {
	Areas map[string]struct {
		Domains  *domain.Tree
		Upstream *common.DNSUpstream
	}
}

func (d *Dispatcher) Name() string {
	return "dispatcher"
}

func (d *Dispatcher) Init(configDir string) error {

	d.Areas = make(map[string]struct {
		Domains  *domain.Tree
		Upstream *common.DNSUpstream
	})

	configData, err := ioutil.ReadFile(configDir + "dispatcher.json")

	if err != nil {
		return err
	}

	conf := DispatcherConf{}

	err = json.Unmarshal(configData, &conf)

	if err != nil {
		return err
	}

	for i := range conf.Areas {
		area := struct {
			Domains  *domain.Tree
			Upstream *common.DNSUpstream
		}{}
		area.Upstream = conf.Areas[i].Upstream
		area.Domains, err = domain.TreeFromFile(conf.Areas[i].DomainFile)
		if err != nil {
			log.Error(err)
			continue
		}

		d.Areas[i] = area
	}

	return nil
}

func (d *Dispatcher) HandleDns(ctx *common.Context) {
	_domain := ctx.Query.Question[0].Name
	for name, area := range d.Areas {
		if area.Domains.Has(domain.Domain(_domain)) {
			log.Debugf("%s switch to area %s", _domain, name)
			ctx.Upstream = area.Upstream
			return
		}
	}

	log.Debugf("%s not matched", _domain)
}

func (d *Dispatcher) Where() uint8 {
	return IN
}
