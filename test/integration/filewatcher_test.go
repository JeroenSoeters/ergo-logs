package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

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
	node     gen.Node
	nodeName gen.Atom
}

type TestReceiver struct {
	act.Actor
	messages []filewatcher.FileContentMessage
}

func (r *TestReceiver) Init(args ...any) error {
	return r.Actor.Init(args...)
}

func NewTestReceiver() *TestReceiver {
	return &TestReceiver{}
}

func (s *FileWatcherTestSuite) SetupSuite() {
	s.nodeName = gen.Atom("test@localhost")
	nodeOpts := gen.NodeOptions{}

	node, err := ergo.StartNode(s.nodeName, nodeOpts)
	s.Require().NoError(err)
	s.node = node
}

func (s *FileWatcherTestSuite) TearDownSuite() {
	if s.node != nil {
		fmt.Println("stopping node")
		s.node.Stop()
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

	fmt.Println("log file " + logFile)

	// Create a test test test receiver
	fmt.Println("[TEST] creating test receiver")
	receiver := NewTestReceiver()
	_, err = s.node.SpawnRegister(gen.Atom("log_processor"), func() gen.ProcessBehavior { return receiver }, gen.ProcessOptions{})
	s.Require().NoError(err)
	fmt.Println("[TEST] test receiver created")

	fmt.Println("[TEST] creating watcher")
	_, err = s.node.Spawn(filewatcher.New, gen.ProcessOptions{}, logFile)
	s.Require().NoError(err)
	fmt.Println("[TEST] watcher created")

	time.Sleep(100 * time.Millisecond)

	// Add a log entry and verify that it is collected
	fmt.Println("[TEST] creating log entry")
	logEntry := "another log entry\n"
	err = os.WriteFile(logFile, []byte(logEntry), 0644)
	s.Require().NoError(err)
	fmt.Println("[TEST] log entry created")

	g.Eventually(func() []filewatcher.FileContentMessage {
		return receiver.messages
	}, 5 * time.Second).Should(gomega.HaveLen(1))
	s.Equal(logEntry, receiver.messages[0].Content)
	s.Equal(logFile, receiver.messages[0].Path)

	s.node.Wait()
}
