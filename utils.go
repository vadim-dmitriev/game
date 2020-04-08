package main

import (
	"github.com/google/uuid"
)

func newUUID() string {
	uid, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	return uid.String()
}
