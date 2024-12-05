package utils

import (
	"log"
	"strconv"
)

func ParseStringToUint(str string) uint {

	u64, err := strconv.ParseUint(str, 0, 64)

	if err != nil {
		log.Printf("Error parsing String to Uint : %v", err)
	}

	return uint(u64)
}
