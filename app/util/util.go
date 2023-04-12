package util

import (
	"time"

	"github.com/gofrs/uuid"
)

var (
	timeNow = time.Now
)

// GetUUIDv4 returns a UUIDv4
func GetUUIDv4() (string, error) {
	u, err := uuid.NewV4()
	return u.String(), err
}

func TimeNow() int64 {
	return timeNow().Unix()
}

func TimeNowMs() int64 {
	t := timeNow()
	return t.UnixNano() / time.Millisecond.Nanoseconds()
}
