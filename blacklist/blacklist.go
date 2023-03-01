package blacklist

import (
	"log"
	"os"
	"strings"
)

var (
	BlacklistPath string = "blacklist/list/"
	BlacklistFile string = BlacklistPath + "nodes"
	BlacklistByte []byte
)

func Init() {
	// Check and create dir "blacklist/list/"
	if _, err := os.Stat(BlacklistPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(BlacklistPath, os.ModePerm)
		} else {
			log.Panic(err)
		}
	}

	// Check and create nodes file
	if _, err := os.Stat(BlacklistFile); err != nil {
		if os.IsNotExist(err) {
			file, _ := os.Create(BlacklistFile)
			defer file.Close()
		}
	}

	// Read blacklist file
	BlacklistByte, _ = os.ReadFile(BlacklistFile)
}

func Find(value string) bool {
	for _, blacklisted := range strings.Split(string(BlacklistByte), "\n") {
		if blacklisted == value {
			return true
		}
	}

	return false
}

func Save(value string) {
	if len(BlacklistByte) > 0 {
		BlacklistByte = append(BlacklistByte, []byte("\n"+value)...)
	} else {
		BlacklistByte = []byte(value)
	}
}

func Write() {
	os.WriteFile(BlacklistFile, BlacklistByte, os.ModePerm)
}
