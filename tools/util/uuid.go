package util

import "github.com/google/uuid"

// CreateUUID create uuid
func CreateUUID() string {
	return "cache-" + uuid.NewString()
}
