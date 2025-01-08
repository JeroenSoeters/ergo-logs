package integration

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jeroensoeters/ergo-logs/internal/messages"
	"github.com/jeroensoeters/ergo-logs/internal/filewatcher"

	"ergo.services/ergo"
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
)

type FileWatcherTestSuite struct {
	suite.Suite
	gomega.WithT
	clientNode, serverNode     gen.Node
}

type TestReceiver struct {
	act.Actor
	messages []messages.FileContentMessage
}

func (r *TestReceiver) Init(args ...any) error {
	return r.Actor.Init(args...)
}

func (r *TestReceiver) HandleMessage(from gen.PID, message any) error {
	switch msg := message.(type) {
	case messages.FileContentMessage:
		r.messages = append(r.messages, msg)
	}
	return nil
}

func NewTestReceiver() *TestReceiver {
	return &TestReceiver{}
}

func (s *FileWatcherTestSuite) SetupSuite() {
	// start client (log shipper)
	nodeOpts := gen.NodeOptions{}
	nodeOpts.Network.Cookie = "123"
	client, err := ergo.StartNode(gen.Atom("client@localhost"), nodeOpts)
	s.Require().NoError(err)
	s.clientNode = client

	// start server (backend)
	server, err := ergo.StartNode(gen.Atom("server@localhost"), nodeOpts)
	s.Require().NoError(err)
	s.serverNode = server
}

func (s *FileWatcherTestSuite) TearDownSuite() {
	if s.clientNode != nil {
		s.clientNode.Stop()
	}
	if s.serverNode != nil {
		s.serverNode.Stop()
	}
}

func (s *FileWatcherTestSuite) SetupTest() {
	// Reset test state before each test
	// We'll add more setup here as we develop our system
}

func TestFileWatcherTestSuite(t *testing.T) {
	suite.Run(t, new(FileWatcherTestSuite))
}

func (s *FileWatcherTestSuite) TestFileWatcher() {
	g := gomega.NewGomegaWithT(s.T())

	// Create a directory for our test logs
	tmpDir, err := os.MkdirTemp("", "ergo-logs-test")
	s.Require().NoError(err)
	defer os.RemoveAll(tmpDir)

	// Create a log file with some initial content
	logFile := filepath.Join(tmpDir, "test.log")
	err = os.WriteFile(logFile, []byte("initial log entry\n"), 0644)
	s.Require().NoError(err)

	// Create a test receiver
	receiver := NewTestReceiver()
	_, err = s.serverNode.SpawnRegister("log_processor", func() gen.ProcessBehavior { return receiver }, gen.ProcessOptions{})
	s.Require().NoError(err)

	_, err = s.clientNode.Spawn(filewatcher.New, gen.ProcessOptions{}, logFile)
	s.Require().NoError(err)

	time.Sleep(100 * time.Millisecond)

	// Add a log entry and verify that it is collected
	logEntry := "another log entry\n"
	err = os.WriteFile(logFile, []byte(logEntry), 0644)
	s.Require().NoError(err)

	g.Eventually(func() []messages.FileContentMessage {
		return receiver.messages
	}, 5*time.Second).Should(gomega.HaveLen(1))
	s.Equal(logEntry, receiver.messages[0].Content)
	s.Equal(logFile, receiver.messages[0].Path)
}
