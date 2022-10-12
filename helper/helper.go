package helper

import (
	"log"
)

func FatalOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
