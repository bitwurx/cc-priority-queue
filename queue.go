package main

// PriorityQueue is a min binary heap implementation of a priority queue data
// structure.
type PriorityQueue struct {
	// Key is the task resource key
	// count is the number of task nodes in the heap
	// heap is the binary heap where task nodes are stored
	Key   string
	count int
	heap  []*Task
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
	pq.minHeapify(0)
	pq.count--

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

	pq.count++
}

// minHeapify the MaxHeapify function on the priority queue nodes.
func (pq *PriorityQueue) minHeapify(i int) {
	MinHeapify(pq.heap, i)
}

// MinHeapify places the target parent node in the proper position in
// the binary heap.
//
// MinHeapify assumes all subtrees are proper binary heaps.
func MinHeapify(nodes []*Task, i int) {
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
		MinHeapify(nodes, min)
	}
}
