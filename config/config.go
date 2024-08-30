package config

import (
	"encoding/json"
	"flag"
	"log"
)

type Config struct {
	Host               string
	Path               string
	Connections        int
	ConnDelayMs        int
	ConnLifetimeMs     int
	Log                bool
	InsecureSkipVerify bool
	UseSSL             bool
	Body               string
	Headers            map[string]string
	ProxyFile          string
	UseWebSocket       bool
}

func LoadConfig() *Config {
	config := &Config{}
	var headers string
	flag.StringVar(&config.Host, "host", "localhost:3001", "Server host")
	flag.StringVar(&config.Path, "path", "/crypt/ws", "Server path")
	flag.IntVar(&config.Connections, "conn", 10, "Number of concurrent connections")
	flag.IntVar(&config.ConnDelayMs, "conn-delay", 100, "Delay between connections in milliseconds")
	flag.IntVar(&config.ConnLifetimeMs, "conn-lifetime", 60000, "Lifetime of each connection before reconnecting in milliseconds")
	flag.BoolVar(&config.Log, "log", false, "Enable logging to console")
	flag.BoolVar(&config.InsecureSkipVerify, "insecure", false, "Skip SSL certificate verification")
	flag.BoolVar(&config.UseSSL, "ssl", false, "Use SSL (wss for WebSocket or https for HTTP)")
	flag.StringVar(&config.Body, "body", "", "Custom body to send with HTTP POST or WebSocket message")
	flag.StringVar(&headers, "header", "", "Custom headers in JSON format")
	flag.StringVar(&config.ProxyFile, "proxy-file", "", "Path to file containing list of proxy servers")
	flag.BoolVar(&config.UseWebSocket, "ws", false, "Use WebSocket instead of HTTP")

	flag.Parse()

	if headers != "" {
		if err := json.Unmarshal([]byte(headers), &config.Headers); err != nil {
			log.Fatalf("Failed to parse headers: %v", err)
		}
	}

	return config
}
