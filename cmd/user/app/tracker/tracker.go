package tracker

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

type requestInfo struct {
	metadata  map[string]interface{}
	startTime time.Time
}

type RequestTracker struct {
	activeRequests sync.WaitGroup
	infoMap        *sync.Map
	logger         *logrus.Logger
	shuttingDown   atomic.Bool
}

func NewRequestTracker(logger *logrus.Logger) *RequestTracker {
	return &RequestTracker{
		infoMap: &sync.Map{},
		logger:  logger,
	}
}

func (rt *RequestTracker) Begin(requestId string, info map[string]any) {
	rt.activeRequests.Add(1)
	rt.infoMap.Store(requestId, requestInfo{
		metadata:  info,
		startTime: time.Now(),
	})
}

func (rt *RequestTracker) End(requestId string) {
	rt.activeRequests.Done()
	rt.infoMap.Delete(requestId)
}

func (rt *RequestTracker) IsShuttingDown() bool {
	return rt.shuttingDown.Load()
}

func (rt *RequestTracker) SetShuttingDown(value bool) {
	rt.shuttingDown.Store(value)
}

func (rt *RequestTracker) WaitForCompletion(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		rt.activeRequests.Wait()
		close(done)
	}()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return nil
		case <-ctx.Done():
			rt.logActiveRequests()
		case <-ticker.C:
			rt.logActiveRequests()
		}
	}
}

func (rt *RequestTracker) logActiveRequests() {
	var longRunning []string

	rt.infoMap.Range(func(key, value interface{}) bool {
		info := value.(requestInfo)
		duration := time.Since(info.startTime)

		if duration > 5*time.Second {
			details := fmt.Sprintf("ID: %v, duration: %v", key, duration.Round(time.Microsecond))

			if path, ok := info.metadata["path"]; ok {
				details += fmt.Sprintf(", path: %v", path)
			}
			if method, ok := info.metadata["method"]; ok {
				details += fmt.Sprintf(", method: %v", method)
			}

			longRunning = append(longRunning, details)
		}

		return true
	})

	if len(longRunning) > 0 {
		rt.logger.Warn("Long-running requests during shutdown", "count", len(longRunning), "long-running", longRunning)
	}
}
