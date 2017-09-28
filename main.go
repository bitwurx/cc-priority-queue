package main

import (
	"github.com/bitwurx/jrpc2"
)

func main() {
	InitDatabase()
	s := jrpc2.NewServer(":8888", "/rpc")
	NewApiV1(&PriorityQueueModel{}, s)
	s.Start()
}
