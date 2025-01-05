package collector

import (
	"errors"
	"fmt"

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
	if err := a.Actor.Init(args); err != nil {
		return err
	}

	if len(args) < 1 {
		return errors.New("missing required filepath argument")
	}

	filepath, ok := args[0].(string)
	if !ok {
		return errors.New("filepath must be a string")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		watcher.Close()
		return fmt.Errorf("failed to create watcher: %w", err)
	}

	if err := watcher.Add(filepath); err != nil {
		watcher.Close()
		return fmt.Errorf("failed to watch file: %w", err)
	}

	c.filepath = filepath
	c.watcher = watcher

	// Start watching for file events
	go c.watchFileEvents(process)

	return nil
}
