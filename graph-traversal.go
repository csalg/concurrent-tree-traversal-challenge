package graphTraversal

type PageResponse struct {
	Content  string `json:"content"`
	Children []int  `json:"children"`
}

type FetchPageJob struct {
	Id              int
	ParentIndexSize int
}

// InitGraphTraversalMonitor is meant to be run as a goroutine
// 1. It keeps all of the state that is in the critical zone and must be processed atomically.
// 2. It implements some simple procedures on this state using a select switch.
func InitGraphTraversalMonitor() {
	if isRunning {
		return
	}
	isRunning = true

	queue := []*FetchPageJob{}
	pending := 0
	total := 0

	for {
		select {
		case newTask := <-enqueueFetchTask:
			queue = append(queue, newTask)
		case newPageIndexSize := <-pageIndexSizeReceived:
			total = total + newPageIndexSize
			pending = pending - 1
		case <-dequeueRequest:
			if len(queue) == 0 {
				dequeueResponse <- nil
			} else {
				// In production I would look for a more performant alternative (probably a linked list implementation).
				nextJob := queue[0]
				queue = queue[1:]
				pending = pending + 1
				dequeueResponse <- nextJob
			}
		case <-doneRequest:
			if len(queue) == 0 && pending == 0 {
				doneResponse <- true
			} else {
				doneResponse <- false
			}
		case <-resultRequest:
			resultResponse <- total
		}
	}
}

// EnqueueJob pushes a page to fetch to the queue
func EnqueueJob(job *FetchPageJob) {
	enqueueFetchTask <- job
}

// RegisterPageResponse notifies that a response has been received.
func RegisterPageResponse(parentIndexSize int, response *PageResponse) {
	pageIndexSizeReceived <- parentIndexSize + len(response.Content)
}

// Dequeue returns a pointer to a page to fetch or nil if the queue is empty.
func Dequeue() *FetchPageJob {
	dequeueRequest <- true
	result := <-dequeueResponse
	return result
}

// IsDone returns true if there are no pages to fetch and no pages being fetched.
func IsDone() bool {
	doneRequest <- true
	result := <-doneResponse
	return result
}

// Result returns the total indexing size of all pages.
func Result() int {
	resultRequest <- true
	result := <-resultResponse
	return result
}

// Private

// Channels
var enqueueFetchTask = make(chan *FetchPageJob)
var pageIndexSizeReceived = make(chan int)

var dequeueRequest = make(chan bool)
var dequeueResponse = make(chan *FetchPageJob)

var doneRequest = make(chan bool)
var doneResponse = make(chan bool)

var resultRequest = make(chan bool)
var resultResponse = make(chan int)

var isRunning = false
