package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/function"
	"net/http"
	"strconv"
	"strings"
)

// TODO: Rename this file? move api to /v1/ ?
// corona/v1/country/ endpoint

// Contains the input data from the mmg api, used twice for requests with a date range
// Once for data with cases and one for data with recovered statistics
type InputDataWithScope struct {
	All struct {
		Country    		string                  `json:"country"`
		Population 		int                     `json:"population"`
		Continent 		string                  `json:"continent"`
		Dates      		map[string]interface{}  `json:"dates"`
	}
}

// For the input data of the request without a defined scope
type InputDataWithoutScope struct {
	All struct {
		Confirmed		int					   	`json:"confirmed"`
		Recovered		int					   	`json:"recovered"`
		Country			string     				`json:"country"`
		Continent		string					`json:"continent"`
		Population		int						`json:"population"`
	}
}

// For the final output
type OutputData struct {
	Country              string  				`json:"country"`
	Continent            string  				`json:"continent"`
	Scope                string  				`json:"scope"`
	Confirmed            int     				`json:"confirmed"`
	Recovered            int    				`json:"recovered"`
	PopulationPercentage float64 				`json:"population_percentage"`
}

// Main function that calls other functions
func CasesPerCountry(w http.ResponseWriter, r *http.Request) {
	// Parts array for the text split parts
	var parts []string
	var scopeParts []string
	// ErrorResponse bool for errors
	var errorResponse bool

	// Checks if the request method is valid (only accepts get)
	if r.Method != http.MethodGet {
		function.ErrorHandle(w, "Method not supported, only get method is supported for this endpoint",
			400, "Bad request")
		return
	}

	// Checks if the request is valid
	errorResponse, parts = function.TextSplitter(r.URL.Path, 5, "/")
	if !errorResponse {
		function.ErrorHandle(w, "Expected format: ../corona/v1/country?scope=start-end",
			400, "Bad request")
		return
	}
	countryName := parts[4]
	// Capitalize the first letter of each word
	countryName = strings.Title(countryName)

	// Scope from the url of the request "?limit=date_start-date_end"
	var scope = r.FormValue("scope")

	// If a scope is not specified
	if scope == "" {
		casesWithoutScope(w,r, countryName)
		return
	}

	// Check the validity of the scope
	errorResponse, scopeParts = function.TextSplitter(scope, 6,"-")
	if !errorResponse {
		function.ErrorHandle(w, "Expected scope format ?scope=yyyy-mm-dd-yyyy-mm-dd",
			400, "Bad request")
		return
	}

	// Check if the dates are valid
	if !function.StartEnd(scopeParts) {
		function.ErrorHandle(w, "Expected scope format ?scope=yyyy-mm-dd-yyyy-mm-dd, parsing",
			400, "Bad request")
		return
	}

	casesWithScope(w,r, countryName, scope, scopeParts)
}

