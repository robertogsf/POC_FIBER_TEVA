package tools

import (
	"log"
)

func IsEmpty(par string, name string) {
	if par == "" {
		log.Fatalf("%s no está definida", name)
	}
}
