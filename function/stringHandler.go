package function

import (
	"strings"
)

/*
	Splits the url into parts and checks if it is the correct length
	Returns the split url as a string array
*/
func TextSplitter(text string, length int, seperator string) (bool, []string){
	/*
		Modified code based on code retrieved from the "RESTstudent" example at
		"https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021/-/blob/master/RESTstudent/cmd/students_server.go"
		Retrieves the country name from the url after the trailing slash
	*/
	parts := strings.Split(text, seperator)

	if len(parts) != length {
		return false, parts
	}

	return true, parts
}

// Check if the characters of a string contains only numbers (check if string is an integer)
// TODO: THIS?

// Make the first char of a string uppercase
// TODO: THIS?