package api

import (
	"bytes"
	"cloud.google.com/go/firestore" // Firestore-specific support
	"context"                       // State handling across API boundaries; part of native GoLang API
	"encoding/json"
	"fmt"
	"google.golang.org/api/iterator"
	"io/ioutil"
	"log"
	"main/function"
	"net/http"
	"strconv"
	"time"
)

// Firebase context and client used by Firestore functions throughout the program.
var Ctx context.Context
var Client *firestore.Client
// Collection name in Firestore
var Collection = "webhooks"
// Counts the amount of webhooks registered
var WebhookCount int
// Map that stores the registered webhooks locally, used for running their functionality
var Webhooks = make(map[string]WebhookData)

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
			webhookViewAll(w)
			return
		// Else view the specified webhook
		} else {
			webhookViewSingle(w, parts[4])
			return
		}
	case http.MethodPost:
		// Handle the post request
		webhookCreate(w, r)
		return
	case http.MethodDelete:
		// Handle the delete request
		webhookDelete(w, parts[4])
		return
	default:
		function.ErrorHandle(w, "Method not supported, only the " +
			"get, delete and put methods is supported for this endpoint",
			400, "Bad request")
	}
}

// Sets the WebhookCount var, and runs the webhook function for the registered webhooks at start
func WebhookStart() {
	var counter int
	var webhookTempData WebhookData

	// Retrieve the data from firestore
	iter := Client.Collection(Collection).Documents(Ctx) // Loop through all entries in collection "webhooks"
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		fmt.Println(doc)
		m := doc.Data()

		webhookTempData.Country = fmt.Sprintf("%v", m["country"])
		webhookTempData.Id = doc.Ref.ID
		webhookTempData.Url = fmt.Sprintf("%v", m["url"])
		webhookTempData.Information = fmt.Sprintf("%v", m["field"])
		webhookTempData.Trigger = fmt.Sprintf("%v", m["trigger"])
		webhookTempData.Timeout, err = strconv.Atoi(fmt.Sprintf("%v", m["timeout"]))
		if err != nil {
			log.Fatalf("Failed converting time from firestore database")
		}

		// Add it to the map
		Webhooks[doc.Ref.ID] = webhookTempData
		// Run the goroutine for it
		go webhookCheck(doc.Ref.ID)

		counter++
	}
	WebhookCount = counter
}

// Runs the webhook functionality each timeout seconds while it exists in the local webhook map
// Ran as a go routine from webhookCreate() and WebhookStart()
func webhookCheck(webhookId string) {
	var exists = true
	var webhook WebhookData
	// Data from the "stringency" information
	var StringencyData [2]OutputStrData
	// Data from the "confirmed" information
	var ConfirmedData [2]OutputData
	var err bool
	var Code string

	// Perform data integrity check here?
	webhook, exists = Webhooks[webhookId]
	// else exits the goroutine
	if !exists {
		return
	}

	// Initial data request, for comparison on first execution
	if webhook.Information == "stringency" {

		// Retrieve country code
		Code, err = CountryCodeWebhook(webhook.Country)
		if !err {
			fmt.Println("Error in retrieving country code for webhook with id: ",
				webhook.Id, "Stopping the routine.")
			return
		}
		// Stringency start data request
		StringencyData[0], err = StringencyWebhookWithoutScope(webhook.Country, Code)
		if !err {
			fmt.Println("Error in retrieving stringency data for webhook with id: ",
				webhook.Id, "Stopping the routine.")
			return
		}

	} else {

		// Confirmed start data request
		ConfirmedData[0], err = CasesWebhook(webhook.Country)
		if !err {
			fmt.Println("Error in retrieving initial data for webhook with id: ", webhook.Id, "Stopping the routine.")
			return
		}

	}

	// Perform initial request to hold the data
	fmt.Println("Webhook runner started for webhook with id: ", webhookId)
	for {
		// Checks if the map still contains the webhook
		_, exists = Webhooks[webhookId]
		// else exits the goroutine
		if !exists {
			break
		}
		// Sleep for the timeout amount of seconds
		time.Sleep(time.Duration(webhook.Timeout) * time.Second)

		if webhook.Information == "stringency" {
			// Stringency data request
			StringencyData[1], err = StringencyWebhookWithoutScope(webhook.Country, Code)
			if !err {
				fmt.Println("Error in retrieving stringency data for webhook with id: ",
					webhook.Id, "Stopping the routine.")
				return
			}
			// Compare it
			if StringencyData[0].Stringency != StringencyData[1].Stringency{
				// If the webhook is to trigger on timeout send the data
				if webhook.Trigger == "ON_TIMEOUT" {
					// Run the output
					webhookSendStringency(webhook.Url, StringencyData[1])
				}
				// Update the data
				StringencyData[0] = StringencyData[1]
			} else {
				// Update the data
				StringencyData[0] = StringencyData[1]
				// Run the output
				webhookSendStringency(webhook.Url, StringencyData[0])
			}
			// Handle output if changed or "ON_UPDATE"
		} else /* Cases */{
			// Request the data
			ConfirmedData[1], err = CasesWebhook(webhook.Country)
			if !err {
				fmt.Println("Error in retrieving data for webhook with id: ", webhook.Id, "Stopping the routine.")
				return
			}
			// Compare it
			if ConfirmedData[0].Confirmed != ConfirmedData[1].Confirmed{
				// If the webhook is to trigger on timeout send the data
				if webhook.Trigger == "ON_TIMEOUT" {
					// Run the output
					webhookSendConfirmed(webhook.Url, ConfirmedData[1])
				}
				// Update the data
				ConfirmedData[0] = ConfirmedData[1]
			} else {
				// Update the data
				ConfirmedData[0] = ConfirmedData[1]
				// Run the output
				webhookSendConfirmed(webhook.Url, ConfirmedData[0])
			}
		}
	}
}

