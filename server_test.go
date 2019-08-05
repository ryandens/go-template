package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHome(t *testing.T) {
	// GIVEN a request to the to the root context
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	responseRecorder := httptest.NewRecorder()
	homeHandler := http.HandlerFunc(HomeHandler)

	// WHEN we serve the request to the HomeHandler
	homeHandler.ServeHTTP(responseRecorder, req)

	// VERIFY we have the correct status code
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("HomeHandler returned wrong status code. Expected %v, actual %v", http.StatusOK, status)
	}

	// VERIFY the response body is correct
	expectedBody := "Home"
	if body := responseRecorder.Body.String(); body != expectedBody {
		t.Errorf("HomeHanlder returned incorrect body. Expected %v, actual %v", expectedBody, body)
	}
}

func TestValidateUserName(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	if !ValidateUserName("safe", responseRecorder) {
		t.Fatal("Should have passed validation")
	}

	if ValidateUserName("alert%281%29", responseRecorder) {
		t.Fatal("Should have failed validation")
	}
}
