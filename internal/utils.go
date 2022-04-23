package internal

import (
	"errors"
	"log"
	"os"
)

func FileExists(path string) bool {
	_, err := os.Stat(path) //Stat returns metadata for said file but does not open it
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	Handle(err)
	return false
}


func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