// Creates a webhook
func webhookCreate(w http.ResponseWriter, r *http.Request) {
	// Expects incoming body in terms of WebhookRegistration struct
	webhook := WebhookData{}
	err := json.NewDecoder(r.Body).Decode(&webhook)
	if err != nil {
		function.ErrorHandle(w, "Bad request, see manual for specification of post",
			400, "Request")
	}

	id, _, err := Client.Collection("webhooks").Add(Ctx,
		map[string]interface{}{
		"url": webhook.Url,
		"timeout": webhook.Timeout,
		"field": webhook.Information,
		"country": webhook.Country,
		"trigger": webhook.Trigger,
		})
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error when adding message ", http.StatusBadRequest)
		return
	}

	// Add the webhook to the webhook map for running
	Webhooks[id.ID] = webhook
	// Run the webhook runner
	go webhookCheck(id.ID)

	WebhookCount ++
	fmt.Println("Webhook with id: " + id.ID + " has been registered.")
	http.Error(w, id.ID, http.StatusCreated)
}

// Views all webhooks
func webhookViewAll(w http.ResponseWriter) {
	var webhooks []WebhookData
	var outputString string

	w.Header().Set("Content-Type", "application/json")

	fmt.Println(Client.Collection(Collection).ID)

	// Retrieve the data from firestore
	var counter int

	_ , err := fmt.Fprintf(w, "[")
	if err != nil {
		http.Error(w, "Error while writing response body.", http.StatusInternalServerError)
	}

	iter := Client.Collection(Collection).Documents(Ctx) // Loop through all entries in collection "webhooks"
	for {
		outputString = ""
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		m := doc.Data() // A message map with string keys. Each key is a unique id

		outputString = fmt.Sprintf(`{"id":"%v","url":"%s","timeout":%v,`, doc.Ref.ID, m["url"], m["timeout"])
		outputString = fmt.Sprintf(`%s"information":"%s","country":"%s","trigger":"%s"},`,
			outputString,
			m["field"],
			m["country"],
			m["trigger"])

		_ , err = fmt.Fprintf(w, "%s", outputString)
		if err != nil {
			http.Error(w, "Error while writing response body.", http.StatusInternalServerError)
		}
		counter++
	}

	WebhookCount = counter
	_ , err = fmt.Fprintf(w, "]")
	if err != nil {
		http.Error(w, "Error while writing response body.", http.StatusInternalServerError)
	}

	// Return the struct array
	fmt.Println(webhooks)
}

// Views the specified webhook if it exists
func webhookViewSingle(w http.ResponseWriter, name string) {
	var outputDataStruct WebhookData

	data, err := Client.Collection(Collection).Doc(name).Get(Ctx)
	if err != nil {
		function.ErrorHandle(w, "No webhook with the specified id found", 400, "Bad request")
		return
	}

	outputData := data.Data()

	outputDataStruct.Country = fmt.Sprintf("%s", outputData["country"])
	outputDataStruct.Information = fmt.Sprintf("%s", outputData["field"])
	outputDataStruct.Timeout, err = strconv.Atoi(fmt.Sprintf("%v", outputData["timeout"]))
	if err != nil {
		function.ErrorHandle(w, "Internal server error", 500, "Parsing")
		return
	}
	outputDataStruct.Trigger = fmt.Sprintf("%s", outputData["trigger"])
	outputDataStruct.Url = fmt.Sprintf("%s", outputData["url"])
	outputDataStruct.Id = name

	// Exporting the data
	w.Header().Set("Content-Type", "application/json")
	// Converts the diagnosticData into json
	outData, _ := json.Marshal(outputDataStruct)
	// Writes the json
	_, err = w.Write(outData)
	if err != nil {
		function.ErrorHandle(w, "Internal server error", 500, "Response")
	}
}

// Deletes the specified webhook
func webhookDelete(w http.ResponseWriter, name string) {
	_, err := Client.Collection(Collection).Doc(name).Delete(Ctx)
	if err != nil {
		http.Error(w, "Deletion of " + name + " failed.", http.StatusInternalServerError)
	}

	// Remove the map value to stop the go routine
	delete(Webhooks, name)

	WebhookCount--
	http.Error(w, "Deletion of " + name + " successful.", http.StatusNoContent)
}

// Serves webhook confirmed data
func webhookSendConfirmed(url string, data OutputData){

	outData, err := json.Marshal(data)
	if err != nil{
		log.Printf("%v", "Error during json marshaling.")
	}

	req, err := http.NewRequest(http.MethodPost, url,
		bytes.NewReader(outData))
	if err != nil {
		log.Printf("%v", "Error during request creation.")
		return
	}
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error in HTTP request: " + err.Error())
		return
	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Something is wrong with invocation response: " + err.Error())
		return
	}

	fmt.Println("Webhook invoked. Received status code " + strconv.Itoa(res.StatusCode) +
		" and body: " + string(response))
}

// Send webhook Stringency data
func webhookSendStringency(url string, data OutputStrData){

	outData, err := json.Marshal(data)
	if err != nil{
		log.Printf("%v", "Error during json marshaling.")
	}

	req, err := http.NewRequest(http.MethodPost, url,
		bytes.NewReader(outData))
	if err != nil {
		log.Printf("%v", "Error during request creation.")
		return
	}
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error in HTTP request: " + err.Error())
		return
	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Something is wrong with invocation response: " + err.Error())
		return
	}

	fmt.Println("Webhook invoked. Received status code " + strconv.Itoa(res.StatusCode) +
		" and body: " + string(response))
}