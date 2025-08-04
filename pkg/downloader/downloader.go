// download the webpage and parse its contents before returning it

package downloader

import (
	"bytes"
	"crawler/pkg/queue"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

const MAX_ERR_COUNT = 10

var (
	errCount = 0 // global count for the number of errors tolerable before exiting
	asciiAndWhiteSpace, _ = regexp.Compile("[a-zA-Z0-9]+")
	skipTags map[string]bool
)

func init() {
	skipTags = map[string]bool {
		"script": true,
		"style":  true,
		"noscript": true,
	}
}

func GetAndParse(url string, queue *queue.Queue) (error) {
	pageContent, err := get(url) 

	if err != nil {
		return err
	}

	parseHTMLAndExtractLinks(pageContent, queue)
	return nil
}

// perform a get request to the url and return its contents
func get(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting response at URL: %s, %v\n", url, err)
		return nil, err
	}

	defer resp.Body.Close()
	resBody, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing request body after request to URL: %s, %v\n", url, err)
		return nil, err
	}

	return resBody, nil
}

// get content of the page, then grab all related links on the page
func parseHTMLAndExtractLinks(content []byte, queue *queue.Queue) {
	htmlTokenizer := html.NewTokenizer(bytes.NewReader(content))
	skipDepth := 0
	tokenCount := 0

	for {
		tokenType := htmlTokenizer.Next()

		if tokenType == html.ErrorToken {
			err := htmlTokenizer.Err()
			if err == io.EOF {
				break // handle the end of file
			}

			fmt.Fprintf(os.Stderr, "Error parsing token: %v\n", err)
			errCount++
			
			if errCount < MAX_ERR_COUNT {
				continue
			} else {
				fmt.Fprintf(os.Stderr, "Too many errors encountered\n")
				os.Exit(1)
			}
		}
		
		// token data gives the text content of the token
		// this would be useful for creating the store later
		// token type gives the tag name
		token := htmlTokenizer.Token()
		tagName := strings.ToLower(token.Data)
		
		// Handle skipped tags (like script, style)
		switch tokenType {
		case html.StartTagToken:
			if skipTags[tagName] {
				skipDepth++
			}
		case html.EndTagToken:
			if skipTags[tagName] && skipDepth > 0 {
				skipDepth--
			}
		}
		
		// Skip processing if we're inside a skipped tag
		if skipDepth > 0 {
			continue
		}
		
		// only count parsed tokens
		tokenCount++
		if tagName == "a" {
			// push the href into the queue
			for _, a := range token.Attr {
				if a.Key == "href" {
					queue.Enqueue(a.Val)
				}
			}
		} else if token.Type.String() == "Text" && asciiAndWhiteSpace.MatchString(token.Data) {
			textContent := strings.TrimSpace(token.Data)
			fmt.Println(textContent)
		}
	}
}