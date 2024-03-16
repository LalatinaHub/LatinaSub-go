package geoip

import (
	"regexp"
)

var (
	symbolRegex = regexp.MustCompile("[^a-zA-Z0-9 ]")
)

func Parse(myIp MyIp) GeoIpJson {
	result := GeoIpJson{
		Ip:          "0.0.0.0",
		CountryName: "Unknown",
		CountryCode: "XX",
		Region:      "Unknown",
		Org:         "LalatinaHub",
	}

	if myIp.Query != "" {
		result.Ip = myIp.Query
	} else {
		result.Ip = myIp.Ip
	}

	for _, country := range CountryList {
		switch country.Code {
		case myIp.Country, myIp.CC:
			result.CountryName = country.Name
			result.CountryCode = country.Code
			result.Region = country.Region
			result.Org = symbolRegex.ReplaceAllString(myIp.Org, "")
			return result
		}
	}

	return result
}
