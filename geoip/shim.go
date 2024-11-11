package geoip

import (
	"net"
	"sync"
)

var (
	geo      *Geoip
	syncOnce sync.Once
)

func GetCountryByByteSlice(ip []byte) string {
	syncOnce.Do(func() {
		geo = New()
	})
	if geo == nil {
		return "??"
	}

	country, ok := geo.GetCountryByAddr(net.IP(ip))
	if !ok {
		return "??"
	} else {
		return country
	}
}
