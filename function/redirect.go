package function

import "net/http"

// Redirects urls missing the trailing slash
func Redirect(w http.ResponseWriter, r *http.Request){
	http.Redirect(w, r, r.URL.Path + "/", http.StatusSeeOther)
}