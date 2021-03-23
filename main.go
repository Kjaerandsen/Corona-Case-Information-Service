package main

import (
	"fmt"
	"log"
	"main/function"
	"main/api"
	"net/http"
	"os"
)

// Main function, opens the port and sends the requests on
func main() {
	// Sets up the port of the application to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Sets the startvalue for the uptime resource in diag
	function.UptimeInit()

	// http request handlers
	http.HandleFunc("/corona/v1/country/", api.CasesPerCountry)
	http.HandleFunc("/corona/v1/diag/", api.Diag)
	http.HandleFunc("/corona/v1/policy/", function.NotImplemented)
	http.HandleFunc("/corona/v1/notifications/", function.NotImplemented)

	// redirect if missing the trailing slash
	http.HandleFunc("/corona/v1/country", function.Redirect)
	http.HandleFunc("/corona/v1/diag", function.Redirect)
	http.HandleFunc("/corona/v1/policy", function.Redirect)
	http.HandleFunc("/corona/v1/notifications", function.Redirect)

	fmt.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}