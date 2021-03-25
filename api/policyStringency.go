package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/function"
	"net/http"
)

// corona/v1/stringency endpoint

// For the country name and alpha 3 code retrieved from the restcountries api
type CountryInfo []struct{
	Name       			string `json:"name"`
	Alpha3Code 			string `json:"alpha3Code"`
}

// For the final output
type OutputStrData struct {
	Country 			string  `json:"country"`
	Scope               string  `json:"scope"`
	Stringency			float64 `json:"stringency"`
	Trend				float64 `json:"trend"`
}

// For the input data of the request to the "covidtrackerapi"
type stringencyData struct {
	Data				map[string]interface{} `json:"data"`
}

// Handle the request
func PolicyStringency(w http.ResponseWriter, r *http.Request) {
	// Parts array for the text split parts
	var parts []string
	var scopeParts []string
	// ErrorResponse bool for errors
	var errorResponse bool

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
		//stringencyWithoutScope(w,r, countryName)
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

	//stringencyWithScope(w,r, countryName, scope, scopeParts)

	fmt.Println(scopeParts)
}

// Get the country code and name from the country name using the restcountries api
func countryCode(w http.ResponseWriter, r *http.Request, name string) (CountryInfo, bool){
	// Perform the request to the api

	/*
		Url request code based on RESTclient found at
		"https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021/-/blob/master/RESTclient/cmd/main.go"
		URL to invoke
	*/
	var data CountryInfo
	url := fmt.Sprintf("https://restcountries.eu/rest/v2/name/%s?fields=borders;currencies", name)

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
	//https://restcountries.eu/rest/v2/name/%s?fields=name;alpha3Code, name
	return data, true
}

// Handles the stringency cases with a scope
func stringencyWithScope(w http.ResponseWriter, r *http.Request) {

}

// Handles the stringency cases without a scope
func stringencyWithoutScope(w http.ResponseWriter, r *http.Request) {

}