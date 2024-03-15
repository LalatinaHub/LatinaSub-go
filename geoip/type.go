package geoip

type MyIp struct {
	Ip      string `json:"ip"`
	Country string `json:"country"`
	CC      string `json:"country,omitempty"`
	Org     string `json:"org,omitempty"`
}

type Countries struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Region string `json:"region"`
}

type GeoIpJson struct {
	Ip          string `json:"ip,omitempty"`
	CountryName string `json:"country_name,omitempty"`
	CountryCode string `json:"country,omitempty"`
	Region      string `json:"region,omitempty"`
	Org         string `json:"org,omitempty"`
}
