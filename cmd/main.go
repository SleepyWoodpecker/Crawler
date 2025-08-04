package main

import (
	"crawler/pkg/downloader"
	"fmt"
)

func main() {
	content, _ := downloader.Get("https://gopl.io")

	fmt.Println(string(content))
}