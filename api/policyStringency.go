package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"main/function"
	"net/http"
	"strconv"
)

// corona/v1/stringency endpoint

// Handle the request
func PolicyStringency(w http.ResponseWriter, r *http.Request) {

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
		function.ErrorHandle(w, "Expected format: ../corona/v1/policy/country_name?scope=start-end",
			400, "Bad request")
		return
	}
	countryName := parts[4]

	// Retrieve the alpha3code of the country
	countryCodeData, err := countryCode(w, r, countryName)
	if !err {
		return
	}

	fmt.Println(countryCodeData)

	// Scope from the url of the request "?limit=date_start-date_end"
	var scope = r.FormValue("scope")

	// If a scope is not specified
	if scope == "" {
		stringencyWithoutScope(w,r, countryName, countryCodeData[0].Alpha3Code)
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

	stringencyWithScope(w,r, countryName, scope, scopeParts, countryCodeData[0].Alpha3Code)

	fmt.Println(scopeParts)
}

// Get the country code and name from the country name using the restcountries api
func countryCode(w http.ResponseWriter, r *http.Request, name string) (CountryInfo, bool){
	// Perform the request to the api

	fmt.Println(name)

	/*
		Url request code based on RESTclient found at
		"https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021/-/blob/master/RESTclient/cmd/main.go"
	*/
	var data CountryInfo
	url := fmt.Sprintf("https://restcountries.eu/rest/v2/name/%s?fields=name;alpha3Code", name)

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		function.ErrorHandle(w, "Error in sending request to external api", 500, "Request")
		return data, false
	}

	// Setting content type -> effect depends on the service provider
	r.Header.Add("content-type", "application/json")
	// Instantiate the client
	client := &http.Client{}

	// Issue request
	res, err := client.Do(r)
	if err != nil {
		function.ErrorHandle(w, "Error in parsing json from external api", 500, "Parsing")
		return data, false
	}

	// If the http statuscode retrieved from restcountries is not 200 / "OK"
	if res.StatusCode != 200 {
		function.ErrorHandle(w, "Error in sending request to external api, country name probably wrong",
			404, "Request")
		return data, false
	}

	// Reading the data
	output, err := ioutil.ReadAll(res.Body)
	if err != nil {
		function.ErrorHandle(w, "Error in parsing json from external api", 500, "Parsing")
		return data, false
	}

	// JSON into struct
	err = json.Unmarshal(output, &data)
	if err != nil {
		function.ErrorHandle(w, "Error in parsing json from external api", 500, "Parsing")
		return data, false
	}

	if data[0].Alpha3Code == "" {
		// TODO: Rename this error
		function.ErrorHandle(w, "Received no data from external api", 500, "Request")
		return data, false
	}

	// Return the data and true if the request was successfull
	return data, true
}

// Function that performs the main request, returns the input data as a struct and a bool if there is an error
func stringencyRequest(w http.ResponseWriter, r *http.Request, date string, alpha string) (StringencyData, bool) {
	var data StringencyData // The data received from the api, and returned

	// Perform the api call
	/*
		Url request code based on RESTclient found at
		"https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021/-/blob/master/RESTclient/cmd/main.go"
	*/
	url := fmt.Sprintf("https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/actions/%s/%s", alpha, date)

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		function.ErrorHandle(w, "Error in sending request to external api", 500, "Request")
		return data, false
	}

	// Setting content type -> effect depends on the service provider
	r.Header.Add("content-type", "application/json")
	// Instantiate the client
	client := &http.Client{}

	// Issue request
	res, err := client.Do(r)
	if err != nil {
		function.ErrorHandle(w, "Error in parsing json from external api", 500, "Parsing")
		return data, false
	}

	// If the http statuscode retrieved from api is not 200 / "OK"
	if res.StatusCode != 200 {
		function.ErrorHandle(w, "Error in sending request to external api, country name probably wrong",
			404, "Request")
		return data, false
	}

	// Reading the data
	output, err := ioutil.ReadAll(res.Body)
	if err != nil {
		function.ErrorHandle(w, "Error in parsing json from external api", 500, "Parsing")
		return data, false
	}

	// JSON into struct
	err = json.Unmarshal(output, &data)
	if err != nil {
		function.ErrorHandle(w, "Error in parsing json from external api", 500, "Parsing")
		return data, false
	}

	// If no errors return the data received and a true bool as in executed successfully
	return data, true
}

