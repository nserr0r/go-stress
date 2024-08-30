package config

import (
	"encoding/json"
	"flag"
	"fmt"
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
	flag.IntVar(&config.Connections, "conn", 10, "Number of concurrent connections (Количество одновременных соединений)")
	flag.IntVar(&config.ConnDelayMs, "conn-delay", 100, "Delay between connections in milliseconds (Задержка между соединениями в миллисекундах)")
	flag.IntVar(&config.ConnLifetimeMs, "conn-lifetime", 60000, "Lifetime of each connection in milliseconds (Время жизни каждого соединения в миллисекундах)")
	flag.BoolVar(&config.Log, "log", false, "Enable logging to console (Включить логирование в консоль)")
	flag.BoolVar(&config.InsecureSkipVerify, "insecure", false, "Skip SSL certificate verification (Пропустить проверку SSL сертификата)")
	flag.BoolVar(&config.UseSSL, "ssl", false, "Use SSL for secure connections (Использовать SSL для безопасных соединений)")
	flag.StringVar(&config.Body, "body", "", "Custom body to send with HTTP POST or WebSocket message (Пользовательское тело запроса для отправки в HTTP POST или WebSocket)")
	flag.StringVar(&headers, "header", "", "Custom headers in JSON format (Пользовательские заголовки в формате JSON, возможно потребуется экранировать кавычки и другие специальные символы)")
	flag.StringVar(&config.ProxyFile, "proxy-file", "", "Path to file containing list of proxy servers (Путь к файлу со списком прокси-серверов)")
	flag.BoolVar(&config.UseWebSocket, "ws", false, "Use WebSocket instead of HTTP (Использовать WebSocket вместо HTTP)")

	flag.Usage = func() {
		fmt.Println("Использование: go-stress [опции]")
		flag.PrintDefaults()
		fmt.Println("Контакты: Nserr0r (nserr0r@gmail.com)")
		fmt.Println("Описание: Инструмент для стресс-тестирования веб-приложений. Не используйте во вредоносных целях.")
	}

	flag.Parse()

	// Парсинг строки заголовков в карту
	if headers != "" {
		if err := json.Unmarshal([]byte(headers), &config.Headers); err != nil {
			fmt.Printf("Не удалось распарсить заголовки: %v\n", err)
		}
	}

	return config
}