// Handles the country cases without a scope
func casesWithoutScope(w http.ResponseWriter, r *http.Request, name string) {
	var data InputDataWithoutScope
	var outputData OutputData

	// Perform the api call
	url := fmt.Sprintf("https://covid-api.mmediagroup.fr/v1/cases?country=%s", name)

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		function.ErrorHandle(w, "Error in sending request to external api", 500, "Request")
		return
	}

	// Setting content type -> effect depends on the service provider
	r.Header.Add("content-type", "application/json")
	// Instantiate the client
	client := &http.Client{}

	// Issue request
	res, err := client.Do(r)
	if err != nil {
		function.ErrorHandle(w, "Error in parsing json from external api", 500, "Parsing")
		return
	}

	// If the http statuscode retrieved from the api is not 200 / "OK"
	if res.StatusCode != 200 {
		function.ErrorHandle(w, "Error in sending request to external api", 404, "Request")
		return
	}

	// Read the data
	output, err := ioutil.ReadAll(res.Body)
	if err != nil {
		function.ErrorHandle(w, "Error in parsing json from external api", 500, "Parsing")
		return
	}

	// JSON into struct
	err = json.Unmarshal(output, &data)
	if err != nil {
		function.ErrorHandle(w, "Error in parsing json from external api", 500, "Parsing")
		return
	}

	if data.All.Country == "" {
		function.ErrorHandle(w, "Country doesn't exist in the covid database",
			400, "Request")
		return
	}

	//fmt.Println(data.All.Confirmed, data.All.Continent, data.All.Country, data.All.Population, data.All.Recovered)
	// Set the output data with values from the api request
	outputData.Recovered = data.All.Recovered
	outputData.Country = data.All.Country
	outputData.Continent = data.All.Country
	outputData.Confirmed = data.All.Confirmed
	// TODO: Fix this, probably doesn't display enough comma values
	outputData.PopulationPercentage, err =
		strconv.ParseFloat(fmt.Sprintf("%.2f",float64(data.All.Confirmed) / float64(data.All.Population)),64)
	if err != nil {
		function.ErrorHandle(w, "Error in handling float for population percentage",
			500, "Internal")
		return
	}
	outputData.Scope = "total"

	// Returns the output data to the user
	returnData(w,outputData)
}

