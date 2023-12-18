package helper

import "github.com/LalatinaHub/LatinaSub-go/geoip"

func CCToEmoji(cc string) string {
	for _, country := range geoip.CountryList {
		if cc == country.Code {
			return string(0x1F1E6+rune(country.Code[0])-'A') + string(0x1F1E6+rune(country.Code[1])-'A')
		}
	}

	return cc
}
