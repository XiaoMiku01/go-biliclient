package utils

import "github.com/google/uuid"

func GetUUid() string {
	newUUID := uuid.New()
	return newUUID.String()
}
