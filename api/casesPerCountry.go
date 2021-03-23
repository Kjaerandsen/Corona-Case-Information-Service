package api

import (
	"fmt"
	"main/function"
	"net/http"
)

// Contains the input data from the mmg api, used twice for requests with a date range
// Once for data with cases and one for data with recovered statistics
type Data struct {
	Country			      	string 					`json:"country"`
	Population				int						`json:"population"`
	Continent       		string  				`json:"continent"`
	Dates 				    map[string]interface{} 	`json:"dates"`
}

// For the final output
type OutputData struct {
	Country              string  `json:"country"`
	Continent            string  `json:"continent"`
	Scope                string  `json:"scope"`
	Confirmed            int     `json:"confirmed"`
	Recovered            int     `json:"recovered"`
	PopulationPercentage float64 `json:"population_percentage"`
}

func CasesPerCountry(w http.ResponseWriter, r *http.Request) {
	// Parts array for the text split parts
	var parts []string
	var scopeParts []string
	// ErrorResponse bool for errors
	var errorResponse bool

	// Checks if the request is valid
	errorResponse, parts = function.TextSplitter(r.URL.Path, 5, "/")
	if !errorResponse {
		function.ErrorHandle(w, "Expected format: ../corona/v1/country?scope=start-end",
			400, "Bad request")
		return
	}
	countryName := parts[4]

	// Scope from the url of the request "?limit=date_start-date_end"
	var scope = r.FormValue("scope")

	// If a scope is not specified
	if scope == "" {
		casesWithoutScope(w,r)
		return
	}

	// Check the validity of the scope
	errorResponse, scopeParts = function.TextSplitter(scope, 6,"-")
	if !errorResponse{
		function.ErrorHandle(w, "Expected scope format ?scope=yyyy-mm-dd-yyyy-mm-dd",
			400, "Bad request")
		return
	}

	// Check if the dates are valid

	casesWithScope(w,r)

	fmt.Println("Name of the country: ", countryName, " Scope: ", scope, scopeParts[0], errorResponse)
}

// Handles the country cases without a scope
func casesWithoutScope(w http.ResponseWriter, r *http.Request) {
	// Perform the api call

}

// Handles the country cases with a scope
func casesWithScope(w http.ResponseWriter, r *http.Request) {
	// Perform the api call

}