package queue

import (
	"concord-pq/task"
)

// PriorityQueue is a min binary heap implementation of a priority queue data structure
type PriorityQueue struct {
	Count int
	Key   string
	nodes []*task.Task
}

// Insert inserts a task into the task nodes in priority order
func (pq *PriorityQueue) Insert(t *task.Task) {
	pq.nodes = append(pq.nodes, t)
	i := len(pq.nodes) - 1

	for i > 0 {
		parent := (i - 1) / 2

		if pq.nodes[i].Priority < pq.nodes[parent].Priority {
			pq.nodes[i] = pq.nodes[parent]
			pq.nodes[parent] = t
			i = parent
		} else {
			break
		}
	}

	pq.Count++
}

// List returns all priority queue nodes
func (pq *PriorityQueue) List() []*task.Task {
	return pq.nodes
}

// NewPriorityQueue returns an initialized priority queue instance
func NewPriorityQueue(key string) *PriorityQueue {
	return &PriorityQueue{Count: 0, Key: key}
}
