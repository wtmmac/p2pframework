package link

import (
	"net"
)

type Link struct {
	connection net.Conn
	upstream chan interface{}
	bus chan interface{}
	running bool
}

// Internal commands

type cmdDisconnect struct {}

// Events

type OnDisconnect struct {
	Link *Link
}

type OnMessage struct {
	Source *Link
	Message string
}

func NewLink(upstream chan interface{}, c net.Conn) *Link {
	link := new(Link)

	link.connection = c
	link.upstream = upstream
	link.bus = make(chan interface{})

	go link.mainLoop()

	return link
}

func (link *Link) SendMessage(msg string) {
	link.connection.Write([]byte(msg + "\n"))
}

func (link *Link) Disconnect() {
	link.bus <- cmdDisconnect{}
}

func (link *Link) readLoop() {
	buffer := make([]byte, 1024)

	for {
		bytes, err := link.connection.Read(buffer)

		if err == nil {
			msg := string(buffer[:bytes-1])
			link.upstream <- OnMessage{Source: link, Message: msg}
		} else {
			link.upstream <- OnDisconnect{link}
			link.connection.Close()
			return
		}
	}
}

func (link *Link) disconnect() {
	link.running = false
}

func (link *Link) onBusMessage(msg interface{}) {
	if _, ok := msg.(cmdDisconnect); ok { link.disconnect() }
}

func (link *Link) mainLoop() {
	link.running = true

	go link.readLoop()

	for link.running {
		msg := <-link.bus

		link.onBusMessage(msg)
	}

	link.connection.Close()
}
