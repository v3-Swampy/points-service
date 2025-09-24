package parsing

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

var logger = logrus.WithField("module", "parsing")

func logHourTimestamp(hourTimestamp int64) *logrus.Entry {
	dt := time.Unix(hourTimestamp, 0).Format(time.DateTime)

	return logger.WithField("ts", fmt.Sprintf("%v (%v)", dt, hourTimestamp))
}
