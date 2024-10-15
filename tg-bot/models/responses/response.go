package responses

import "log"

func HandlePanic() {
	if r := recover(); r != nil {
		log.Printf("Panic: %v", r)
	}
}
