package commands

import (
	"log"
	"net/http"
	"time"
)

// WaitForApplication waits until specified URL returns status code 200
// or timeout exceeds. If timeout equals 0, function will wait indefinitely.
func WaitForApplication(url string, timeout time.Duration) bool {
	httpClient := http.Client{Timeout: timeout}

	resultChannel := make(chan bool)

	var timeoutChannel <- chan time.Time

	if timeout == 0 {
		timeoutChannel = nil
	} else {
		timeoutChannel = time.After(timeout)
	}

	go sendRequest(httpClient, url, resultChannel)

	for {
		select {
		case r := <-resultChannel:
			if r {
				return true
			} else {
				log.Println("Sleeping for 2 seconds...")
				time.Sleep(2 * time.Second)
				go sendRequest(httpClient, url, resultChannel)
			}
		case <-timeoutChannel:
			log.Printf("WaitForApplication: request to url %s has timed out", url)
			return false
		}
	}

}

func sendRequest(httpClient http.Client, url string, resultChannel chan bool) {
	response, err := httpClient.Get(url)

	if err != nil {
		log.Printf("sendRequest: URL: %s, received HTTP error: %s", url, err)
		resultChannel <- false
		return
	}

	result := response.StatusCode == http.StatusOK

	if ! result {
		log.Printf("sendRequest: received status code: %d\n", response.StatusCode)
	}

	response.Body.Close()

	resultChannel <- result

}