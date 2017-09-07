package queue_test

import (
	"testing"

	"concord-pq/queue"
	"concord-pq/task"
)

func TestNewPriorityQueue(t *testing.T) {
	pq := queue.NewPriorityQueue("key")
	if pq.Key != "key" {
		t.Fatal("expected key to be 'key'")
	}
	if pq.Count != 0 {
		t.Fatal("expected node count to be 0")
	}
}

func TestPriorityQueueInsert(t *testing.T) {
	tasks := []*task.Task{
		&task.Task{Priority: 22.5},
		&task.Task{Priority: 3.5},
		&task.Task{Priority: 16.5},
	}
	pq := queue.NewPriorityQueue("key")
	for _, task := range tasks {
		pq.Insert(task)
	}
	if pq.Count != 3 {
		t.Fatal("expected node count to be 3")
	}
	nodes := pq.List()
	if nodes[0].Priority != 3.5 || nodes[1].Priority != 22.5 || nodes[2].Priority != 16.5 {
		t.Fatal("invalid node order")
	}
}
