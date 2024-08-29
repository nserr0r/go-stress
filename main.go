package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	config := LoadConfig()
	logger := NewConsoleLogger(config.Log)

	var proxyManager *ProxyManager
	if config.ProxyFile != "" {
		proxyManager = NewProxyManager(config, logger)
	}

	connectionManager := NewWebSocketManager(config, logger, proxyManager)

	var wg sync.WaitGroup

	if !config.Log {
		go displayStatus(connectionManager, proxyManager)
	}

	for i := 0; i < config.Connections; i++ {
		wg.Add(1)
		go func(connID int) {
			defer wg.Done()
			connectionManager.ManageConnection(connID)
		}(i)

		time.Sleep(time.Duration(config.ConnDelayMs) * time.Millisecond)
	}

	wg.Wait()
	select {}
}

const (
	colorRed        = "\033[31m"
	colorGreen      = "\033[32m"
	colorYellow     = "\033[33m"
	colorLightGreen = "\033[92m"
	colorReset      = "\033[0m"
)

func displayStatus(manager ConnectionManager, proxyManager *ProxyManager) {
	for {
		active, completed := manager.Status()
		maxConnectionsDigits := len(fmt.Sprintf("%d", manager.MaxConnections()))
		statusString := fmt.Sprintf("\r%sActive Connections: %s%0*d %s| Completed Connections: %s%0*d%s",
			colorYellow, colorRed, maxConnectionsDigits, active, colorYellow, colorGreen, maxConnectionsDigits, completed, colorReset)

		if proxyManager != nil {
			proxyStatus := proxyManager.getProxyStatusString()
			statusString += proxyStatus
		}

		fmt.Print(statusString)
		time.Sleep(40 * time.Millisecond)
	}
}
