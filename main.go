package main

import (
	"context" // State handling across API boundaries; part of native GoLang API
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"main/api" // /v1?
	"main/function"
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

	// Connects to the firestore database
	api.Ctx = context.Background()
	sa := option.WithCredentialsFile("./service-account.json")
	app, err := firebase.NewApp(api.Ctx, nil, sa)
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}

	api.Client, err = app.Firestore(api.Ctx)

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	defer api.Client.Close()

	// Sets the global variable for the amount of registered webhooks
	api.WebhookStart()

	// http request handlers
	http.HandleFunc("/corona/v1/country/", api.CasesPerCountry)
	http.HandleFunc("/corona/v1/diag/", api.Diag)
	http.HandleFunc("/corona/v1/policy/", api.PolicyStringency)
	http.HandleFunc("/corona/v1/notifications/", api.Notifications)

	// redirect if missing the trailing slash
	http.HandleFunc("/corona/v1/country", function.Redirect)
	http.HandleFunc("/corona/v1/diag", function.Redirect)
	http.HandleFunc("/corona/v1/policy", function.Redirect)
	http.HandleFunc("/corona/v1/notifications", function.Redirect)

	fmt.Println("Listening on port " + port)
	fmt.Println("Endpoints available: ")
	fmt.Println("/corona/v1/country")
	fmt.Println("/corona/v1/policy")
	fmt.Println("/corona/v1/diag")
	fmt.Println("/corona/v1/notifications")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}