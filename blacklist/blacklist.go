package blacklist

var (
	Blacklist []string
)

func Find(uid string) bool {
	for _, blacklistId := range Blacklist {
		if blacklistId == uid {
			return true
		}
	}

	return false
}

func Save(uid string) {
	Blacklist = append(Blacklist, uid)
}

func Clear() {
	Blacklist = []string{}
}
