package main

import (
	"github.com/bitwurx/jrpc2"
)

func main() {
	InitDatabase()
	s := jrpc2.NewServer(":8080", "/rpc", nil)
	NewApiV1(&PriorityQueueModel{}, s)
	s.Start()
}
