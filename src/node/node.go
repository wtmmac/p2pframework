package node

import (
	"node/link"
	"node/server"
	"fmt"
)

type Node struct {
	server *server.Server
	links *link.List
	upstream chan interface{}
	bus chan interface{}
	running bool
}

func NewNode(upstream chan interface{}) *Node {
	node := new(Node)

	node.upstream = upstream
	node.links = link.NewList()
	node.bus = make(chan interface{})

	node.server = server.NewServer(node.bus, 8000)

	return node
}

func (node *Node) Start() {
	node.server.Start()

	go node.mainLoop()
}

func (node *Node) onMessage(msg link.OnMessage) {
	fmt.Printf("Message: %s\n", msg.Message)

	node.links.Each(func(link *link.Link) {
		if link != msg.Source {
			link.SendMessage(msg.Message)
		}
	})
}

func (node *Node) onConnect(msg server.OnConnection) {
	link := link.NewLink(node.bus, msg.Connection)

	node.links.Add(link)

	fmt.Printf("%d nodes connected.\n", node.links.Len())
}

func (node *Node) onDisconnect(msg link.OnDisconnect) {
	node.links.Remove(msg.Link)

	fmt.Printf("%d nodes connected.\n", node.links.Len())
}

func (node *Node) onBusMessage(msg interface{}) {
	if m, ok := msg.(server.OnConnection); ok {
		node.onConnect(m)
	}

	if m, ok := msg.(link.OnDisconnect); ok {
		node.onDisconnect(m)
	}

	if m, ok := msg.(link.OnMessage); ok {
		node.onMessage(m)
	}
}

func (node *Node) mainLoop() {
	node.running = true

	for node.running {
		msg := <-node.bus

		node.onBusMessage(msg)
	}

	node.upstream <- nil
}
