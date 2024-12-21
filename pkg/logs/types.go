package logs

import "time"

type LogEntry struct {
	Timestamp time.Time
	Message   string
}
