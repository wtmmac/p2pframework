package main

import (
	"node"
)

func main() {
	bus := make(chan interface{})

	node := node.NewNode(bus)

	node.Start()

	<-bus
}
