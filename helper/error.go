package helper

import (
	"fmt"
	"runtime/debug"
)

func CatchError(print bool) interface{} {
	message := recover()

	if message != nil && print {
		fmt.Println("[-] Error:", message)
		fmt.Println("[-] Stack:", string(debug.Stack()))
	}
	return message
}
