package proxy

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
)

type ProxyInfo struct {
	URL       *url.URL
	IsWorking bool
}

type ProxyManager struct {
	proxies      []*ProxyInfo
	currentIndex int
	mu           sync.Mutex
}

func NewProxyManager(proxyFile string) *ProxyManager {
	manager := &ProxyManager{}

	file, err := os.Open(proxyFile)
	if err != nil {
		log.Fatalf("Failed to open proxy file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxyURL, err := url.Parse(scanner.Text())
		if err != nil {
			log.Printf("Invalid proxy URL: %s", scanner.Text())
			continue
		}
		manager.proxies = append(manager.proxies, &ProxyInfo{
			URL:       proxyURL,
			IsWorking: true,
		})
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading proxy file: %v", err)
	}

	return manager
}

func (m *ProxyManager) GetNextProxy() *ProxyInfo {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.proxies) == 0 {
		return nil
	}

	for {
		proxy := m.proxies[m.currentIndex]
		m.currentIndex = (m.currentIndex + 1) % len(m.proxies)
		if proxy.IsWorking {
			return proxy
		}
	}
}

// Implementing the ProxyReporter interface
func (m *ProxyManager) GetProxyStatusString() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	working := 0
	dead := 0
	for _, proxy := range m.proxies {
		if proxy.IsWorking {
			working++
		} else {
			dead++
		}
	}

	return fmt.Sprintf(" | Working Proxies: %d %s| Dead Proxies: %d %s", working, "\033[92m", dead, "\033[31m")
}

