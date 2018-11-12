package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// In this file we want to test a method that uses external API to retrieve instance ID
// to do this we will need to mock a metadata API and overwrite original endpoint
// Additionally tested method should not panic, but should return correct error

// API stub to respond with this instance ID
var testInstanceId = "i-1234567890"

// Metadata endpoint stub
func MetadataEndpointStub() *httptest.Server {
	// return stub API endpoint
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var response string
		switch r.RequestURI {
		case "/latest/meta-data/instance-id": // instanceIdURI is defined in main package
			response = testInstanceId
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
}

// Setup environment for testing metadata API
// and return teardown function
func setupGetInstanceIdTest(t *testing.T) func(t *testing.T) {
	server := MetadataEndpointStub()

	// Override default metadata endpoint URL
	originalURL := metadataEndpointURL
	metadataEndpointURL = server.URL

	return func(t *testing.T) {
		server.Close()
		metadataEndpointURL = originalURL
	}
}

// Run test of getInstanceId
func TestGetInstanceIdPass(t *testing.T) {
	// Run setup
	teardown := setupGetInstanceIdTest(t)
	defer teardown(t)

	// We should get "i-1234567890" here as instance ID from API stub
	id, err := getInstanceId()

	// Error must be nil
	if err != nil {
		t.Errorf("Test %s failed. Expected error [nil], got [%s]", t.Name(), err.Error())
		return
	}

	if id != testInstanceId {
		t.Errorf("Test %s failed. Expected instance id [%s], got [%s]", t.Name(), testInstanceId, id)
	}
}

// Check error we get if we cannot connect to API
// code should not panic in this case, but return appropriate error
func TestGetInstanceIdConectionFailure(t *testing.T) {
	// Skip setup since we do not need running server

	// Should get non-nil error and id value then is irrelevant
	_, err := getInstanceId()

	// Check error non-nil or not
	if err == nil {
		t.Errorf("Test %s failed. Expected non-nil error, got nil", t.Name())
		return
	}
}

// Check for unexpected status code error
func TestGetInstanceIdStatusCodeFailure(t *testing.T) {
	// Run setup
	teardown := setupGetInstanceIdTest(t)
	defer teardown(t)

	// Change the request URI
	instanceIdURI = "/404"

	// Should get non-nil error and id value then is irrelevant
	_, err := getInstanceId()

	// Check error non-nil or not
	if err == nil {
		t.Errorf("Test %s failed. Expected non-nil error, got nil", t.Name())
		return
	}

	// Ensure proper error returned
	if err != ErrorNotFound {
		t.Errorf("Test %s failed. Expected error [%s], got [%s]", t.Name(), ErrorNotFound, err)
		return
	}
}
