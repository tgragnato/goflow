package sampler

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type cache struct {
	items map[string]string
	sync.RWMutex
}

func newcache() *cache {
	cache := &cache{
		map[string]string{},
		sync.RWMutex{},
	}
	go func() {
		for range time.Tick(time.Hour) {
			go cache.update()
		}
	}()
	return cache
}

func (c *cache) get(ip net.IP) (string, bool) {
	c.RLock()
	defer c.RUnlock()
	v, ok := c.items[ip.String()]
	return v, ok
}

func (c *cache) set(ip net.IP) {
	c.Lock()
	defer c.Unlock()
	c.items[ip.String()] = getHostnameFromIp(ip)
	fmt.Printf("New sampler discovered: (%s)=> %s\n", ip.String(), c.items[ip.String()])
}

func getHostnameFromIp(ip net.IP) string {
	hostname, err := net.LookupAddr(ip.String())
	if err != nil || len(hostname) <= 0 {
		return ""
	}
	if hostname[0][len(hostname[0])-1] == '.' {
		return hostname[0][:len(hostname[0])-1]
	}
	return hostname[0]
}

func (c *cache) update() {
	c.Lock()
	defer c.Unlock()
	for k := range c.items {
		newHostname := getHostnameFromIp(net.ParseIP(k))
		if newHostname != "" && newHostname != c.items[k] {
			c.items[k] = newHostname
		}
	}
}
