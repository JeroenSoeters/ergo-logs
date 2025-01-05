package filewatcher

import (
	"errors"
	"fmt"
	"os"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
	"github.com/fsnotify/fsnotify"
)

type FileContentMessage struct {
	Path    string
	Content string
}

type FileWatcher struct {
	act.Actor
	filepath      string
	watcher       *fsnotify.Watcher
	offset        int64
	processorName gen.Atom
}

func New() gen.ProcessBehavior {
	return &FileWatcher{
		offset:        0,
		processorName: gen.Atom("log_processor"),
	}
}

func (w *FileWatcher) ProcessInit(process gen.Process, args ...any) error {
	fmt.Println("FileWatcher initializing")
	if len(args) < 1 {
		return errors.New("missing required filepath argument")
	}

	filepath, ok := args[0].(string)
	if !ok {
		return errors.New("filepath must be a string")
	}

	fmt.Println("FileWatcher starting")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		watcher.Close()
		return fmt.Errorf("failed to create watcher: %w", err)
	}

	w.filepath = filepath
	w.watcher = watcher

	fmt.Println("FileWatcher watching", w.filepath)

	go w.watchFileEvents()

	return nil
}

func (w *FileWatcher) watchFileEvents() {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			fmt.Println("FileWatcher event", event)
			if !ok {
				return
			}

			if event.Has(fsnotify.Write) {
				file, err := os.Open(w.filepath)
				if err != nil {
					//TODO: handle errors
					continue
				}
				defer file.Close()

				fileInfo, err := file.Stat()
				if err != nil {
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

					msg := FileContentMessage{
						Content: string(content),
						Path:    w.filepath,
					}

					procID := gen.ProcessID{
						Name: w.processorName,
						Node: w.Actor.Node().Name(), //TODO: look up remote node name
					}

					err = w.SendProcessID(procID, msg)
					if err != nil {
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
	fmt.Println("FileWatcher terminating")
	if w.watcher != nil {
		w.watcher.Close()
	}
}
