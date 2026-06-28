package geoip

import (
	"net"
	"testing"
)

func TestGeoip(t *testing.T) {
	t.Parallel()

	geoip := New()
	if geoip == nil {
		t.Fatal("expected non-nil geoip")
	}

	t.Run("IPv4 Country Mapping Tests", func(t *testing.T) {
		t.Parallel()
		for _, test := range []struct {
			addr, cc string
			ok       bool
		}{
			{"129.97.208.23", "CA", true},
			{"127.0.0.1", "", false},
			{"255.255.255.255", "", false},
			{"0.0.0.0", "", false},
			{"223.252.127.255", "JP", true},
			{"223.252.127.255", "JP", true},
		} {
			country, ok := geoip.GetCountryByAddr(net.ParseIP(test.addr))
			if country != test.cc {
				t.Errorf("addr %s: expected country %q, got %q", test.addr, test.cc, country)
			}
			if ok != test.ok {
				t.Errorf("addr %s: expected ok=%v, got ok=%v", test.addr, test.ok, ok)
			}
		}
	})

	t.Run("IPv6 Country Mapping Tests", func(t *testing.T) {
		t.Parallel()
		for _, test := range []struct {
			addr, cc string
			ok       bool
		}{
			{"2620:101:f000:0:250:56ff:fe80:168e", "CA", true},
			{"fd00:0:0:0:0:0:0:1", "", false},
			{"0:0:0:0:0:0:0:0", "", false},
			{"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", "", false},
			{"2a07:2e47:ffff:ffff:ffff:ffff:ffff:ffff", "FR", true},
			{"2a07:2e40::", "FR", true},
		} {
			country, ok := geoip.GetCountryByAddr(net.ParseIP(test.addr))
			if country != test.cc {
				t.Errorf("addr %s: expected country %q, got %q", test.addr, test.cc, country)
			}
			if ok != test.ok {
				t.Errorf("addr %s: expected ok=%v, got ok=%v", test.addr, test.ok, ok)
			}
		}
	})
}
