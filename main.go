package main

import (
	"sync"
	"time"

	"main/config"
	"main/connection"
	"main/logger"
	"main/proxy"
	"main/status"
)

func main() {
	cfg := config.LoadConfig()
	log := logger.NewConsoleLogger(cfg.Log)

	var proxyManager *proxy.ProxyManager
	var dialerFactory *proxy.DialerFactory

	if cfg.ProxyFile != "" {
		proxyManager = proxy.NewProxyManager(cfg.ProxyFile)
		dialerFactory = proxy.NewDialerFactory(cfg, log)
		dialerFactory.ProxyManager = proxyManager
	} else {
		dialerFactory = proxy.NewDialerFactory(cfg, log)
	}

	var connectionManager connection.ConnectionManager
	if cfg.UseWebSocket {
		connectionManager = connection.NewWebSocketManager(cfg, log, dialerFactory)
	} else {
		connectionManager = connection.NewHTTPManager(cfg, log, dialerFactory)
	}

	if !cfg.Log {
		activeConnections := connectionManager.GetActiveConnections()
		completedConnections := connectionManager.GetCompletedConnections()

		var statusManager *status.StatusManager
		if proxyManager != nil {
			statusManager = status.NewStatusManager(&activeConnections, &completedConnections, cfg.Connections, proxyManager)
		} else {
			statusManager = status.NewStatusManager(&activeConnections, &completedConnections, cfg.Connections, nil)
		}
		go statusManager.DisplayStatus()
	}

	var wg sync.WaitGroup
	for i := 0; i < cfg.Connections; i++ {
		wg.Add(1)
		go func(connID int) {
			defer wg.Done()
			connectionManager.ManageConnection(connID, cfg.Body)
		}(i)

		time.Sleep(time.Duration(cfg.ConnDelayMs) * time.Millisecond)
	}

	wg.Wait()
	select {}
}

