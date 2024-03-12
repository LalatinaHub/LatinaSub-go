package geoip

import (
	"regexp"
)

var (
	symbolRegex = regexp.MustCompile("[^a-zA-Z0-9 ]")
)

func Parse(myIp MyIp) GeoIpJson {
	result := GeoIpJson{
		Ip:          myIp.Ip,
		CountryName: "Unknown",
		CountryCode: "XX",
		Region:      "Unknown",
		Org:         "LalatinaHub",
	}

	for _, country := range CountryList {
		if myIp.CC == country.Code {
			result.CountryName = country.Name
			result.CountryCode = country.Code
			result.Region = country.Region
			result.Org = symbolRegex.ReplaceAllString(myIp.Org, "")
			break
		}
	}

	return result
}
