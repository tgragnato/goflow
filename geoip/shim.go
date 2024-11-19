package geoip

import (
	"fmt"
	"math"
	"net"

	"github.com/oschwald/geoip2-golang"
)

var (
	geo           *Geoip
	asnReader     *geoip2.Reader
	countryReader *geoip2.Reader
)

func GetASNByByteSlice(ip []byte) (number uint32, organization string) {
	if asnReader == nil {
		return
	}

	asn, err := asnReader.ASN(net.IP(ip))
	if err != nil {
		return
	}

	if asn.AutonomousSystemNumber <= math.MaxUint32 {
		number = uint32(asn.AutonomousSystemNumber)
	}
	organization = asn.AutonomousSystemOrganization
	return
}

func GetCountryByByteSlice(ip []byte) string {
	if country, ok := geo.GetCountryByAddr(net.IP(ip)); ok {
		return country
	}
	if countryReader == nil {
		return "??"
	}

	country, err := countryReader.Country(net.IP(ip))
	if err == nil {
		return country.Country.IsoCode
	}
	return "??"
}

func Init(GeoipASN string, GeoipCountry string) {
	geo = New()
	if db, err := geoip2.Open(GeoipASN); err == nil {
		asnReader = db
		fmt.Printf("ASN database loaded: epoch(%d)\n", asnReader.Metadata().BuildEpoch)
	} else {
		fmt.Println("Error opening ASN database: ", err.Error())
	}
	if db, err := geoip2.Open(GeoipCountry); err == nil {
		countryReader = db
		fmt.Printf("Country database loaded: epoch(%d)\n", countryReader.Metadata().BuildEpoch)
	} else {
		fmt.Println("Error opening country database: ", err.Error())
	}
}
