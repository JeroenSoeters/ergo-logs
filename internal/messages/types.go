package messages

import (
	"ergo.services/ergo/net/edf"
)

type FileContentMessage struct {
	Path    string
	Content string
}

func init() {
	// register network messages
	if err := edf.RegisterTypeOf(FileContentMessage{}); err != nil {
		panic(err)
	}
}
