package geoip

import (
	"net"
	"regexp"

	"github.com/oschwald/maxminddb-golang"
)

var (
	symbolRegex    = regexp.MustCompile("[^a-zA-Z0-9 ]")
	maxmindCity, _ = maxminddb.Open("./City.mmdb")
	maxmindASN, _  = maxminddb.Open("./ASN.mmdb")
)

func Parse(ip string) GeoIpJson {
	result := GeoIpJson{
		Ip:          ip,
		CountryName: "Unknown",
		CountryCode: "XX",
		Region:      "Unknown",
		Org:         "LalatinaHub",
	}

	if maxmindCity == nil || maxmindASN == nil {
		panic("[-] Maxmind reader failed!")
	}

	var (
		geoipCity = GeoIpCity{}
		geoipASN  = GeoIpASN{}
	)
	maxmindCity.Lookup(net.ParseIP(ip), &geoipCity)
	maxmindASN.Lookup(net.ParseIP(ip), &geoipASN)

	if geoipCity.Country.ISOCode != "" {
		for _, country := range CountryList {
			if geoipCity.Country.ISOCode == country.Code {
				result.CountryName = country.Name
				result.CountryCode = country.Code
				result.Region = country.Region
				result.Org = symbolRegex.ReplaceAllString(geoipASN.ASNName, "")
				break
			}
		}
	}

	return result
}
