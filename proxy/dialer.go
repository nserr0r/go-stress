package proxy

import (
	"crypto/tls"
	"net/http"
	"net/url"

	"main/config"
	"main/logger"

	"github.com/gorilla/websocket"
)

type DialerFactory struct {
	config       *config.Config
	logger       logger.Logger
	ProxyManager *ProxyManager
}

func NewDialerFactory(config *config.Config, logger logger.Logger) *DialerFactory {
	return &DialerFactory{
		config: config,
		logger: logger,
	}
}

func (f *DialerFactory) CreateHTTPClient(proxyURL *url.URL) *http.Client {
	transport := &http.Transport{}
	if proxyURL != nil {
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	if f.config.UseSSL {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: f.config.InsecureSkipVerify,
		}
	}

	return &http.Client{
		Transport: transport,
	}
}

func (f *DialerFactory) CreateWebSocketDialer(proxyURL *url.URL) *websocket.Dialer {
	dialer := websocket.DefaultDialer
	if proxyURL != nil {
		dialer.Proxy = http.ProxyURL(proxyURL)
	}

	if f.config.UseSSL {
		dialer.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: f.config.InsecureSkipVerify,
		}
	}

	return dialer
}

