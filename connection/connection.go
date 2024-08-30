package connection

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"main/config"
	"main/logger"
	"main/proxy"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type ConnectionManager interface {
	ManageConnection(connID int, body string)
	GetActiveConnections() int64
	GetCompletedConnections() int64
}

type HTTPManager struct {
	config               *config.Config
	logger               logger.Logger
	dialerFactory        *proxy.DialerFactory
	activeConnections    int64
	completedConnections int64
}

func NewHTTPManager(config *config.Config, logger logger.Logger, dialerFactory *proxy.DialerFactory) *HTTPManager {
	return &HTTPManager{
		config:        config,
		logger:        logger,
		dialerFactory: dialerFactory,
	}
}

func (m *HTTPManager) ManageConnection(connID int, body string) {
	scheme := "http"
	if m.config.UseSSL {
		scheme = "https"
	}

	urlStr := fmt.Sprintf("%s://%s%s", scheme, m.config.Host, m.config.Path)

	for {
		var proxyURL *url.URL
		if m.config.ProxyFile != "" {
			proxyInfo := m.dialerFactory.ProxyManager.GetNextProxy()
			if proxyInfo != nil {
				proxyURL = proxyInfo.URL
			}
		}

		client := m.dialerFactory.CreateHTTPClient(proxyURL)

		var req *http.Request
		var err error

		if body != "" {
			req, err = http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(body)))
		} else {
			req, err = http.NewRequest("GET", urlStr, nil)
		}

		if err != nil {
			m.logger.LogError(fmt.Sprintf("Connection %d: failed to create request: %v", connID, err))
			time.Sleep(5000 * time.Millisecond)
			continue
		}

		// Set headers
		for key, value := range m.config.Headers {
			req.Header.Set(key, value)
		}

		atomic.AddInt64(&m.activeConnections, 1)

		resp, err := client.Do(req)
		if err != nil {
			m.logger.LogError(fmt.Sprintf("Connection %d: failed to perform request: %v", connID, err))
			atomic.AddInt64(&m.activeConnections, -1)
			time.Sleep(5000 * time.Millisecond)
			continue
		}

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			m.logger.LogError(fmt.Sprintf("Connection %d: failed to read response: %v", connID, err))
			resp.Body.Close()
			atomic.AddInt64(&m.activeConnections, -1)
			time.Sleep(5000 * time.Millisecond)
			continue
		}

		m.logger.LogInfo(fmt.Sprintf("Connection %d: received response: %s", connID, string(bodyBytes)))
		resp.Body.Close()

		atomic.AddInt64(&m.completedConnections, 1)
		atomic.AddInt64(&m.activeConnections, -1)

		time.Sleep(time.Duration(m.config.ConnLifetimeMs) * time.Millisecond)
	}
}

func (m *HTTPManager) GetActiveConnections() int64 {
	return atomic.LoadInt64(&m.activeConnections)
}

func (m *HTTPManager) GetCompletedConnections() int64 {
	return atomic.LoadInt64(&m.completedConnections)
}

type WebSocketManager struct {
	config               *config.Config
	logger               logger.Logger
	dialerFactory        *proxy.DialerFactory
	activeConnections    int64
	completedConnections int64
}

func NewWebSocketManager(config *config.Config, logger logger.Logger, dialerFactory *proxy.DialerFactory) *WebSocketManager {
	return &WebSocketManager{
		config:        config,
		logger:        logger,
		dialerFactory: dialerFactory,
	}
}

func (m *WebSocketManager) ManageConnection(connID int, body string) {
	scheme := "ws"
	if m.config.UseSSL {
		scheme = "wss"
	}

	urlStr := fmt.Sprintf("%s://%s%s", scheme, m.config.Host, m.config.Path)

	for {
		var proxyURL *url.URL
		if m.config.ProxyFile != "" {
			proxyInfo := m.dialerFactory.ProxyManager.GetNextProxy()
			if proxyInfo != nil {
				proxyURL = proxyInfo.URL
			}
		}

		dialer := m.dialerFactory.CreateWebSocketDialer(proxyURL)

		conn, _, err := dialer.Dial(urlStr, nil)
		if err != nil {
			m.logger.LogError(fmt.Sprintf("Connection %d: failed to connect: %v", connID, err))
			time.Sleep(5000 * time.Millisecond)
			continue
		}

		atomic.AddInt64(&m.activeConnections, 1)
		m.logger.LogInfo(fmt.Sprintf("Connection %d: successfully connected", connID))

		for key, value := range m.config.Headers {
			conn.WriteJSON(map[string]string{key: value})
		}

		lifetimeTimer := time.NewTimer(time.Duration(m.config.ConnLifetimeMs) * time.Millisecond)
		done := make(chan struct{})

		go func() {
			defer close(done)
			if body != "" {
				if err := conn.WriteMessage(websocket.TextMessage, []byte(body)); err != nil {
					m.logger.LogError(fmt.Sprintf("Connection %d: failed to send message: %v", connID, err))
					return
				}

				if _, _, err := conn.ReadMessage(); err != nil {
					m.logger.LogError(fmt.Sprintf("Connection %d: failed to read message: %v", connID, err))
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

func (m *WebSocketManager) GetActiveConnections() int64 {
	return atomic.LoadInt64(&m.activeConnections)
}

func (m *WebSocketManager) GetCompletedConnections() int64 {
	return atomic.LoadInt64(&m.completedConnections)
}
