package plugin

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pcmid/mdns/core/common"
	"github.com/pcmid/mdns/plugin/lib/domain"
	log "github.com/sirupsen/logrus"
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

type Area struct {
	Domains  *domain.Tree
	Upstream *common.DNSUpstream
}

type Dispatcher struct {
	Areas map[string]Area
}

func (d *Dispatcher) Name() string {
	return "dispatcher"
}

func (d *Dispatcher) Init(config map[string]interface{}) error {

	d.Areas = make(map[string]Area)

	areas := config["areas"].(map[string]interface{})

	for name, conf := range areas {
		area := Area{}

		area.Upstream = new(common.DNSUpstream)
		_ = mapstructure.Decode(conf.(map[string]interface{})["upstream"], area.Upstream)

		var err error
		area.Domains, err = domain.TreeFromFile(conf.(map[string]interface{})["domain_file"].(string))
		if err != nil {
			log.Error(err)
			continue
		}

		d.Areas[name] = area
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
