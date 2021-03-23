package function

import (
	"log"
	"net/http"
)

// Returns the http status code from the parameter url, used in the diag endpoint
func GetHttpStatus(webUrl string) int {
	// Get http status code
	resp, err := http.Get(webUrl)
	if err != nil {
		// TODO: Look over error handling
		log.Fatal(err)
		return 500
	}
	return resp.StatusCode
}