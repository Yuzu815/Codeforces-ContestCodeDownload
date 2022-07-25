package logMode

import "github.com/sirupsen/logrus"

type prefixHook struct {
	UID string
}

func (h *prefixHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *prefixHook) Fire(entry *logrus.Entry) error {
	entry.Data["UID"] = h.UID
	return nil
}
