package ipapi

import "encoding/json"

func Parse(str string) Ipapi {
	var ipapi Ipapi
	json.Unmarshal([]byte(str), &ipapi)

	for _, country := range CountryList {
		if ipapi.CountryCode == country.Code {
			ipapi.Region = country.Region
			break
		}
	}

	return ipapi
}
