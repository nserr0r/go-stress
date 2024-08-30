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

	flag.StringVar(&config.Host, "host", "localhost:3001", "Server host (Хост сервера, например, localhost:3001)")
	flag.StringVar(&config.Path, "path", "/crypt/ws", "Server path (Путь на сервере, например, /api/test)")
	flag.IntVar(&config.Connections, "conn", 100, "Number of concurrent connections (Количество одновременных соединений)")
	flag.IntVar(&config.ConnDelayMs, "conn-delay", 10, "Delay between connections in milliseconds (Задержка между соединениями в миллисекундах)")
	flag.IntVar(&config.ConnLifetimeMs, "conn-lifetime", 1000, "Lifetime of each connection in milliseconds (Время жизни каждого соединения в миллисекундах)")
	flag.BoolVar(&config.Log, "log", false, "Enable logging to console (Включить логирование в консоль)")
	flag.BoolVar(&config.InsecureSkipVerify, "insecure", true, "Skip SSL certificate verification (Пропустить проверку SSL сертификата)")
	flag.BoolVar(&config.UseSSL, "ssl", false, "Use SSL for secure connections (Использовать SSL для безопасных соединений)")
	flag.StringVar(&config.Body, "body", "", "Custom body to send with HTTP POST or WebSocket message (Пользовательское тело запроса для отправки в HTTP POST или WebSocket)")
	flag.StringVar(&headers, "header", "", "Custom headers in JSON format (Пользовательские заголовки в формате JSON)")
	flag.StringVar(&config.ProxyFile, "proxy-file", "", "Path to file containing list of proxy servers (Путь к файлу со списком прокси-серверов)")
	flag.BoolVar(&config.UseWebSocket, "ws", false, "Use WebSocket instead of HTTP (Использовать WebSocket вместо HTTP)")

	flag.Usage = func() {
		log.Println("Использование: go-stress [опции]")
		flag.PrintDefaults()
		log.Println("Контакты: Nserr0r (nserr0r@gmail.com)")
		log.Println("Описание: Инструмент для стресс-тестирования веб-приложений. Не используйте во вредоносных целях.")
	}

	flag.Parse()

	// Парсинг строки заголовков в карту
	if headers != "" {
		if err := json.Unmarshal([]byte(headers), &config.Headers); err != nil {
			log.Fatalf("Не удалось распарсить заголовки: %v", err)
		}
	}

	return config
}
