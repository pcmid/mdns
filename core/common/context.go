package common

import (
	"github.com/miekg/dns"
)

type Context struct {
	Client   dns.ResponseWriter
	Err      error
	Upstream *DNSUpstream
	Query    *dns.Msg
	Response *dns.Msg

	IsQuery bool

	returned bool
}

func (c *Context) Return() {
	if c.returned == true {
		return
	}

	_ = c.Client.WriteMsg(c.Response)
	_ = c.Client.Close()
	c.returned = true
}

func (c *Context) Abort() {
	_ = c.Client.Close()
	c.returned = true
}
