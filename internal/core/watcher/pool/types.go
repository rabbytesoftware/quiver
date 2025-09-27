package pool

import "github.com/sirupsen/logrus"

type Message struct {
	Level   logrus.Level
	Message string
}

type Subscriber func(
	logrus.Level,
	string,
)
