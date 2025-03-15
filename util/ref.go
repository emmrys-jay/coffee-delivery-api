package util

import "github.com/google/uuid"

func GenerateReference() string {
	return uuid.New().String()
}
