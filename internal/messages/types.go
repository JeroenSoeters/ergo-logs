package messages

import (
	"ergo.services/ergo/net/edf"
	"time"
)

type FileContentMessage struct {
	Path    string
	Content string
}

type SyslogEntry struct {
	Timestamp time.Time
	Hostname  string
	Tag       string
	Content   string
}

type SyslogEntriesEvent struct {
	Entries []SyslogEntry
}

func init() {
	// register network messages
	if err := edf.RegisterTypeOf(FileContentMessage{}); err != nil {
		panic(err)
	}
	if err := edf.RegisterTypeOf(SyslogEntriesEvent{}); err != nil {
		panic(err)
	}
}
