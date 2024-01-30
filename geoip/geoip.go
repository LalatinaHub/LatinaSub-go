package geoip

import (
	"io"
	"net"
	"net/http"
	"regexp"

	"github.com/oschwald/maxminddb-golang"
)

const (
	CityMmdbUrl = "https://git.io/GeoLite2-City.mmdb"
	ASNMmdbUrl  = "https://git.io/GeoLite2-ASN.mmdb"
)

var (
	symbolRegex = regexp.MustCompile("[^a-zA-Z0-9 ]")

	CityMmdbResp, _ = http.Get(CityMmdbUrl)
	ASNMmdbResp, _  = http.Get(ASNMmdbUrl)

	CityMmdbBytes, _ = io.ReadAll(CityMmdbResp.Body)
	ASNMmdbBytes, _  = io.ReadAll(ASNMmdbResp.Body)

	maxmindCity, _ = maxminddb.FromBytes(CityMmdbBytes)
	maxmindASN, _  = maxminddb.FromBytes(ASNMmdbBytes)
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