// Handles the stringency cases without a scope
func stringencyWithoutScope(w http.ResponseWriter, r *http.Request, name string, alpha string) {
	var output OutputStrData
	var dates = function.DateToday()
	var data StringencyData

	// Perform the api requests
	data, err := stringencyDataRequest(w, r, dates, alpha)
	if !err {
		function.ErrorHandle(w, "Found no data in the stringency api for seven, 10 and 13 days ago",
			400, "Request")
		return
	}

	output.Stringency = data.Data.StringencyActual
	if output.Stringency == 0 {
		output.Stringency = data.Data.Stringency
	}
	output.Scope= "total"
	output.Trend = 0
	output.Country = name

	// Returns the output data to the user
	returnStringencyData(w, output)
}

// Handles the stringency cases with a scope
func stringencyWithScope(w http.ResponseWriter, r *http.Request, name string, scope string,
	scopeParts []string, alpha string) {
	var output OutputStrData
	var outputTrend float64
	var errorCheck error

	// Perform the api request for the start date
	data, err := stringencyRequest(w, r, fmt.Sprintf("%s-%s-%s", scopeParts[0], scopeParts[1], scopeParts[2]),
		alpha)
	if !err {
		function.ErrorHandle(w, "Internal Server Error when handling dataset", 500, "Internal")
		return
	}

	// Perform the api call for the end date
	// Perform the api request for the start date
	data2, err := stringencyRequest(w, r, fmt.Sprintf("%s-%s-%s", scopeParts[3], scopeParts[4], scopeParts[5]),
		alpha)
	if !err {
		function.ErrorHandle(w, "Internal Server Error when handling dataset", 500, "Internal")
		return
	}

	// If one of the dates contain no data default to -1 as stringency value and 0 as trend value
	if data.Data.Message != "" || data2.Data.Message != "" {
		output.Stringency = -1
		output.Trend = 0
	} else {
		// If StringencyActual is empty default to Stringency
		output.Stringency = data2.Data.StringencyActual
		if output.Stringency == 0 {
			output.Stringency = data2.Data.Stringency
		}

		// If StringencyActual is empty default to Stringency
		outputTrend = data.Data.StringencyActual
		if outputTrend == 0 {
			outputTrend = data.Data.Stringency
		}

		output.Trend = output.Stringency - outputTrend
	}

	// Limit the decimals in the trend float
	output.Trend, errorCheck = strconv.ParseFloat(fmt.Sprintf("%.2f", output.Trend), 64)
	if errorCheck != nil {
		function.ErrorHandle(w, "Error in handling float for trend data",
			500, "Internal")
		return
	}

	output.Scope= scope
	output.Country = name

	// Returns the output data to the user
	returnStringencyData(w, output)
}

// Loops through the dates and checks for valid data, returns valid data if found, else empty data and a false bool
func stringencyDataRequest(w http.ResponseWriter, r *http.Request, dates [3]string, alpha string) (StringencyData, bool) {
	for i:=0; i<3; i++ {
		data, err := stringencyRequest(w, r, dates[i], alpha)

		// If data.Data.Message is empty that means that the datapoint exists for the date
		// If the datapoint exists and there is no error return the data and go on
		if data.Data.Message == "" && err {
			return data, true
		}
	}
	return StringencyData{}, false
}

