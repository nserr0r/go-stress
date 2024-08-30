package status

import (
	"fmt"
	"sync/atomic"
	"time"
)

const (
	colorRed        = "\033[31m"
	colorGreen      = "\033[32m"
	colorYellow     = "\033[33m"
	colorLightGreen = "\033[92m"
	colorReset      = "\033[0m"
)

type ProxyReporter interface {
	GetProxyStatusString() string
}

type StatusManager struct {
	activeConnections    *int64
	completedConnections *int64
	maxConnections       int
	proxyReporter        ProxyReporter
}

func NewStatusManager(activeConn, completedConn *int64, maxConn int, proxyReporter ProxyReporter) *StatusManager {
	return &StatusManager{
		activeConnections:    activeConn,
		completedConnections: completedConn,
		maxConnections:       maxConn,
		proxyReporter:        proxyReporter,
	}
}

func (s *StatusManager) DisplayStatus() {
	for {
		active := atomic.LoadInt64(s.activeConnections)
		completed := atomic.LoadInt64(s.completedConnections)
		maxConnectionsDigits := len(fmt.Sprintf("%d", s.maxConnections))
		statusString := fmt.Sprintf("\r%sActive Connections: %s%0*d %s| Completed Connections: %s%0*d%s",
			colorYellow, colorRed, maxConnectionsDigits, active, colorYellow, colorGreen, maxConnectionsDigits, completed, colorReset)

		if s.proxyReporter != nil {
			proxyStatus := s.proxyReporter.GetProxyStatusString()
			statusString += proxyStatus
		}

		fmt.Print(statusString)
		time.Sleep(40 * time.Millisecond)
	}
}

