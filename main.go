package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const PAGES_TO_FETCH = 10
const PORT = 8099

func main() {
	// This implementation works synchronously and does not handle network shenanigans or auth.

	fetchJob, _ := newPageFetchJob(1, 0)
	queue := []FetchPageJob{*fetchJob}
	indexingSizeSum := 0

	for {
		job := queue[0]
		// I am actually not aware of a standard queue implementation so I am just cutting the slice. In a production
		// system I would put more thought into whether this is performant, look for a linked list implementation or
		// something like that.
		queue = queue[1:]

		res, err := GetPageById(job.Id)
		if err != nil {
			panic(err)
		}
		fmt.Println(PrettyPrint(*res))

		indexSize := len(res.Content) + job.ParentIndexSize
		fmt.Println(indexSize)
		indexingSizeSum = indexingSizeSum + indexSize

		for _, childId := range res.Children {
			fetchJob, err := newPageFetchJob(childId, indexSize)
			if err != nil {
				panic(err)
			}
			queue = append(queue, *fetchJob)
		}
		if len(queue) == 0 {
			break
		}
	}
	fmt.Println("The total indexing size is:", indexingSizeSum)
}

type FetchPageJob struct {
	Id              int
	ParentIndexSize int
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

type Response struct {
	Content  string `json:"content"`
	Children []int  `json:"children"`
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
