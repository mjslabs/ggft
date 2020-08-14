package utils

import "log"

// CheckError - Exit on error
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
