Go-Stress

Go-Stress is a powerful and flexible tool for stress-testing web applications, capable of generating concurrent HTTP and WebSocket connections. It supports the use of proxies, custom request headers, and bodies, making it ideal for a wide range of load testing scenarios.
Features

    HTTP and WebSocket Support: Test your HTTP and WebSocket endpoints with a single tool.
    Custom Headers and Body: Send custom headers and body content for both HTTP and WebSocket connections.
    Proxy Support: Route your requests through a list of proxies, with support for both HTTP and SOCKS5 proxies.
    SSL/TLS: Supports secure connections via SSL/TLS with an optional --insecure flag to skip certificate verification.
    Connection Management: Control the number of concurrent connections, delay between connections, and the lifetime of each connection.
    Real-Time Status: Monitor active and completed connections in real time, with color-coded output.

Installation
Prerequisites

    Go 1.16+ installed on your system.

Building the Project

bash

make build

Installing the Binary

bash

sudo make install

This will install the go-stress binary to /usr/local/bin.
Usage

bash

go-stress [options]

Command-Line Options
Option	Description	Default
-host	Server host (e.g., localhost:3001)	localhost:3001
-path	Server path (e.g., /api/test)	/crypt/ws
-conn	Number of concurrent connections	10
-conn-delay	Delay between establishing new connections (in milliseconds)	100
-conn-lifetime	Lifetime of each connection before reconnecting (in milliseconds)	60000
-log	Enable logging to console	false
-ssl	Use SSL for secure connections	false
-insecure	Skip SSL certificate verification	false
-body	Custom body content to send with HTTP POST or WebSocket messages	""
-header	Custom headers in JSON format (e.g., {'Authorization':'Bearer token', 'Content-Type':'application/json'})	""
-proxy-file	Path to a file containing a list of proxy servers	""
-ws	Use WebSocket instead of HTTP	false
Examples
Basic HTTP Test

bash

go-stress -host=example.com -path=/api/test -conn=100 -conn-delay=50

WebSocket Test with SSL

bash

go-stress -host=example.com -path=/ws -conn=50 -ssl=true -ws=true

Using Proxies

bash

go-stress -host=example.com -path=/api/test -conn=100 -proxy-file=proxies.txt

Sending Custom Headers and Body

bash

go-stress -host=example.com -path=/api/test -conn=50 -header="{\"Content-Type\":\"application/json\", \"Authorization\":\"Bearer token\"}" -body="{'key':'value'}"

Status Monitoring

The status of active and completed connections will be displayed in real-time in the console, with color-coded output:

    Red: Active connections
    Green: Completed connections
    Light Green: Working proxies (if enabled)
    Yellow: Pending connections or proxies

Contributing

Contributions are welcome! Please submit issues and pull requests to the GitHub repository.

License

This project is licensed under the GNU General Public License v3.0 (GPL-3.0). You may freely distribute and modify this project under the terms of the GPL-3.0 license.

For more details, see the LICENSE file or visit the GNU website.

