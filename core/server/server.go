package server

import (
	"crypto/tls"
	"github.com/pcmid/mdns/core/common"
	"github.com/pcmid/mdns/plugin"
	"golang.org/x/net/proxy"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	addr string

	upstream *common.DNSUpstream

	inPlugin  []plugin.Plugin
	outPlugin []plugin.Plugin
}

func NewServer(addr string, upstream *common.DNSUpstream, plugins []plugin.Plugin) *Server {
	s := &Server{
		addr:     addr,
		upstream: upstream,
	}

	for _, p := range plugins {
		if p.Where()&plugin.IN != 0 {
			s.inPlugin = append(s.inPlugin, p)
		}

		if p.(plugin.Plugin).Where()&plugin.OUT != 0 {
			s.outPlugin = append(s.outPlugin, p)
		}
	}

	return s
}

func (s *Server) Run() {
	mux := dns.NewServeMux()
	mux.Handle(".", s)

	wg := new(sync.WaitGroup)
	wg.Add(1)

	log.Info("Start mdns on " + s.addr)

	for _, p := range [2]string{"tcp", "udp"} {
		go func(p string) {
			err := dns.ListenAndServe(s.addr, p, mux)
			if err != nil {
				log.Fatal("Listen "+p+" failed: ", err)
				os.Exit(1)
			}
		}(p)
	}

	wg.Wait()
}

func (s *Server) ServeDNS(w dns.ResponseWriter, q *dns.Msg) {
	ctx := &common.Context{
		Client:   w,
		Upstream: s.upstream,
		Query:    q,
		IsQuery:  true,
	}

	if len(ctx.Query.Question) <= 0 {
		log.Warnf("nil request from %s", ctx.Client.RemoteAddr())
		ctx.Abort()
		return
	}

	for _, p := range s.inPlugin {
		p.HandleDns(ctx)
	}

	s.Exchange(ctx)
	ctx.IsQuery = false

	if ctx.Response == nil {
		log.Errorf("nil response for %s", ctx.Query.Question[0].Name)
		dns.HandleFailed(w, q)
		return
	}

	if ctx.Err != nil {
		log.Error(ctx.Err)
	}

	for _, p := range s.outPlugin {
		p.HandleDns(ctx)
	}

	ctx.Return()
}

func (s *Server) Exchange(ctx *common.Context) {
	if ctx.Upstream == nil {
		ctx.Upstream = s.upstream
	}

	upstream := ctx.Upstream

	var conn net.Conn
	if upstream.SOCKS5Address != "" {
		s, e := proxy.SOCKS5(upstream.Protocol, upstream.SOCKS5Address, nil, proxy.Direct)
		if e != nil {
			log.Errorf("get socks5 proxy dialer failed: %v", e)
			return
		}
		conn, e = s.Dial(upstream.Protocol, upstream.Address)
		if e != nil {
			log.Errorf("dial DNS upstream with SOCKS5 proxy failed: %v", e)
			return
		}
	} else if upstream.Protocol == "tcp-tls" {
		var err error
		conf := &tls.Config{
			InsecureSkipVerify: false,
		}
		s := strings.Split(upstream.Address, "@")
		if len(s) == 2 {
			var servername, port string
			if servername, port, err = net.SplitHostPort(s[0]); err != nil {
				log.Errorf("DNS-over-TLS servername:port@serverAddress config failed: %v", err)
				return
			}
			conf.ServerName = servername
			upstream.Address = s[1] + ":" + port
		}
		if conn, err = tls.Dial("tcp", upstream.Address, conf); err != nil {
			log.Errorf("dial DNS-over-TLS upstream failed: %v", err)
			return
		}
	} else {
		var err error
		if conn, err = net.Dial(upstream.Protocol, upstream.Address); err != nil {
			log.Errorf("dial DNS upstream failed: %v", err)
			return
		}
	}

	dnsTimeout := time.Duration(upstream.Timeout) * time.Second / 3

	_ = conn.SetDeadline(time.Now().Add(dnsTimeout))
	_ = conn.SetReadDeadline(time.Now().Add(dnsTimeout))
	_ = conn.SetWriteDeadline(time.Now().Add(dnsTimeout))

	dc := &dns.Conn{Conn: conn}
	defer func() {
		_ = dc.Close()
	}()
	err := dc.WriteMsg(ctx.Query)
	if err != nil {
		log.Errorf("%s: send question message failed: %v", upstream.Name, err)
		return
	}
	temp, err := dc.ReadMsg()

	if err != nil {
		log.Errorf("%s: read record message failed: %v", upstream.Name, err)
		return
	}
	if temp == nil {
		log.Errorf("%s: Response message is nil, maybe timeout, please check your query or dns configuration", upstream.Name)
		return
	}
	ctx.Response = temp
	return
}
