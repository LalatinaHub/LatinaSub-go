package subscription

type SubStruct struct {
	Id            int    `json:"id"`
	Remarks       string `json:"remarks"`
	Site          string `json:"site"`
	Url           string `json:"url"`
	Update_method string `json:"update_method"`
	Enabled       bool   `json:"enabled"`
}
