package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"golang.org/x/net/proxy"
)

type ProxyStatus int

const (
	ProxyAlive ProxyStatus = iota
	ProxyDead
	ProxyUnused
)

type ProxyInfo struct {
	Address string
	Status  ProxyStatus
}

type ProxyManager struct {
	config     Config
	logger     Logger
	proxies    []ProxyInfo
	proxyMutex sync.Mutex
}

func NewProxyManager(config Config, logger Logger) *ProxyManager {
	manager := &ProxyManager{
		config: config,
		logger: logger,
	}

	if config.ProxyFile != "" {
		manager.loadProxies(config.ProxyFile)
	}

	return manager
}

func (p *ProxyManager) loadProxies(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		p.logger.LogError(fmt.Sprintf("Failed to open proxy file: %v", err))
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxyAddr := strings.TrimSpace(scanner.Text())
		if proxyAddr != "" {
			p.proxies = append(p.proxies, ProxyInfo{Address: proxyAddr, Status: ProxyUnused})
		}
	}

	if err := scanner.Err(); err != nil {
		p.logger.LogError(fmt.Sprintf("Error reading proxy file: %v", err))
	}
}

func (p *ProxyManager) getNextProxy() *ProxyInfo {
	p.proxyMutex.Lock()
	defer p.proxyMutex.Unlock()

	for i := range p.proxies {
		if p.proxies[i].Status == ProxyUnused || p.proxies[i].Status == ProxyAlive {
			return &p.proxies[i]
		}
	}
	return nil
}

func (p *ProxyManager) markProxyAsDead(proxy *ProxyInfo) {
	p.proxyMutex.Lock()
	defer p.proxyMutex.Unlock()
	proxy.Status = ProxyDead
}

func (p *ProxyManager) createDialerWithProxy(proxyInfo *ProxyInfo) *websocket.Dialer {
	proxyURL, err := url.Parse(proxyInfo.Address)
	if err != nil {
		p.logger.LogError(fmt.Sprintf("Invalid proxy URL %s: %v", proxyInfo.Address, err))
		return nil
	}

	switch proxyURL.Scheme {
	case "http", "https":
		return &websocket.Dialer{
			Proxy: http.ProxyURL(proxyURL),
		}
	case "socks5":
		auth := &proxy.Auth{}
		if proxyURL.User != nil {
			auth.User = proxyURL.User.Username()
			if password, ok := proxyURL.User.Password(); ok {
				auth.Password = password
			}
		}
		dialSocksProxy, err := proxy.SOCKS5("tcp", proxyURL.Host, auth, proxy.Direct)
		if err != nil {
			p.logger.LogError(fmt.Sprintf("Failed to connect to SOCKS5 proxy %s: %v", proxyInfo.Address, err))
			return nil
		}
		proxyInfo.Status = ProxyAlive
		return &websocket.Dialer{
			NetDial: dialSocksProxy.Dial,
		}
	default:
		p.logger.LogError(fmt.Sprintf("Unsupported proxy type %s", proxyURL.Scheme))
		return nil
	}
}

func (p *ProxyManager) getProxyStatusString() string {
	p.proxyMutex.Lock()
	defer p.proxyMutex.Unlock()

	var aliveCount, deadCount, unusedCount int

	for _, proxy := range p.proxies {
		switch proxy.Status {
		case ProxyAlive:
			aliveCount++
		case ProxyDead:
			deadCount++
		case ProxyUnused:
			unusedCount++
		}
	}

	return fmt.Sprintf(" | Alive Proxies: %s%03d%s | Dead Proxies: %s%03d%s | Unused Proxies: %s%03d%s",
		colorLightGreen, aliveCount, colorReset,
		colorRed, deadCount, colorReset,
		colorYellow, unusedCount, colorReset)
}
