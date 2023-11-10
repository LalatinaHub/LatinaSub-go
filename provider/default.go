package provider

import (
	"encoding/base64"
	"os"
	"regexp"
	"strings"
)

func GetFirstLine(content string) (string, string) {
	lines := strings.Split(content, "\n")
	if len(lines) == 1 {
		return lines[0], ""
	}
	others := strings.Join(lines[1:], "\n")
	return lines[0], others
}

func GetTrimedFile(path string) []byte {
	content, _ := os.ReadFile(path)
	return []byte(TrimBlank(string(content)))
}

func TrimBlank(str string) string {
	str = strings.Trim(str, " ")
	str = strings.Trim(str, "\a")
	str = strings.Trim(str, "\b")
	str = strings.Trim(str, "\f")
	str = strings.Trim(str, "\r")
	str = strings.Trim(str, "\t")
	str = strings.Trim(str, "\v")
	return str
}

func replaceIllegalBase64(content string) string {
	result := content
	result = strings.ReplaceAll(result, "-", "+")
	result = strings.ReplaceAll(result, "_", "/")
	return result
}

func DecodeBase64Safe(content string) string {
	reg1 := regexp.MustCompile(`^(?:[A-Za-z0-9-_+/]{4})*[A-Za-z0-9_+/]{4}$`)
	reg2 := regexp.MustCompile(`^(?:[A-Za-z0-9-_+/]{4})*[A-Za-z0-9_+/]{3}(=)?$`)
	reg3 := regexp.MustCompile(`^(?:[A-Za-z0-9-_+/]{4})*[A-Za-z0-9_+/]{2}(==)?$`)
	var result []string
	result = reg1.FindStringSubmatch(content)
	if len(result) > 0 {
		decode, err := base64.StdEncoding.DecodeString(replaceIllegalBase64(content))
		if err == nil {
			return string(decode)
		}
	}
	result = reg2.FindStringSubmatch(content)
	if len(result) > 0 {
		equals := ""
		if result[1] == "" {
			equals = "="
		}
		decode, err := base64.StdEncoding.DecodeString(replaceIllegalBase64(content + equals))
		if err == nil {
			return string(decode)
		}
	}
	result = reg3.FindStringSubmatch(content)
	if len(result) > 0 {
		equals := ""
		if result[1] == "" {
			equals = "=="
		}
		decode, err := base64.StdEncoding.DecodeString(replaceIllegalBase64(content + equals))
		if err == nil {
			return string(decode)
		}
	}
	return content
}

func ParseContent(contentRaw string) string {
	content := DecodeBase64Safe(contentRaw)
	return content
}
