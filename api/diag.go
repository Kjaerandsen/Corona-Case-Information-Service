package api

import (
	"encoding/json"
	"fmt"
	"main/function"
	"net/http"
	"strings"
)

// Handles the diag /corona/v1/diag/ request
func Diag(w http.ResponseWriter, r *http.Request) {

	// Checks if the method is get, if not sends an error
	// TODO: turn this into a separate function
	if r.Method != http.MethodGet {
		function.ErrorHandle(w,
			"Method not allowed, only get is allowed on this resource", 405, "Method")
	}

	// Checks if the url is correct
	// TODO: Turn this into a separate function with a parameter expected message
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 5 {
		function.ErrorHandle(w, "500 Internal Server Error", 500, "Internal")
		return
	}

	// Creates the diagnostic information
	w.Header().Set("Content-Type", "application/json")
	diagnosticData := &Diagnostic{
		// TODO: replace the method for this request as it returns a 403 forbidden with head request
		Mmediagroupapi: fmt.Sprintf("%d", function.GetHttpStatus("https://covid-api.mmediagroup.fr/v1/cases")),
		Covidtrackerapi: fmt.Sprintf("%d", function.GetHttpStatus("https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/date-range/2021-03-02/2021-03-19")),
		Register: len(Webhooks),
		Version: "v1",
		Uptime: fmt.Sprintf("%ds", int(function.Uptime().Seconds())),
	}

	// Converts the diagnosticData into json
	data, _ := json.Marshal(diagnosticData)
	// Writes the json
	_, err := w.Write(data)
	// Error handling with code response
	if err != nil {
		function.ErrorHandle(w, "500 Internal Server Error", 500, "Internal")
	}
}