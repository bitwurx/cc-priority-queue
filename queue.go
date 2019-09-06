package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

// Task is a unit of work that is queued in the priority queue.
type Task struct {
	// Id is the unique version 1 uuid assigned for task identification.
	// Priority is the queue priority order.
	Id       string  `json:"_key"`
	Priority float64 `json:"priority"`
}

// PriorityQueue is a min binary heap implementation of a priority queue data
// structure.
type PriorityQueue struct {
	// Key is the task resource key.
	// count is the number of task nodes in the heap.
	// heap is the binary heap where task nodes are stored.
	Key   string  `json:"_key"`
	count int     `json:"count"`
	heap  []*Task `json:"heap"`
}

// NewPriorityQueue returns an initialized priority queue instance.
func NewPriorityQueue(key string) *PriorityQueue {
	return &PriorityQueue{key, 0, make([]*Task, 0)}
}

// List returns all priority queue nodes.
func (pq *PriorityQueue) List() []*Task {
	return pq.heap
}

// Peek returns the min heap node without modifying the heap.
func (pq *PriorityQueue) Peek() *Task {
	return pq.heap[0]
}

// Pop removes and returns the min heap node.
func (pq *PriorityQueue) Pop() *Task {
	if pq.count == 0 {
		return nil
	}
	min := pq.heap[0]
	pq.heap[0] = pq.heap[pq.count-1]
	pq.heap = pq.heap[:pq.count-1]
	pq.minHeapify(pq.heap, 0)
	pq.count--
	log.Printf("popped task [%v] from queue [%s]", min, pq.Key)

	return min
}

// Push inserts a task into the task nodes in priority order.
func (pq *PriorityQueue) Push(t *Task) {
	pq.heap = append(pq.heap, t)
	i := len(pq.heap) - 1

	for i > 0 {
		parent := (i - 1) / 2

		if pq.heap[i].Priority < pq.heap[parent].Priority {
			pq.heap[i] = pq.heap[parent]
			pq.heap[parent] = t
			i = parent
		} else {
			break
		}
	}

	log.Printf("pushed task [%v] to queue [%s]", t, pq.Key)
	pq.count++
}

// Remove the node from the priority queue with the provided id.
func (pq *PriorityQueue) Remove(id string) error {
	if pq.count == 0 {
		return nil
	}
	nodeIndex := -1
	for i, node := range pq.heap {
		if node.Id == id {
			nodeIndex = i
		}
	}
	if nodeIndex == -1 {
		return errors.New("id not found")
	}
	pq.heap[nodeIndex] = pq.heap[pq.count-1]
	pq.heap = pq.heap[:pq.count-1]

	i := nodeIndex
	for i > 0 {
		parent := (i - 1) / 2
		node := pq.heap[i]

		if pq.heap[i].Priority < pq.heap[parent].Priority {
			pq.heap[i] = pq.heap[parent]
			pq.heap[parent] = node
			i = parent
		} else {
			break
		}
	}

	pq.minHeapify(pq.heap, nodeIndex)
	pq.count--

	return nil
}

// Save writes the priority queue to the database.
func (pq *PriorityQueue) Save(pqModel Model) (DocumentMeta, error) {
	nodes := make([]*Task, 0)
	for _, node := range pq.heap {
		nodes = append(nodes, &Task{Id: node.Id, Priority: node.Priority})
	}
	pq.heap = nodes
	return pqModel.Save(pq)
}

// minHeapify the MaxHeapify function on the priority queue nodes.
func (pq *PriorityQueue) minHeapify(nodes []*Task, i int) {
	left := (i * 2) + 1
	right := (i * 2) + 2
	min := i

	if left < len(nodes) && nodes[left].Priority < nodes[i].Priority {
		min = left
	}
	if right < len(nodes) && nodes[right].Priority < nodes[i].Priority {
		min = right
	}
	if min != i {
		node := nodes[i]
		nodes[i] = nodes[min]
		nodes[min] = node
		pq.minHeapify(nodes, min)
	}

	pq.heap = nodes
}

// MarshalJSON serializes the priority queue key, count, and nodes
// members.
func (pq *PriorityQueue) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(
		fmt.Sprintf(`{"_key": "%s", "count": %d, "heap": %s}`, pq.Key, pq.count, (func() string {
			nodes := bytes.NewBuffer([]byte("["))
			for i, node := range pq.heap {
				nodes.WriteString(
					fmt.Sprintf(`{"_key": "%s", "priority": %.1f}`, node.Id, node.Priority),
				)
				if i < (pq.count - 1) {
					nodes.WriteByte(',')
				}
			}
			nodes.WriteString("]")
			return nodes.String()
		})(),
		))
	return buf.Bytes(), nil
}

// UnmarshalJSON deserializes the stored priority queue meta data into
// a priority queue instance.
func (pq *PriorityQueue) UnmarshalJSON(b []byte) error {
	data := make(map[string]interface{})
	json.Unmarshal(b, &data)
	pq.Key = data["_key"].(string)
	pq.count = int(data["count"].(float64))
	for _, node := range data["heap"].([]interface{}) {
		v, _ := node.(map[string]interface{})
		task := &Task{
			Id:       v["_key"].(string),
			Priority: v["priority"].(float64),
		}
		pq.heap = append(pq.heap, task)
	}
	return nil
}
