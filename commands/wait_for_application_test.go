package commands

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)


func executeQuery(handler http.HandlerFunc, timeout time.Duration) bool {
	server := httptest.NewServer(handler)
	defer server.Close()

	return WaitForApplication(server.URL, timeout)

}

func TestServerStartingNormally(t *testing.T) {
	handlerFunc := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "Hello, Test")
	})

	result := executeQuery(handlerFunc, 0)

	if !result {
		t.FailNow()
	}
}

func TestServerNotResponding(t *testing.T) {
	handlerFunc := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(2 * time.Second)
		fmt.Fprintln(writer, "Hello, Test")
	})

	result := executeQuery(handlerFunc, 100 * time.Millisecond)

	if result {
		t.FailNow()
	}
}

func TestServerWrongStatusCode(t *testing.T) {
	handlerFunc := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNotFound)
	})

	result := executeQuery(handlerFunc, 100 * time.Millisecond)

	if result {
		t.FailNow()
	}
}

func TestServerStartsSendingOKStatusCodeAfterSomeTime(t *testing.T) {
	applicationStarted := time.After(5 * time.Second)

	handlerFunc := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		applicationStillStarting := time.After(500 * time.Millisecond)

		select {
		case <- applicationStillStarting:
			writer.WriteHeader(http.StatusNotFound)
		case <-applicationStarted:
			writer.WriteHeader(http.StatusOK)
		}
	})

	result := executeQuery(handlerFunc, 10 * time.Second)

	if ! result {
		t.FailNow()
	}

}