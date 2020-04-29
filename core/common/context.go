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

	Returned bool
}

func (c *Context) Return() {
	if c.Returned == true {
		return
	}

	_ = c.Client.WriteMsg(c.Response)
	_ = c.Client.Close()
	c.Returned = true
}

func (c *Context) Abort() {
	_ = c.Client.Close()
	c.Returned = true
}
