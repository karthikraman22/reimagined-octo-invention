package util

import "github.com/google/uuid"

func ToUuid(sUuid string) *uuid.UUID {
	uid, err := uuid.Parse(sUuid)
	if err != nil {
		return nil
	}
	return &uid
}
