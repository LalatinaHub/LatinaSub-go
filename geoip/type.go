package geoip

type MyIp struct {
	Ip      string `json:"ip"`
	IpType  int    `json:"ip_type,omitempty"`
	Country string `json:"country"`
	CC      string `json:"country_abbr,omitempty"`
	Region  string `json:"continent,omitempty"`
	ASN     string `json:"asn,omitempty"`
	Org     string `json:"asn_org,omitempty"`
}

type GeoIpCity struct {
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
}

type GeoIpASN struct {
	ASNName string `maxminddb:"autonomous_system_organization"`
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
