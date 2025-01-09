package collector

import (
	"fmt"
"regexp"
"strings"
"time"

	"github.com/jeroensoeters/ergo-logs/internal/messages"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
	"github.com/fsnotify/fsnotify"
)

type LogCollector struct {
	act.Actor
	filepath string
	watcher  *fsnotify.Watcher
}

func New() gen.ProcessBehavior {
	return &LogCollector{}
}

func (c *LogCollector) Init(args ...any) error {
	c.Log().Info("log collector started with args %w", args)

	return nil
}

func (c *LogCollector) HandleMessage(message any) error {
	switch msg := message.(type) {
	case messages.FileContentMessage:
		entries := c.parseSyslogEntries(msg.Content)
		if len(entries) > 0 {
			event := messages.SyslogEntriesEvent{
				Entries: entries,
			}
			c.Send(c.Self(), event)
		}
		return nil
	default:
		return fmt.Errorf("unknown message type: %T", message)
	}
}

func (c *LogCollector) parseSyslogEntries(content string) []messages.SyslogEntry {
	var entries []messages.SyslogEntry
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		if entry, ok := c.parseSyslogLine(line); ok {
			entries = append(entries, entry)
		}
	}

	return entries
}

func (c *LogCollector) parseSyslogLine(line string) (messages.SyslogEntry, bool) {
	// Basic syslog format: <time> <hostname> <tag>: <content>
	// Example: Jan 23 06:25:43 localhost nginx[1234]: GET /index.html
	pattern := `^(\w+\s+\d+\s+\d+:\d+:\d+)\s+(\S+)\s+([^:]+):\s+(.+)$`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return messages.SyslogEntry{}, false
	}

	// Parse timestamp
	timestamp, err := time.Parse("Jan 2 15:04:05", matches[1])
	if err != nil {
		c.Log().Error("failed to parse timestamp: %v", err)
		return messages.SyslogEntry{}, false
	}

	// If parsing succeeded, create and return the entry
	return messages.SyslogEntry{
		Timestamp: timestamp,
		Hostname:  matches[2],
		Tag:       matches[3],
		Content:   matches[4],
	}, true
}
