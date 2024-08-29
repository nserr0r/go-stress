package main

import (
	"crypto/tls"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketManager struct {
	config               Config
	logger               Logger
	activeConnections    int64
	completedConnections int64
	proxyManager         *ProxyManager
}

func NewWebSocketManager(config Config, logger Logger, proxyManager *ProxyManager) *WebSocketManager {
	return &WebSocketManager{
		config:       config,
		logger:       logger,
		proxyManager: proxyManager,
	}
}

func (m *WebSocketManager) ManageConnection(connID int) {
	scheme := "ws"
	if m.config.UseSSL {
		scheme = "wss"
	}

	url := fmt.Sprintf("%s://%s%s", scheme, m.config.Host, m.config.Path)

	for {
		var dialer *websocket.Dialer
		if m.proxyManager != nil {
			proxyInfo := m.proxyManager.getNextProxy()
			if proxyInfo == nil {
				m.logger.LogError("No available proxies.")
				return
			}
			dialer = m.proxyManager.createDialerWithProxy(proxyInfo)
			if dialer == nil {
				m.proxyManager.markProxyAsDead(proxyInfo)
				continue
			}
		} else {
			dialer = websocket.DefaultDialer
		}

		// Add TLS configuration if using SSL
		if m.config.UseSSL {
			dialer.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: m.config.InsecureSkipVerify, // Опция для отключения проверки сертификата
			}
		}

		conn, _, err := dialer.Dial(url, nil)
		if err != nil {
			m.logger.LogError(fmt.Sprintf("Connection %d: failed to connect: %v", connID, err))
			time.Sleep(5000 * time.Millisecond)
			continue
		}

		atomic.AddInt64(&m.activeConnections, 1)
		m.logger.LogInfo(fmt.Sprintf("Connection %d: successfully connected", connID))

		lifetimeTimer := time.NewTimer(time.Duration(m.config.ConnLifetimeMs) * time.Millisecond)
		done := make(chan struct{})

		go func() {
			defer close(done)
			if m.config.Messages > 0 {
				for i := 0; i < m.config.Messages; i++ {
					msg := m.generateOrUseMessage(connID)
					if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
						m.logger.LogError(fmt.Sprintf("Connection %d: failed to send message: %v", connID, err))
						return
					}

					if _, _, err := conn.ReadMessage(); err != nil {
						m.logger.LogError(fmt.Sprintf("Connection %d: failed to read message: %v", connID, err))
						return
					}

					time.Sleep(time.Duration(m.config.DelayMs) * time.Millisecond)
				}

				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					m.logger.LogError(fmt.Sprintf("Connection %d: lost connection, reconnecting...", connID))
					return
				}
			} else {
				for {
					time.Sleep(30000 * time.Millisecond)
					if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
						m.logger.LogError(fmt.Sprintf("Connection %d: connection lost, reconnecting...", connID))
						return
					}
				}
			}
		}()

		select {
		case <-lifetimeTimer.C:
			m.logger.LogInfo(fmt.Sprintf("Connection %d: connection lifetime exceeded, reconnecting...", connID))
			conn.Close()
			atomic.AddInt64(&m.completedConnections, 1)
			atomic.AddInt64(&m.activeConnections, -1)
		case <-done:
			conn.Close()
			atomic.AddInt64(&m.completedConnections, 1)
			atomic.AddInt64(&m.activeConnections, -1)
		}
	}
}

func (m *WebSocketManager) Status() (active int64, completed int64) {
	return atomic.LoadInt64(&m.activeConnections), atomic.LoadInt64(&m.completedConnections)
}

func (m *WebSocketManager) MaxConnections() int {
	return m.config.Connections
}

func (m *WebSocketManager) generateOrUseMessage(connID int) string {
	if m.config.GenMsg {
		return fmt.Sprintf("Random message from connection %d", connID)
	}
	if m.config.CustomMsg != "" {
		return m.config.CustomMsg
	}
	return fmt.Sprintf("Hello from connection %d", connID)
}
