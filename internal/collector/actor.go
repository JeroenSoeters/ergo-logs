package collector

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
	"github.com/fsnotify/fsnotify"
)

type Props struct {
	LogFile string
}

type LogCollector struct {
	act.Actor
	filepath string
	watcher  *fsnotify.Watcher
}

func NewLogCollector() gen.ProcessBehavior {
	return &LogCollector{}
}
