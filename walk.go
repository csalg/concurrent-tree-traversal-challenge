package main

import "sync"

type FetchPageJob struct {
	Id              int
	ParentIndexSize int
}

type WalkQueue struct {
	mutex sync.Mutex
	items []*FetchPageJob
	total int
}

func (q *WalkQueue) Enqueue(job *FetchPageJob) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.items = append(q.items, job)
}

func (q *WalkQueue) Dequeue() *FetchPageJob {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.items) == 0 {
		return nil
	}
	result := q.items[0]
	q.items = q.items[1:]
	return result
}

func (q *WalkQueue) RegisterPage(size int) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.total += size
}

func (q *WalkQueue) GetTotalIndexingSize() int{
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return q.total
}