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
}

func TestPriorityQueuePeek(t *testing.T) {
	tasks := []*task.Task{
		&task.Task{Priority: 13.5},
		&task.Task{Priority: 11.5},
	}
	pq := queue.NewPriorityQueue("key")
	for _, task := range tasks {
		pq.Push(task)
	}
	if pq.Peek() != tasks[1] {
		t.Fatal("peek returned an unexpected task")
	}
}

func TestPriorityQueuePush(t *testing.T) {
	tasks := []*task.Task{
		&task.Task{Priority: 22.5},
		&task.Task{Priority: 3.5},
		&task.Task{Priority: 16.5},
	}
	pq := queue.NewPriorityQueue("key")
	for _, task := range tasks {
		pq.Push(task)
	}
	heap := [3]float64{}
	for i, task := range pq.List() {
		heap[i] = task.Priority
	}
	if heap != [3]float64{3.5, 22.5, 16.5} {
		t.Fatal("invalid node order")
	}
}

func TestMinHeapify(t *testing.T) {
	heap := [7]float64{}
	tasks := []*task.Task{
		&task.Task{Priority: 75.5},
		&task.Task{Priority: 63.5},
		&task.Task{Priority: 22.5},
		&task.Task{Priority: 72.5},
		&task.Task{Priority: 65.5},
		&task.Task{Priority: 55.5},
		&task.Task{Priority: 80.5},
	}
	queue.MinHeapify(tasks, 0)
	for i, task := range tasks {
		heap[i] = task.Priority
	}
	if heap != [7]float64{22.5, 63.5, 55.5, 72.5, 65.5, 75.5, 80.5} {
		t.Fatal("invalid node order")
	}
}

func TestPriorityQueuePop(t *testing.T) {
	tasks := []*task.Task{
		&task.Task{Priority: 22.5},
		&task.Task{Priority: 5.5},
		&task.Task{Priority: 125.5},
	}
	pq := queue.NewPriorityQueue("key")
	for _, task := range tasks {
		pq.Push(task)
	}

	if pq.Pop().Priority != tasks[1].Priority {
		t.Fatal("pop returned an unexpected task")
	}
	if pq.Pop().Priority != tasks[0].Priority {
		t.Fatal("pop returned an unexpected task")
	}
	if pq.Pop().Priority != tasks[2].Priority {
		t.Fatal("pop returned an unexpected task")
	}
	if pq.Pop() != nil {
		t.Fatal("expected pop to return nil")
	}
}
