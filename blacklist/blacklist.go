package blacklist

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	BlacklistPath                                      string = "blacklist/list/"
	BlacklistNodeByte, BlacklistSubByte, BlacklistByte []byte
	BlacklistPtr                                       *[]byte
)

func Init() {
	// Check and create dir "blacklist/list/"
	if _, err := os.Stat(BlacklistPath); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(BlacklistPath, os.ModePerm)
		} else {
			log.Panic(err)
		}
	}

	// Check and create sub file
	if _, err := os.Stat(BlacklistPath + "sub"); err != nil {
		if os.IsNotExist(err) {
			file, _ := os.Create(BlacklistPath + "sub")
			defer file.Close()
		}
	}

	// Check and create node file
	if _, err := os.Stat(BlacklistPath + "node"); err != nil {
		if os.IsNotExist(err) {
			file, _ := os.Create(BlacklistPath + "node")
			defer file.Close()
		}
	}

	// Read blacklist file
	BlacklistNodeByte, _ = os.ReadFile(BlacklistPath + "node")
	BlacklistSubByte, _ = os.ReadFile(BlacklistPath + "sub")
}

func Find(blacklistType, value string) bool {
	// Assign pointer
	if blacklistType == "node" {
		BlacklistPtr = &BlacklistNodeByte
	} else {
		BlacklistPtr = &BlacklistSubByte
	}

	for _, blacklisted := range strings.Split(string(*BlacklistPtr), "\n") {
		if blacklisted == value {
			return true
		}
	}

	return false
}

func Save(blacklistType, value string) {
	if Find(blacklistType, value) {
		return
	}

	if len(*BlacklistPtr) > 0 {
		*BlacklistPtr = append(*BlacklistPtr, []byte("\n"+value)...)
	} else {
		*BlacklistPtr = []byte(value)
	}

	fmt.Println(fmt.Sprintf("[blacklist] %s blacklisted!", value))
}

func Write() {
	os.WriteFile(BlacklistPath+"node", BlacklistNodeByte, os.ModePerm)
	os.WriteFile(BlacklistPath+"sub", BlacklistSubByte, os.ModePerm)
}
