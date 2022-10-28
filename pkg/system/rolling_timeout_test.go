package system

import (
	"context"
	"testing"
	"time"

	"github.com/filecoin-project/bacalhau/pkg/logger"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RollingTimeoutSuite struct {
	suite.Suite
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestRollingTimeoutSuite(t *testing.T) {
	suite.Run(t, new(RollingTimeoutSuite))
}

// Before all suite
func (suite *RollingTimeoutSuite) SetupAllSuite() {

}

// Before each test
func (suite *RollingTimeoutSuite) SetupTest() {
	logger.ConfigureTestLogging(suite.T())
	require.NoError(suite.T(), InitConfigForTesting())
}

func (suite *RollingTimeoutSuite) TearDownTest() {
}

func (suite *RollingTimeoutSuite) TearDownAllSuite() {

}

type RollingTimeoutTester struct {
	initialValue int
	workDelay    time.Duration
	getNextValue func(prevValue int) int
	isAlive      func(prevValue int) bool
	getError     func(prevValue int) error
}

func (tester *RollingTimeoutTester) GetInitialValue() int {
	return tester.initialValue
}

func (tester *RollingTimeoutTester) CheckFunction(prevValue int) (int, bool, error) {
	return tester.getNextValue(prevValue), tester.isAlive(prevValue), tester.getError(prevValue)
}

func (tester *RollingTimeoutTester) WorkFunction(ctx context.Context) (string, error) {
	time.Sleep(tester.workDelay)
	return "hello world", nil
}

func TestSanity(t *testing.T) {
	handler := RollingTimeoutHandler{
		initialValue: 0,
		workDelay:    time.Millisecond * 100,
		getNextValue: func(prevValue int) int {
			return prevValue + 1
		},
		isAlive: func(prevValue int) bool {
			return true
		},
		getError: func(prevValue int) error {
			return nil
		},
	}
}
