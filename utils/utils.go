package utils

import (
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
)

// UID новый идентификатор
func UID() string {
	return uuid.NewV5(uuid.NewV4(), baseUID).String()
}

// baseUID базовый UID
var baseUID string

// getHostUID возвращает идентификатор хоста
func getHostUID() string {
	hostname, err := os.Hostname()

	if err != nil {
		return time.Now().Format(time.RFC3339Nano)
	}

	return hostname
}

func init() {
	baseUID = getHostUID()
}
