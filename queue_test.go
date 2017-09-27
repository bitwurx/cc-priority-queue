package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewPriorityQueue(t *testing.T) {
	pq := NewPriorityQueue("key")
	if pq.Key != "key" {
		t.Fatal("expected key to be 'key'")
	}
}

func TestPriorityQueuePeek(t *testing.T) {
	nodes := []*Task{
		&Task{Priority: 13.5},
		&Task{Priority: 11.5},
	}
	pq := NewPriorityQueue("key")
	for _, task := range nodes {
		pq.Push(task)
	}
	if pq.Peek() != nodes[1] {
		t.Fatal("peek returned an unexpected task")
	}
}

func TestPriorityQueuePush(t *testing.T) {
	nodes := []*Task{
		&Task{Priority: 22.5},
		&Task{Priority: 3.5},
		&Task{Priority: 16.5},
	}
	pq := NewPriorityQueue("key")
	for _, task := range nodes {
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
	nodes := []*Task{
		&Task{Priority: 75.5},
		&Task{Priority: 63.5},
		&Task{Priority: 22.5},
		&Task{Priority: 72.5},
		&Task{Priority: 65.5},
		&Task{Priority: 55.5},
		&Task{Priority: 80.5},
	}
	MinHeapify(nodes, 0)
	for i, task := range nodes {
		heap[i] = task.Priority
	}
	if heap != [7]float64{22.5, 63.5, 55.5, 72.5, 65.5, 75.5, 80.5} {
		t.Fatal("invalid node order")
	}
}

func TestPriorityQueuePop(t *testing.T) {
	nodes := []*Task{
		&Task{Priority: 22.5},
		&Task{Priority: 5.5},
		&Task{Priority: 125.5},
	}
	pq := NewPriorityQueue("key")
	for _, task := range nodes {
		pq.Push(task)
	}

	if pq.Pop().Priority != nodes[1].Priority {
		t.Fatal("pop returned an unexpected task")
	}
	if pq.Pop().Priority != nodes[0].Priority {
		t.Fatal("pop returned an unexpected task")
	}
	if pq.Pop().Priority != nodes[2].Priority {
		t.Fatal("pop returned an unexpected task")
	}
	if pq.Pop() != nil {
		t.Fatal("expected pop to return nil")
	}
}

func TestPriorityQueueRemove(t *testing.T) {
	heap := [4]float64{}
	nodes := []*Task{
		&Task{Id: "1", Priority: 21.2},
		&Task{Id: "2", Priority: 4.1},
		&Task{Id: "3", Priority: 16.9},
		&Task{Id: "4", Priority: 3.7},
		&Task{Id: "5", Priority: 1.7},
	}
	pq := NewPriorityQueue("test")
	for _, task := range nodes {
		pq.Push(task)
	}
	if err := pq.Remove("0"); err == nil {
		t.Fatal("expected id not found error")
	}
	if err := pq.Remove("4"); err != nil {
		t.Fatal(err)
	}
	for i, node := range pq.List() {
		heap[i] = node.Priority
	}
	if heap != [4]float64{1.7, 4.1, 16.9, 21.2} {
		t.Fatal("got unexpected node order")
	}
}

func TestPriorityQueueMarshalJSON(t *testing.T) {
	pq := NewPriorityQueue("key-123")
	task := NewTask([]byte(`{"priority": 3.5}`))
	pq.Push(task)
	data, err := json.Marshal(pq)
	if err != nil {
		t.Fatal(err)
	}
	dataString := fmt.Sprintf(`{"_key":"key-123","count":1,"heap":[{"_key":"%s","priority":3.5}]}`, task.Id)
	if string(data) != dataString {
		t.Fatal("got unexpected marshal json data string")
	}
}

func TestPriorityQueueUnmarshalJSON(t *testing.T) {
	pq := new(PriorityQueue)
	json.Unmarshal([]byte(`{"_key":"key-xyz","count":1,"heap":[{"_key":"%s","priority":3.5}]}`), pq)
	if pq.Key != "key-xyz" {
		t.Fatal("expected key to be 'key-xyz'")
	}
	if pq.Peek().Priority != 3.5 {
		t.Fatal("expected heap node priority to be 3.5")
	}
}

func TestPriorityQueueSave(t *testing.T) {
	var model Model
	if testing.Short() {
		model = MockModel{}
	} else {
		model = &PriorityQueueModel{}
	}
	pq := NewPriorityQueue("some-key")
	pq.Push(&Task{Priority: 13.5, Key: "some-key", Status: 42})
	if _, err := pq.Save(model); err != nil {
		t.Fatal(err)
	}
}
