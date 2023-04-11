package ipapi

import (
	"encoding/json"
	"strings"
)

func Parse(str string) Ipapi {
	var ipapi Ipapi
	json.Unmarshal([]byte(str), &ipapi)

	// Detect is IPv6
	if strings.Contains(ipapi.Ip, ":") {
		ipapi.Ip = ""
	}

	if ipapi.CountryCode != "" {
		for _, country := range CountryList {
			if ipapi.CountryCode == country.Code {
				ipapi.Region = country.Region
				ipapi.CountryName = country.Name
				break
			}
		}
	} else {
		ipapi.CountryName = "Unknown"
		ipapi.CountryCode = "XX"
		ipapi.Region = "Unknown"
		ipapi.Org = "Lalatina"
	}

	return ipapi
}
