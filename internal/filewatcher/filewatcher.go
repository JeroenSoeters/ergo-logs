package filewatcher

import (
	"errors"
	"fmt"
	"os"

	"github.com/jeroensoeters/ergo-logs/internal/messages"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
	"github.com/fsnotify/fsnotify"
)

type FileWatcher struct {
	act.Actor

	filepath      string
	watcher       *fsnotify.Watcher
	offset        int64
	processorName gen.Atom
}

func New() gen.ProcessBehavior {
	return &FileWatcher{
		processorName: gen.Atom("log_processor"),
	}
}

func (w *FileWatcher) Init(args ...any) error {
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

	err = watcher.Add(filepath)
	if err != nil {
		return fmt.Errorf("failed to watch path %s: %w", filepath, err)
	}

	w.filepath = filepath
	w.watcher = watcher

	go w.watchFileEvents()

	w.Log().Info("started FileWatcher process with args %v", args)

	return nil
}

func (w *FileWatcher) watchFileEvents() {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			if event.Has(fsnotify.Write) {
				file, err := os.Open(w.filepath)
				if err != nil {
					w.Log().Warning("error opening file %s: %w", w.filepath, err)
					//TODO: handle errors
					continue
				}
				defer file.Close()

				fileInfo, err := file.Stat()
				if err != nil {
					w.Log().Warning("error stat-ing file %s: %w", w.filepath, err)
					//TODO: handle Errors
					continue
				}

				if fileInfo.Size() > w.offset {
					content := make([]byte, fileInfo.Size()-w.offset)
					_, err := file.ReadAt(content, w.offset)
					if err != nil {
						//TODO: handle Errors
						continue
					}

					msg := messages.FileContentMessage{
						Content: string(content),
						Path:    w.filepath,
					}

					procID := gen.ProcessID{
						Name: w.processorName,
						Node: gen.Atom("server@localhost"), //TODO: look up remote node name
					}

					err = w.Send(procID, msg)
					if err != nil {
						fmt.Printf("Error sending message: %s", err)
						//TODO: handle errors
						continue
					}

					w.offset = fileInfo.Size()
				}
			}
		case _, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			//TODO: handle errors
		}
	}
}

func (w *FileWatcher) ProcessTerminate(reason error) {
	if w.watcher != nil {
		w.watcher.Close()
	}
}
