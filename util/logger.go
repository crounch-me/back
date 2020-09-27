package util

import (
	"time"

	"github.com/sirupsen/logrus"
)

func GetLogger() *logrus.Entry {
	return logrus.WithTime(time.Now())
}
