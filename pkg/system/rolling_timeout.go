package system

import (
	"context"
	"fmt"
	"time"
)

// wait for a time period - then call a "check" function
// if the check function returns true it means
// we are still active and should carry on waiting
// so start another wait for the configured time period
// if it returns false it means we have timed out
// example: watch the size of this folder
// if it has not grown by > X then timeout
type RollingTimeoutHandler[T, V any] interface {
	// return the initial value
	GetInitialValue() T
	// the function that is called each time to check
	// if we have made progress - we pass it the previous
	// value that was returned by the check function
	CheckFunction(prevValue T) (T, bool, error)
	// the function that will be doing work
	// it is a cancel context based on the context
	// passed into the rolling timout
	WorkFunction(ctx context.Context) (V, error)
}

func RollingTimeout[T, V any](
	handler RollingTimeoutHandler[T, V],
	ctx context.Context,
	delay time.Duration,
) (*V, error) {
	// get the initial value to pass into the check function
	previousValue := handler.GetInitialValue()
	// we need a way of canceling the work if we
	// decide we have timed out
	ctxWithCancel, cancelFunction := context.WithCancel(ctx)
	// the work function has completed
	resultChan := make(chan V, 1)
	// the work function has errored
	errorChan := make(chan error, 1)
	// we tick every "delay" time period
	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	go func() {
		innerResult, err := handler.WorkFunction(ctxWithCancel)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- innerResult
		}
	}()

	for {
		select {
		case <-ticker.C:
			// we have passed the delay and not heard back from the work function
			// so let's call our check with the previuos value to see if we
			// are really stuck or not
			innerPreviousValue, isAlive, err := handler.CheckFunction(previousValue)
			if err != nil {
				// if the check function has errored then all bets are off
				errorChan <- err
			} else if !isAlive {
				// if we are not alive then we have timed out
				errorChan <- fmt.Errorf("we have timed out")
			} else {
				// we are still alive so we need to reset the previous value
				// and carry on waiting
				previousValue = innerPreviousValue
			}
		case result := <-resultChan:
			cancelFunction()
			return &result, nil
		case err := <-errorChan:
			cancelFunction()
			return nil, err
		}
	}
}
