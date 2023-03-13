package ipapi

import "encoding/json"

func Parse(str string) Ipapi {
	var ipapi Ipapi
	json.Unmarshal([]byte(str), &ipapi)

	if ipapi.CountryCode != "" {
		for _, country := range CountryList {
			if ipapi.CountryCode == country.Code {
				ipapi.Region = country.Region
				break
			}
		}
	} else {
		ipapi.CountryCode = "XX"
		ipapi.Region = "Unknown"
		ipapi.Org = "Lalatina"
	}

	return ipapi
}
