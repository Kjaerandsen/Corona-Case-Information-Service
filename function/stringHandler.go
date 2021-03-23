package function

import (
	"fmt"
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

	fmt.Println(len(parts))

	if len(parts) != length {
		return false, parts
	}

	return true, parts
}

/*
	Checks if the input dates are ordered, uses the data from the textsplitter function
 */
func DateCheck(dateparts string) bool{
	// Check if the start date is before or equal to the end date
	if dateparts[0] > dateparts[3] {
		return false
	} else if dateparts[0] == dateparts[3] {
		if dateparts[1] > dateparts[4] {
			return false
		} else if dateparts[1] == dateparts[4] && dateparts[2] >= dateparts[5] {
			return false
		}
	}

	return true
}