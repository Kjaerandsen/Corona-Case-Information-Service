package function

import (
	"fmt"
	"net/http"
)

// Handles errors in the program by sending proper http responses and messages
func ErrorHandle(w http.ResponseWriter, errorMsg string, errorCode int, errorType string) {
	fmt.Println("Mistakes were made")
	// Set the content type to json
	w.Header().Set("Content-Type", "application/json")
	// Writes the error code
	w.WriteHeader(errorCode)
	// Sends the error code and message
	_, err := w.Write([]byte("Error: " + errorType + ": " + errorMsg))
	if err != nil {
		fmt.Println("Error in encoding error message for errorhandling in function.Errorhandle")
		return
	}
}