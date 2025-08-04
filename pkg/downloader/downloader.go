// download the webpage and parse its contents before returning it

package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// perform a get request to the url and return its contents
func Get(url string) ([]byte, error) {
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