// Simply returns the output data struct it receives as json to the user
func returnStringencyData(w http.ResponseWriter, data OutputStrData) {
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

// Webhook version

// Get the country code and name from the country name using the restcountries api
func CountryCodeWebhook(name string) (string, bool){
	/*
		Url request code based on RESTclient found at
		"https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021/-/blob/master/RESTclient/cmd/main.go"
	*/
	var data CountryInfo
	url := fmt.Sprintf("https://restcountries.eu/rest/v2/name/%s?fields=name;alpha3Code", name)

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("Error in sending request to external api")
		return "", false
	}

	// Setting content type -> effect depends on the service provider
	r.Header.Add("content-type", "application/json")
	// Instantiate the client
	client := &http.Client{}

	// Issue request
	res, err := client.Do(r)
	if err != nil {
		log.Println("Error in parsing json from external api")
		return "", false
	}

	// If the http statuscode retrieved from restcountries is not 200 / "OK"
	if res.StatusCode != 200 {
		log.Println("Error in sending request to external api, country name probably wrong")
		return "", false
	}

	// Reading the data
	output, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Error in parsing json from external api")
		return "", false
	}

	// JSON into struct
	err = json.Unmarshal(output, &data)
	if err != nil {
		log.Println("Error in parsing json from external api")
		return "", false
	}

	if data[0].Alpha3Code == "" {
		log.Println("Received no data from external api")
		return "", false
	}

	// Return the data and true if the request was successfull
	return data[0].Alpha3Code, true
}

// Handles the stringency cases without a scope
func StringencyWebhookWithoutScope(name string, alpha string) (OutputStrData, bool){
	var output OutputStrData
	var dates = function.DateToday()
	var data StringencyData

	// Perform the api requests
	data, err := stringencyWebhookDataRequest(dates, alpha)
	if !err {
		log.Println("Found no data in the stringency api for seven, 10 and 13 days ago")
		return OutputStrData{}, false
	}

	output.Stringency = data.Data.StringencyActual
	if output.Stringency == 0 {
		output.Stringency = data.Data.Stringency
	}
	output.Scope= "total"
	output.Trend = 0
	output.Country = name

	// Returns the output data to the user
	return output, true
}

// Loops through the dates and checks for valid data, returns valid data if found, else empty data and a false bool
func stringencyWebhookDataRequest(dates [3]string, alpha string) (StringencyData, bool) {
	for i:=0; i<3; i++ {
		data, err := stringencyWebhookRequest(dates[i], alpha)

		// If data.Data.Message is empty that means that the datapoint exists for the date
		// If the datapoint exists and there is no error return the data and go on
		if data.Data.Message == "" && err {
			return data, true
		}
	}
	return StringencyData{}, false
}

// Function that performs the main request, returns the input data as a struct and a bool if there is an error
func stringencyWebhookRequest(date string, alpha string) (StringencyData, bool) {
	var data StringencyData // The data received from the api, and returned

	// Perform the api call
	/*
		Url request code based on RESTclient found at
		"https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021/-/blob/master/RESTclient/cmd/main.go"
	*/
	url := fmt.Sprintf("https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/actions/%s/%s", alpha, date)

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("Error in sending request to external api")
		return data, false
	}

	// Setting content type -> effect depends on the service provider
	r.Header.Add("content-type", "application/json")
	// Instantiate the client
	client := &http.Client{}

	// Issue request
	res, err := client.Do(r)
	if err != nil {
		log.Println("Error in parsing json from external api")
		return data, false
	}

	// If the http statuscode retrieved from api is not 200 / "OK"
	if res.StatusCode != 200 {
		log.Println("Error in sending request to external api, country name probably wrong")
		return data, false
	}

	// Reading the data
	output, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Error in parsing json from external api")
		return data, false
	}

	// JSON into struct
	err = json.Unmarshal(output, &data)
	if err != nil {
		log.Println("Error in parsing json from external api")
		return data, false
	}

	// If no errors return the data received and a true bool as in executed successfully
	return data, true
}
