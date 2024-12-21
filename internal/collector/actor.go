package collector

import (
	"ergo.services/ergo/gen"
	"github.com/fsnotify/fsnotify"
)

type LogCollector struct {
	gen.Server
	filepath    string
	watcher     *fsnotify.Watcher
	subscribers []etf.Pid
}
