package function

import "net/http"

// Redirects urls missing the trailing slash
func Redirect(w http.ResponseWriter, r *http.Request){
	http.Redirect(w, r, r.URL.Path + "/", http.StatusSeeOther)
}

// Returns a not implemented code, used when building endpoints
func NotImplemented(w http.ResponseWriter, r *http.Request) {
	ErrorHandle(w, "501 Not Implemented", 501, "Implementation")
}