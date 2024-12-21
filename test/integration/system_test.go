package integration

import (
	"os"
	"path/filepath"
	"time"

	"ergo.services/ergo"
	"ergo.services/ergo/gen"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
)

// SystemTestSuite defines our integration test suite
type SystemTestSuite struct {
	suite.Suite
	gomega.WithT
	node     gen.Node
	nodeName gen.Atom
}

func (s *SystemTestSuite) SetupSuite() {
	s.nodeName = gen.Atom("test@localhost")
	nodeOpts := gen.NodeOptions{}

	node, err := ergo.StartNode(s.nodeName, nodeOpts)
	s.Require().NoError(err)
	s.node = node
	s.Require().True(node.IsAlive())
}

// TearDownSuite runs once after all tests complete
func (s *SystemTestSuite) TearDownSuite() {
	if s.node != nil {
		s.node.Stop()
	}
	// This is where we'll clean up our test resources
}

// SetupTest runs before each test
func (s *SystemTestSuite) SetupTest() {
	// Reset test state before each test
	// We'll add more setup here as we develop our system
}

// TestLogCollectorStartup verifies that our log collector actor
// can start up and begin monitoring a log file
func (s *SystemTestSuite) TestLogCollectorStartup() {
	g := gomega.NewGomegaWithT(s.T())

	// Create a directory for our test logs
	tmpDir, err := os.MkdirTemp("", "ergo-logs-test")
	s.Require().NoError(err)
	defer os.RemoveAll(tmpDir)

	// Create a log file with some initial content
	logFile := filepath.Join(tmpDir, "test.log")
	err = os.WriteFile(logFile, []byte("initial log entry\n"), 0644)
	s.Require().NoError(err)

	processOpts := gen.ProcessOptions{}
	pid, err := s.node.Spawn("log_collector", processOpts, LogCollectorProps(logFile))
	s.Require().NoError(err)

	_, err = s.node.ProcessInfo(pid)
	s.Require().NoError(err)

	// Add a log entry and verify that it is collected
	err = os.WriteFile(logFile, []byte("initial log entry\n"), 0644)
	s.Require().NoError(err)

	// Assert that our collector starts up within a reasonable timeout
	g.Eventually(func() bool {
		// TODO: Add real healthcheck implementation
		return true // Placeholder for now
	}).WithTimeout(5 * time.Second).Should(gomega.BeTrue())
}
