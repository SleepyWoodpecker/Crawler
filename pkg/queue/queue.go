package queue

import (
	"sync"
)

type Queue struct {
	Urls []string
	mu sync.Mutex
	length int
	totalQueued int
}

func New() *Queue {
	return &Queue{
		Urls: make([]string, 0, 20), // set an arbitrary initial capacity of 20 on the queue
		length: 0,
		totalQueued: 0,
	}
}

func (q *Queue) Enqueue(url string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.Urls = append(q.Urls, url)
	q.length++
	q.totalQueued++
}

// not too sure how to tell between an empty queue and a finished program
func (q *Queue) Dequeue() string {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.isEmpty() {
		return ""
	}

	url := q.Urls[0]
	q.Urls = q.Urls[1:]
	q.length--

	return url
}

func (q *Queue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.isEmpty()
}

func (q *Queue) isEmpty() bool {
	return len(q.Urls) == 0
}

func (q *Queue) TotalQueued() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.totalQueued
}