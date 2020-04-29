package plugin

import (
	"bufio"
	"github.com/miekg/dns"
	"github.com/pcmid/mdns/core/common"
	"github.com/pcmid/mdns/plugin/lib/domain"
	"github.com/pcmid/mdns/plugin/lib/ipset"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

func init() {
	Register(&Ipset{})
}

type Ipset struct {
	Domains map[string]*domain.Tree
	Set     map[string]*ipset.IPSet
}

func (i *Ipset) Name() string {
	return "ipset"
}

func (i *Ipset) Init(config map[string]interface{}) error {

	i.Domains = make(map[string]*domain.Tree)
	i.Set = make(map[string]*ipset.IPSet)

	sets := config["sets"].(map[string]interface{})

	for name := range sets {

		domains, err := domain.TreeFromFile(sets[name].(map[string]interface{})["domain_file"].(string))

		i.Domains[name] = domains

		set, err := ipset.New(name, "hash:net", &ipset.Params{
			Timeout: 0,
		})
		if set == nil {
			log.Error(err)
			continue
		}
		err = set.Create()
		if err != nil {
			log.Error(err)
		}

		IpFile := sets[name].(map[string]interface{})["ip_file"].(string)

		for _, ip := range parseIPList(IpFile) {
			if err = set.Add(ip, 0); err != nil {
				log.Error(err)
			}

		}

		i.Set[name] = set
	}

	return nil
}

func (i *Ipset) Where() uint8 {
	return OUT
}

func (i *Ipset) HandleDns(ctx *common.Context) {

	if ctx.Response != nil && len(ctx.Response.Answer) <= 0 {
		return
	}
	log.Debug(dns.Field(ctx.Response.Answer[0], 1))

	for setName := range i.Domains {
		if i.Domains[setName].Has(domain.Domain(ctx.Response.Question[0].Name)) {
			for _, ans := range ctx.Response.Answer {
				if ans.Header().Rrtype != dns.TypeA {
					continue
				}

				err := i.Set[setName].Add(dns.Field(ans, 1), func(a, b int) int {
					if a > b {
						return a
					}
					return b
				}(int(ans.Header().Ttl), 3600))
				log.Debugf("ipset add %s to %s", ctx.Response.Question[0].Name, setName)
				if err != nil {
					log.Error(err)
				}
			}
			break
		}
	}
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
