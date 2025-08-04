package main

import (
	"crawler/pkg/downloader"
	"crawler/pkg/queue"
	"fmt"
)

const CRAWL_LIMIT = 100

func main() {
	q := queue.New()

	// crawl the first site
	downloader.GetAndParse("https://www.cc.gatech.edu/", q)

	for !q.IsEmpty() && q.TotalQueued < 100 {
		url := q.Dequeue()
		
		if len(url) == 0 {
			fmt.Println("Exiting because of empty content")
			continue
		}

		println(url)
		downloader.GetAndParse(url, q)
	}

	fmt.Println("Finished running")
}