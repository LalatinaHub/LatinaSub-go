package helper

import (
	"encoding/json"
)

func IsJson(str string) bool {
	var js json.RawMessage
	json.Unmarshal([]byte(str), &js)

	return len(js) > 0
}
