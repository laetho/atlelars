package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Table-driven test for the handlerAction function
func TestHandlerAction(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		payload        ActionRequest
		expectedCode   int
		expectedBody   ActionResponse
		expectedHeader string
	}{
		{
			name:           "Valid POST request with action",
			method:         http.MethodPost,
			url:            "/action/test-kind",
			payload:        ActionRequest{Action: "test-action"},
			expectedCode:   http.StatusOK,
			expectedBody:   ActionResponse{Message: "Action received: test-action"},
			expectedHeader: "application/json",
		},
		{
			name:           "Invalid method GET",
			method:         http.MethodGet,
			url:            "/action/test-kind",
			payload:        ActionRequest{},
			expectedCode:   http.StatusMethodNotAllowed,
			expectedBody:   ActionResponse{},
			expectedHeader: "text/plain; charset=utf-8",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create the request body based on the test case
			var body *bytes.Buffer
			if tc.method == http.MethodPost {
				reqBody, _ := json.Marshal(tc.payload)
				body = bytes.NewBuffer(reqBody)
			} else {
				body = bytes.NewBuffer(nil)
			}

			// Create a new HTTP request based on the test case
			req, err := http.NewRequest(tc.method, tc.url, body)
			if err != nil {
				t.Fatalf("Could not create HTTP request: %v", err)
			}

			// Set content type to application/json if POST request
			if tc.method == http.MethodPost {
				req.Header.Set("Content-Type", "application/json")
			}

			// Create a response recorder to capture the handler's response
			rr := httptest.NewRecorder()

			// Create a dummy mux to handle requests
			mux := http.NewServeMux()
			mux.HandleFunc("/action/test-kind", handlerAction)

			// Call the handler with our recorder and request
			mux.ServeHTTP(rr, req)

			// Check that the status code matches the expected one
			if status := rr.Code; status != tc.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedCode)
			}

			// Check that the response content type is correct
			if contentType := rr.Header().Get("Content-Type"); contentType != tc.expectedHeader {
				t.Errorf("handler returned wrong content type: got %v want %v", contentType, tc.expectedHeader)
			}

			// Only check body for valid requests (200 OK)
			if tc.expectedCode == http.StatusOK {
				var response ActionResponse
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Could not decode response: %v", err)
				}

				// Check if the response message matches the expected message
				if response.Message != tc.expectedBody.Message {
					t.Errorf("handler returned unexpected body: got %v want %v",
						response.Message, tc.expectedBody.Message)
				}
			}
		})
	}
}
