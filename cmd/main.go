package main

import (
	"crawler/pkg/downloader"
)

func main() {
	downloader.GetAndParse("https://gopl.io")
}