package geoip

import (
	"net"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGeoip(t *testing.T) {
	Convey("Geoip", t, func() {
		geoip := New()
		So(geoip, ShouldNotBeNil)

		Convey("IPv4 Country Mapping Tests", func() {
			for _, test := range []struct {
				addr, cc string
				ok       bool
			}{
				{
					"129.97.208.23", //uwaterloo
					"CA",
					true,
				},
				{
					"127.0.0.1",
					"",
					false,
				},
				{
					"255.255.255.255",
					"",
					false,
				},
				{
					"0.0.0.0",
					"",
					false,
				},
				{
					"223.252.127.255", //test high end of range
					"JP",
					true,
				},
				{
					"223.252.127.255", //test low end of range
					"JP",
					true,
				},
			} {
				country, ok := geoip.GetCountryByAddr(net.ParseIP(test.addr))
				So(country, ShouldEqual, test.cc)
				So(ok, ShouldResemble, test.ok)
			}
		})

		Convey("IPv6 Country Mapping Tests", func() {
			for _, test := range []struct {
				addr, cc string
				ok       bool
			}{
				{
					"2620:101:f000:0:250:56ff:fe80:168e", //uwaterloo
					"CA",
					true,
				},
				{
					"fd00:0:0:0:0:0:0:1",
					"",
					false,
				},
				{
					"0:0:0:0:0:0:0:0",
					"",
					false,
				},
				{
					"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
					"",
					false,
				},
				{
					"2a07:2e47:ffff:ffff:ffff:ffff:ffff:ffff", //test high end of range
					"FR",
					true,
				},
				{
					"2a07:2e40::", //test low end of range
					"FR",
					true,
				},
			} {
				country, ok := geoip.GetCountryByAddr(net.ParseIP(test.addr))
				So(country, ShouldEqual, test.cc)
				So(ok, ShouldResemble, test.ok)
			}
		})
	})
}
