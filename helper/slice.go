package helper

func SliceContains(slice []string, find string) bool {
	for _, value := range slice {
		if value == find {
			return true
		}
	}

	return false
}

func FilterEmptyString(slice []string) []string {
	var cleanString []string

	for _, x := range slice {
		if x != "" {
			cleanString = append(cleanString, x)
		}
	}

	return cleanString
}

func FilterDuplicateString(slice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range slice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
