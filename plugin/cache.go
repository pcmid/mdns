// source: https://github.com/shawn1m/overture/blob/master/core/cache/cache.go
// Copyright (c) 2014 The SkyDNS Authors. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be
// found in the LICENSE file.
package plugin

import (
	"encoding/json"
	"github.com/miekg/dns"
	"github.com/pcmid/mdns/core/common"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"sync"
	"time"
)

type elem struct {
	expiration time.Time
	msg        *dns.Msg
}

type Cache struct {
	sync.RWMutex

	capacity int
	table    map[dns.Question]*elem
	mTTL     uint32
}

type CacheConfig struct {
	Capacity int
	MTTL     uint32
}

func (c *Cache) Name() string {
	return "cache"
}

func (c *Cache) Init(configDir string) error {
	configData, _ := ioutil.ReadFile(configDir + "cache.json")
	conf := CacheConfig{}
	err := json.Unmarshal(configData, &conf)

	if err != nil {
		return err
	}

	c.table = make(map[dns.Question]*elem)
	c.capacity = conf.Capacity
	c.mTTL = conf.MTTL

	return nil
}

func (c *Cache) Where() uint8 {
	return BOTH
}

func (c *Cache) HandleDns(ctx *common.Context) {
	if ctx.IsQuery {
		if msg := c.Hit(ctx.Query.Question[0], ctx.Query.Id); msg != nil {
			ctx.Response = msg
			ctx.Return()
		}
	} else {
		if ctx.Response != nil {
			c.InsertMessage(ctx.Response)
		}
	}
}

func init() {
	Register(&Cache{})
}

func (c *Cache) Capacity() int { return c.capacity }

func (c *Cache) Remove(q dns.Question) {
	c.Lock()
	delete(c.table, q)
	c.Unlock()
}

func (c *Cache) EvictRandom() {
	cacheLength := len(c.table)
	if cacheLength <= c.capacity {
		return
	}
	i := c.capacity - cacheLength
	for k := range c.table {
		delete(c.table, k)
		i--
		if i == 0 {
			break
		}
	}
}

func (c *Cache) InsertMessage(m *dns.Msg) {
	if c.capacity <= 0 || m == nil {
		return
	}

	c.Lock()
	defer c.Unlock()

	ttl := c.mTTL
	if len(m.Answer) > 0 && ttl < m.Answer[0].Header().Ttl {
		ttl = m.Answer[0].Header().Ttl
	}

	ttlDuration := time.Duration(ttl) * time.Second
	if _, ok := c.table[m.Question[0]]; !ok {
		c.table[m.Question[0]] = &elem{time.Now().UTC().Add(ttlDuration), m.Copy()}
	}
	log.Debugf("Cached: %s, TTL: %d", &m.Question[0], ttl)
	c.EvictRandom()
}

func (c *Cache) Search(q dns.Question) (*dns.Msg, time.Time, bool) {
	if c.capacity <= 0 {
		return nil, time.Time{}, false
	}
	c.RLock()
	defer c.RUnlock()

	if e, ok := c.table[q]; ok {
		e1 := e.msg.Copy()
		return e1, e.expiration, true
	}
	return nil, time.Time{}, false
}
func (c *Cache) Hit(q dns.Question, msgid uint16) *dns.Msg {
	m, exp, hit := c.Search(q)
	if hit {
		// Cache hit! \o/
		if time.Since(exp) < 0 {
			m.Id = msgid
			m.Compress = true
			// Even if something ended up with the TC bit *in* the cache, set it to off
			m.Truncated = false
			for _, a := range m.Answer {
				a.Header().Ttl = uint32(time.Since(exp).Seconds() * -1)
			}
			return m
		}
		// Expired! /o\
		c.Remove(q)
	}
	return nil
}
