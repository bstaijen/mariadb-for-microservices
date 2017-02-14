package util

import (
	"log"
)

// PanicIfError panics when error is not nil
func PanicIfError(err error) {
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
}
