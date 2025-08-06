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
	downloader.GetAndParse("https://www.ucla.edu/", q)

	for q.TotalQueued() < CRAWL_LIMIT {
		url := q.Dequeue()
		
		if len(url) == 0 {
			continue
		}

		println(url)
		go downloader.GetAndParse(url, q)
	}

	fmt.Println("Finished running")
}