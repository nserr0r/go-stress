Go-Stress

Go-Stress — это мощный и гибкий инструмент для стресс-тестирования веб-приложений, созданный для проверки отказоустойчивости моих приложений. Этот инструмент позволяет генерировать параллельные HTTP и WebSocket соединения, поддерживает использование прокси-серверов, а также задавать пользовательские заголовки и содержимое запроса.

Важно: Использование этого инструмента в целях, нарушающих закон, категорически не рекомендуется.
Основные возможности

    Поддержка HTTP и WebSocket: Тестируйте свои HTTP и WebSocket эндпоинты с помощью одного инструмента.
    Пользовательские заголовки и содержимое запроса: Отправляйте пользовательские заголовки и тело запроса как для HTTP, так и для WebSocket соединений.
    Поддержка прокси-серверов: Направляйте ваши запросы через список прокси-серверов, включая поддержку HTTP и SOCKS5 прокси.
    SSL/TLS: Поддержка безопасных соединений через SSL/TLS с опциональным флагом --insecure для пропуска проверки сертификатов.
    Управление соединениями: Контролируйте количество одновременных соединений, задержку между подключениями и время жизни каждого соединения.
    Реальное время: Отслеживайте активные и завершенные соединения в реальном времени с цветовой индикацией в консоли.

Установка

Предварительные требования

    Установленный Go 1.16 или выше.

Сборка проекта

    make build

Установка исполняемого файла

    sudo make install

Это установит исполняемый файл go-stress в /usr/local/bin.
Использование

    go-stress [опции]

Опции командной строки
Опция	Описание	По умолчанию
    
    -host	Хост сервера (например, localhost:3001)	localhost:3001
    
    -path	Путь на сервере (например, /api/test)	/crypt/ws
    
    -conn	Количество одновременных соединений	10
    
    -conn-delay	Задержка между установкой новых соединений (в миллисекундах)	100
    
    -conn-lifetime	Время жизни каждого соединения перед переподключением (в миллисекундах)	1000
    
    -log	Включить логирование в консоль	false
    
    -ssl	Использовать SSL для безопасных соединений	false
    
    -insecure	Пропустить проверку SSL сертификата	true
    
    -body	Пользовательское содержимое запроса для отправки в HTTP POST или WebSocket сообщениях	""
    
    -header	Пользовательские заголовки в формате JSON (например, {'Authorization':'Bearer token'})	""
    
    -proxy-file	Путь к файлу, содержащему список прокси-серверов	""
    
    -ws	Использовать WebSocket вместо HTTP	false
    
    -help	Вывести информацию о командах и опциях (на русском языке)	

Примеры использования

Базовый HTTP тест

    go-stress -host=example.com -path=/api/test -conn=100 -conn-delay=50

WebSocket тест с SSL

    go-stress -host=example.com -path=/ws -conn=50 -ssl=true -ws=true

Использование прокси-серверов

    go-stress -host=example.com -path=/api/test -conn=100 -proxy-file=proxies.txt

Отправка пользовательских заголовков и тела запроса

    go-stress -host=example.com -path=/api/test -conn=50 -header="{"Content-Type":"application/json", "Authorization":"Bearer token"}"

Мониторинг статуса

Статус активных и завершенных соединений будет отображаться в реальном времени в консоли с цветовой индикацией:

  Красный: Активные соединения
  Зеленый: Завершенные соединения
  Светло-зеленый: Рабочие прокси (если включены)
  Желтый: Ожидающие соединения или прокси

Лицензия

Этот проект лицензирован под GNU General Public License v3.0 (GPL-3.0). Вы можете свободно распространять и модифицировать этот проект в соответствии с условиями лицензии GPL-3.0.

Для получения более подробной информации смотрите файл LICENSE или посетите сайт GNU.
Контакт

Создатель: Nserr0r (nserr0r@gmail.com)
