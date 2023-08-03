package util

import (
	"github.com/google/uuid"
	"time"
)

func GenerateTimestamp() int64 {
	return time.Now().UnixMilli()
}

func GenerateMsgID() string {
	return uuid.NewString()
}
