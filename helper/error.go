package helper

import "fmt"

func CatchError(print bool) interface{} {
	message := recover()

	if message != nil && print {
		fmt.Println("[-] Error:", message)
	}
	return message
}
