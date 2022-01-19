package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

func main() {
	// This implementation works asynchronously but still does not handle network shenanigans or auth.
	fetchJob, _ := newPageFetchJob(1, 0)
	c := make(chan SharedState)
	fmt.Println("wtf")

	var wg sync.WaitGroup
	wg.Add(1)
	go Handler(c, &wg)
	c <- SharedState{
		Queue:  []FetchPageJob{*fetchJob},
		Result: 0}
	wg.Wait()

	fmt.Println("wtf")
}

// SharedState is the state that is in the critical zone and is passed around the workers via a channel
type SharedState struct {
	Queue  []FetchPageJob
	Result int
}

func Handler(c <-chan SharedState, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Inside handler")
	sharedState := <-c

	job := sharedState.Queue[0]
	queue = sharedState.Queue[1:]

	res, err := GetPageById(job.Id)
	if err != nil {
		panic(err)
	}
	fmt.Println(PrettyPrint(*res))

	indexSize := len(res.Content) + job.ParentIndexSize
	fmt.Println(indexSize)
	result = sharedState.Result + indexSize

	for _, childId := range res.Children {
		fetchJob, err := newPageFetchJob(childId, indexSize)
		if err != nil {
			panic(err)
		}
		queue = append(queue, *fetchJob)
	}

}

func newPageFetchJob(id, parentIndexSize int) (*FetchPageJob, error) {
	if id < 0 {
		return nil, errors.New(fmt.Sprintf("Expected non-negative value for id, but got %d", id))
	}
	if parentIndexSize < 0 {
		return nil, errors.New(fmt.Sprintf("Expected non-negative value for parentIndexSize, but got %d", parentIndexSize))
	}
	return &FetchPageJob{Id: id, ParentIndexSize: parentIndexSize}, nil
}

func GetPageById(id int) (*Response, error) {

	res, err := http.Get(fmt.Sprintf("http://localhost:%d/api/page?id=%d", PORT, id))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Cannot unmarshal JSON")
	}

	return &result, nil
}

// PrettyPrint prints struct using tabs for indentation
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
