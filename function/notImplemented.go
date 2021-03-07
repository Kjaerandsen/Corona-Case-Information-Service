package function

import "net/http"

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	_, err := w.Write([]byte("501 Not Implemented"))
	if err != nil {
		// TODO add error message here
	}
}