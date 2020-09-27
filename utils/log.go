package utils

import "github.com/sirupsen/logrus"

// NewLogger создает логер
func NewLogger(id string) *logrus.Entry {
	return logrus.NewEntry(logrus.StandardLogger()).WithField("id", id)
}