// Handles the country cases with a scope
func casesWithScope(w http.ResponseWriter, r *http.Request, name string, scope string, scopeParts []string) {

	// Perform the api call for cases of covid
	//https://covid-api.mmediagroup.fr/v1/history?country=France&status=Confirmed
	var data InputDataWithScope	// Used for handling the input data
	var outputData OutputData	// Used for handling the output data
	var key string 				// Used as the key when traversing the date map
	var outValue int				// Used for cases and recovered when adding to the output data

	// Perform the api call
	url := fmt.Sprintf("https://covid-api.mmediagroup.fr/v1/history?country=%s&status=Confirmed", name)

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		function.ErrorHandle(w, "Error in sending request to external api", 500, "Request")
		return
	}

	// Setting content type -> effect depends on the service provider
	r.Header.Add("content-type", "application/json")
	// Instantiate the client
	client := &http.Client{}

	// Issue request
	res, err := client.Do(r)
	if err != nil {
		function.ErrorHandle(w, "Error in parsing json from external api", 500, "Parsing")
		return
	}

	// If the http statuscode retrieved from the api is not 200 / "OK"
	if res.StatusCode != 200 {
		function.ErrorHandle(w, "Error in sending request to external api", 404, "Request")
		return
	}

	// Print output
	output, err := ioutil.ReadAll(res.Body)
	if err != nil {
		function.ErrorHandle(w, "Error in parsing json from external api", 500, "Parsing")
		return
	}

	// JSON into struct
	err = json.Unmarshal(output, &data)
	if err != nil {
		function.ErrorHandle(w, "Error in parsing json from external api", 500, "Parsing")
		return
	}

	// If the country doesn't exist exit
	if data.All.Country == "" {
		function.ErrorHandle(w, "Country doesn't exist in the covid database",
			400, "Request")
		return
	}

	key = fmt.Sprintf("%s-%s-%s", scopeParts[0], scopeParts[1], scopeParts[2])

	// Check the data for the first date
	value, exists := data.All.Dates[key]
	if !exists {
		function.ErrorHandle(w, "Scope date 1 not found in covid dataset, first date reported is 2020-01-22",
			404, "Not found")
		return
	}

	outValue, err = strconv.Atoi(fmt.Sprintf("%v" ,value))
	if err != nil {
		function.ErrorHandle(w, "Internal Server Error when handling dataset", 500, "Internal")
		return
	}

	outputData.Confirmed -= outValue

	key = fmt.Sprintf("%s-%s-%s", scopeParts[3], scopeParts[4], scopeParts[5])

	// Check the data for the second date
	value, exists = data.All.Dates[key]
	if !exists {
		function.ErrorHandle(w,
			"Scope date 2 not found in covid dataset, please try yesterdays date if today is not yet reported",
			404, "Not found")
		return
	}

	outValue, err = strconv.Atoi(fmt.Sprintf("%v", value))
	if err != nil {
		function.ErrorHandle(w, "Internal Server Error when handling dataset", 500, "Internal")
		return
	}

	outputData.Confirmed += outValue

	// Perform the api call for recovered cases
	//https://covid-api.mmediagroup.fr/v1/history?country=%s&status=Recovered
	// Perform the api call
	url = fmt.Sprintf("https://covid-api.mmediagroup.fr/v1/history?country=%s&status=Recovered", name)

	r, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		function.ErrorHandle(w, "Error in sending request to external api", 500, "Request")
		return
	}

	// Setting content type -> effect depends on the service provider
	r.Header.Add("content-type", "application/json")
	// Instantiate the client
	client = &http.Client{}

	// Issue request
	res, err = client.Do(r)
	if err != nil {
		function.ErrorHandle(w, "Error in parsing json from external api", 500, "Parsing")
		return
	}

	// If the http statuscode retrieved from the api is not 200 / "OK"
	if res.StatusCode != 200 {
		function.ErrorHandle(w, "Error in sending request to external api", 404, "Request")
		return
	}

	// Print output
	output, err = ioutil.ReadAll(res.Body)
	if err != nil {
		function.ErrorHandle(w, "Error in parsing json from external api", 500, "Parsing")
		return
	}

	// JSON into struct
	err = json.Unmarshal(output, &data)
	if err != nil {
		function.ErrorHandle(w, "Error in parsing json from external api", 500, "Parsing")
		return
	}

	// If the country doesn't exist exit
	if data.All.Country == "" {
		function.ErrorHandle(w, "Country doesn't exist in the covid database",
			400, "Request")
		return
	}

	key = fmt.Sprintf("%s-%s-%s", scopeParts[0], scopeParts[1], scopeParts[2])

	// Check the data for the first date
	value, exists = data.All.Dates[key]
	if !exists {
		function.ErrorHandle(w, "Scope date 1 not found in covid dataset, first date reported is 2020-01-22",
			404, "Not found")
		return
	}

	outValue, err = strconv.Atoi(fmt.Sprintf("%v" ,value))
	if err != nil {
		function.ErrorHandle(w, "Internal Server Error when handling dataset", 500, "Internal")
		return
	}

	outputData.Recovered -= outValue

	key = fmt.Sprintf("%s-%s-%s", scopeParts[3], scopeParts[4], scopeParts[5])

	// Check the data for the second date
	value, exists = data.All.Dates[key]
	if !exists {
		function.ErrorHandle(w,
			"Scope date 2 not found in covid dataset, please try yesterdays date if today is not yet reported",
			404, "Not found")
		return
	}

	outValue, err = strconv.Atoi(fmt.Sprintf("%v", value))
	if err != nil {
		function.ErrorHandle(w, "Internal Server Error when handling dataset", 500, "Internal")
		return
	}

	outputData.Recovered += outValue

	// Add the data from the api call to outputData here
	outputData.Country = data.All.Country
	outputData.Continent = data.All.Continent
	// TODO: Fix this, probably doesn't display enough comma values
	outputData.PopulationPercentage, err =
		strconv.ParseFloat(fmt.Sprintf("%.2f", float64(outputData.Confirmed) / float64(data.All.Population)), 64)
	if err != nil {
		function.ErrorHandle(w, "Error in handling float for population percentage",
			500, "Internal")
		return
	}
	outputData.Scope = scope

	// Returns the output data to the user
	returnData(w, outputData)
}

// Simply returns the output data struct it receives as json to the user
func returnData(w http.ResponseWriter, data OutputData) {
	// Exporting the data
	w.Header().Set("Content-Type", "application/json")
	// Converts the diagnosticData into json
	outData, _ := json.Marshal(data)
	// Writes the json
	_, err := w.Write(outData)
	if err != nil {
		function.ErrorHandle(w, "Internal server error", 500, "Response")
	}
}