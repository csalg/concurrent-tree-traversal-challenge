package main

import (
	"fmt"
	"sync"
)

func main() {
	apiClient := NewApiClient()
	walkQueue := NewWalkQueue()
	rootPage := FetchPageJob{Id: 1, ParentIndexSize: 0}

	var wg sync.WaitGroup
	wg.Add(1)
	go worker(apiClient, walkQueue, &wg, &rootPage, RETRIES_PER_PAGE)
	wg.Wait()

	fmt.Println("")
	fmt.Printf("The total indexing size of all pages is %d bytes\n", walkQueue.GetTotalIndexingSize())
	fmt.Printf("Fetched indexing size for %d pages \n", walkQueue.GetPagesFetched())
}

// The way I handle the spawning and retries here is recursively. This is definitely not something I would do in production,
// as the number of idle goroutines will scale exponentially.
// In production I would have a worker pool with n coroutines running concurrently, each with their own channel, and a
// queue with the ids of the idle channels.
func worker(apiClient *APIClient, walkQueue *WalkQueue, wg *sync.WaitGroup, pageToFetch *FetchPageJob, retries int) {
	defer wg.Done()
	Debug(fmt.Sprintf("Fetching id %d with parent indexing size %d", pageToFetch.Id, pageToFetch.ParentIndexSize))
	Debug(fmt.Sprintf("Retries %d", retries))

	if retries == 0 {
		panic("Number of retries exceeded")
	}

	// Fetch page
	page, err := apiClient.GetPage(pageToFetch.Id)
	if err != nil {
		// Retry
		wg.Add(1)
		retries--
		go worker(apiClient, walkQueue, wg, pageToFetch, retries)
		return
	}
	// Register index size
	// Remark: len(str String) outputs the number of bytes in the string
	indexSize := pageToFetch.ParentIndexSize + len(page.Content)
	walkQueue.RegisterPageIndexSize(indexSize)

	// Enqueue children
	for _, id := range page.Children {
		Debug(fmt.Sprintf("Enqueuing %d", id))
		walkQueue.Enqueue(&FetchPageJob{ParentIndexSize: indexSize, Id: id})
	}
	// Spawn new workers
	for {
		pageToFetch = walkQueue.Dequeue()
		if pageToFetch == nil {
			return
		}
		wg.Add(1)
		go worker(apiClient, walkQueue, wg, pageToFetch, RETRIES_PER_PAGE)
	}
}

func Debug(str string) {
	if DEBUG {
		fmt.Println(str)
	}
}
