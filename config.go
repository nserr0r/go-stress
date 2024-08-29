package main

import (
	"flag"
)

type Config struct {
	Host               string
	Path               string
	Connections        int
	Messages           int
	DelayMs            int
	ConnDelayMs        int
	ConnLifetimeMs     int
	Log                bool
	CustomMsg          string
	GenMsg             bool
	ProxyFile          string
	UseSSL             bool
	InsecureSkipVerify bool
}

func LoadConfig() Config {
	var config Config
	flag.StringVar(&config.Host, "host", "localhost:443", "WebSocket server host")
	flag.StringVar(&config.Path, "path", "/ws", "WebSocket server path")
	flag.IntVar(&config.Connections, "conn", 10, "Number of concurrent WebSocket connections")
	flag.IntVar(&config.Messages, "msg", 0, "Number of messages per client (set to 0 for no messages)")
	flag.IntVar(&config.DelayMs, "delay", 100, "Delay between messages in milliseconds")
	flag.IntVar(&config.ConnDelayMs, "conn-delay", 100, "Delay between connections in milliseconds")
	flag.IntVar(&config.ConnLifetimeMs, "conn-lifetime", 1000, "Lifetime of each connection before reconnecting in milliseconds")
	flag.BoolVar(&config.Log, "log", false, "Enable logging to console")
	flag.StringVar(&config.CustomMsg, "msg-text", "", "Custom message to send (ignored if empty)")
	flag.BoolVar(&config.GenMsg, "gen-msg", false, "Generate random message if true")
	flag.StringVar(&config.ProxyFile, "proxy-file", "", "Path to file containing list of proxy servers")
	flag.BoolVar(&config.UseSSL, "ssl", false, "Use SSL for WebSocket connections (wss)")
	flag.BoolVar(&config.InsecureSkipVerify, "insecure", true, "Skip SSL certificate verification (not recommended)")

	flag.Parse()
	return config
}
