package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// If we run this file on real AWS EC2 instance
// or any other cloud provider that exposes metadata
// on http://169.254.169.254/latest/meta-data/
// it should give us back instance id

var metadataEndpointURL = "http://169.254.169.254"
var instanceIdURI = "/latest/meta-data/instance-id"

// Default http.Client has no timeout
// so let's set one
var httpClient = &http.Client{
	Timeout: time.Second * 5, // 5 seconds timeout is long enough
}

func main() {
	// Get our instance ID
	id, err := getInstanceId()
	if err != nil {
		// if some error occur we just panic
		panic(err)
	}

	// Output instance ID and exit
	fmt.Println("Instance ID:", id)
}

// getInstanceId retrieves metadata from standard
// endpoint
func getInstanceId() (string, error) {
	response, err := httpClient.Get(metadataEndpointURL + instanceIdURI)
	if err != nil {
		return "", err
	}

	// Check status code
	// since this is GET then expected is 200
	if response.StatusCode != http.StatusOK {
		return "", ErrorNotFound
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// We've got the body as bytes, so let's return it now as string with nil error
	return string(bodyBytes), nil
}

type ApiError string

func (e ApiError) Error() string {
	return string(e)
}

const (
	ErrorNotFound ApiError = "Unexpected 404"
)
