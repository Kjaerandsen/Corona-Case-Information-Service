package api

import (
	"main/function"
	"net/http"
)

// Main function that calls other functions
func Notifications(w http.ResponseWriter, r *http.Request) {
	var parts []string
	// ErrorResponse bool for errors
	var errorResponse bool

	// Checks if the request is valid & puts the path into the parts string array
	errorResponse, parts = function.TextSplitter(r.URL.Path, 5, "/")
	if !errorResponse {
		function.ErrorHandle(w, "Expected format: ../corona/v1/country?scope=start-end",
			400, "Bad request")
		return
	}

	// If the request is a get request
	if r.Method == http.MethodGet {
		// If the data after the slash is empty run the "View all registered webhooks" command
		if parts[4] == "" {
			webhookViewAll(w, r)
			return
		// Else view the specified webhook
		} else {
			webhookViewSingle(w, r, parts[4])
			return
		}
	}

	// If the request is a post request
	if r.Method == http.MethodPost {
		// Handle the post request
		webhookInvocate(w, r, parts[4])
		return
	}

	// If the request is a delete request
	if r.Method == http.MethodDelete {
		// Handle the delete request
		webhookDelete(w, r, parts[4])
		return
	}
}

// Views all webhooks
func webhookViewAll(w http.ResponseWriter, r *http.Request) {

}

// Views the specified webhook if it exists
func webhookViewSingle(w http.ResponseWriter, r *http.Request, name string) {

}

// Deletes the specified webhook
func webhookDelete(w http.ResponseWriter, r *http.Request, name string) {

}

func webhookInvocate(w http.ResponseWriter, r *http.Request, name string) {

}