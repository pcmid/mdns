package plugin

import (
	"bufio"
	"encoding/json"
	"github.com/miekg/dns"
	"github.com/pcmid/mdns/core/common"
	"github.com/pcmid/mdns/plugin/lib/domain"
	"github.com/pcmid/mdns/plugin/lib/ipset"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type IpSet struct {
	Domains map[string]*domain.Tree
	Set     map[string]*ipset.IPSet
}

type IpSetConfig []struct {
	Name       string `json:"name"`
	DomainFile string `json:"domain_file"`
}

func (i *IpSet) Name() string {
	return "ipset"
}

func (i *IpSet) Init(configDir string) error {
	jsonData, err := ioutil.ReadFile(configDir + "ipset.json")

	if err != nil {
		return err
	}

	var config IpSetConfig

	err = json.Unmarshal(jsonData, &config)

	if err != nil {
		return err
	}

	i.Domains = make(map[string]*domain.Tree)
	i.Set = make(map[string]*ipset.IPSet)

	for _, conf := range config {

		domains, err := domain.TreeFromFile(conf.DomainFile)
		if err != nil {
			continue
		}
		i.Domains[conf.Name] = domains

		set, err := ipset.New(conf.Name, "hash:net", &ipset.Params{})
		if set == nil {
			log.Error(err)
			continue
		}

		err = set.Create()
		if err != nil {
			log.Error(err)
		}

		err = set.Flush()
		if err != nil {
			log.Error(err)
		}
		i.Set[conf.Name] = set
	}

	return nil
}

func (i *IpSet) Where() uint8 {
	return OUT
}

func (i *IpSet) HandleDns(ctx *common.Context) {
	if ctx.Response != nil && len(ctx.Response.Answer) <= 0 {
		return
	}
	log.Debug(dns.Field(ctx.Response.Answer[0], 1))

	for setName := range i.Domains {
		if i.Domains[setName].Has(domain.Domain(ctx.Response.Question[0].Name)) {
			for _, ans := range ctx.Response.Answer {
				err := i.Set[setName].Add(dns.Field(ans, 1), 0)
				log.Debugf("ipset add %s to %s", ctx.Response.Question[0].Name, setName)
				if err != nil {
					log.Error(err)
				}
			}
			break
		}
	}
}

func init() {
	Register(&IpSet{})
}

func parseIPList(file string) []string {
	ipf, err := os.Open(file)
	if err != nil {
		log.Error(err)
		return nil
	}
	defer func() {
		_ = ipf.Close()
	}()

	buf := bufio.NewReader(ipf)

	ipList := make([]string, 0)

	for {

		line, errRead := buf.ReadString('\n')

		if errRead != nil && errRead != io.EOF {
			continue
		}

		line = strings.TrimSpace(line)

		ipList = append(ipList, line)

		if errRead == io.EOF {
			break
		}
	}

	return ipList
}

func parseDomainList(file string) *domain.Tree {
	dt, err := domain.TreeFromFile(file)

	if err != nil {
		log.Error(err)
		return nil
	}

	return dt
}
