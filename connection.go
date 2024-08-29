package main

type Connection interface {
	Connect(url string) error
	Disconnect() error
	SendMessage(message string) error
	ReceiveMessage() (string, error)
}

type ConnectionManager interface {
	ManageConnection(connID int)
	Status() (active int64, completed int64)
	MaxConnections() int
}
