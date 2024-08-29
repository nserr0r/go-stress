package main

import (
	"flag"
)

type Config struct {
	Host           string
	Path           string
	Connections    int
	Messages       int
	DelayMs        int
	ConnDelayMs    int
	ConnLifetimeMs int
	Log            bool
	CustomMsg      string
	GenMsg         bool
	ProxyFile      string
}

func LoadConfig() Config {
	var config Config
	flag.StringVar(&config.Host, "host", "localhost:3001", "WebSocket server host")
	flag.StringVar(&config.Path, "path", "/crypt/websocket", "WebSocket server path")
	flag.IntVar(&config.Connections, "conn", 100, "Number of concurrent WebSocket connections")
	flag.IntVar(&config.Messages, "msg", 0, "Number of messages per client (set to 0 for no messages)")
	flag.IntVar(&config.DelayMs, "delay", 0, "Delay between messages in milliseconds")
	flag.IntVar(&config.ConnDelayMs, "conn-delay", 0, "Delay between connections in milliseconds")
	flag.IntVar(&config.ConnLifetimeMs, "conn-lifetime", 100, "Lifetime of each connection before reconnecting in milliseconds")
	flag.BoolVar(&config.Log, "log", false, "Enable logging to console")
	flag.StringVar(&config.CustomMsg, "msg-text", "", "Custom message to send (ignored if empty)")
	flag.BoolVar(&config.GenMsg, "gen-msg", false, "Generate random message if true")
	flag.StringVar(&config.ProxyFile, "proxy-file", "", "Path to file containing list of proxy servers")

	flag.Parse()
	return config
}
