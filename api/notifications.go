package api

import (
	"encoding/json"
	"fmt"
	"main/function"
	"net/http"
	"strconv"
)

// Webhook data placeholder
var WebhooksMap = make(map[string]WebhookData)

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

	// Checks the method of the request and handles it accordingly
	switch r.Method {
	case http.MethodGet:
		// If the data after the slash is empty run the "View all registered webhooks" command
		if parts[4] == "" {
			webhookViewAll(w, r)
			return
		// Else view the specified webhook
		} else {
			webhookViewSingle(w, parts[4])
			return
		}
	case http.MethodPost:
		// Handle the post request
		webhookCreate(w, r, parts[4])
		return
	case http.MethodDelete:
		// Handle the delete request
		webhookDelete(w, r, parts[4])
		return
	default:
		function.ErrorHandle(w, "Method not supported, only the " +
			"get, delete and put methods is supported for this endpoint",
			400, "Bad request")
	}
}

func webhookCreate(w http.ResponseWriter, r *http.Request, name string) {
	// Expects incoming body in terms of WebhookRegistration struct
	webhook := WebhookData{}
	err := json.NewDecoder(r.Body).Decode(&webhook)
	if err != nil {
		function.ErrorHandle(w, "Bad request, see manual for specification of post",
			400, "Request")
	}
	// Check the data and add an id
	webhook.Id = fmt.Sprintf("%v", len(WebhooksMap))

	fmt.Println(webhook)

	// TODO: check if already exists first
	WebhooksMap[webhook.Id] = webhook

	fmt.Println("Webhook " + webhook.Url + " has been registered.")
	//http.Error(w, strconv.Itoa(len(Webhooks)-1), http.StatusCreated)
	http.Error(w, strconv.Itoa(len(WebhooksMap)-1), http.StatusCreated)
}

// Views all webhooks
func webhookViewAll(w http.ResponseWriter, r *http.Request) {
	// Return all webhooks
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(WebhooksMap)
	if err != nil {
		function.ErrorHandle(w, "Internal server error",
			500, "Internal")
	}
}

// Views the specified webhook if it exists
func webhookViewSingle(w http.ResponseWriter, name string) {

	// If the map key exists return the value
	output, exists := WebhooksMap[name]
	if !exists{
		// If the id doesn't exist in the database return an error
		function.ErrorHandle(w, "There is no webhook with the specified id registered",
			400, "Bad request")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(output)
	if err != nil {
		function.ErrorHandle(w, "Internal server error",
			500, "Internal")
	}

}


// Deletes the specified webhook
func webhookDelete(w http.ResponseWriter, r *http.Request, name string) {
	// Check if the key exists
	_, exists := WebhooksMap[name]
	if !exists{
		// If the id doesn't exist in the database return an error
		function.ErrorHandle(w, "There is no webhook with the specified id registered",
			400, "Bad request")
		return
	}
	// Remove it from the database
	delete(WebhooksMap, name)
	http.Error(w, "Webhook: " + name + " deleted", http.StatusOK)
}