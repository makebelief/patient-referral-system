package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"

)

// Index handles requests to "/" and "/Home"
func Index(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	// case "/login", "/Login", "/LOGIN", "/ABOUT US", "/about us":
	// 	if r.URL.Path == "/login" || r.URL.Path == "/Login" || r.URL.Path == "LOGIN"  {
	// 		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	// 	}
	// 	if r.Method == http.MethodGet {
	// 		serveTemplate(w, "templates/login.html")
	// 	}
	case "/index", "/Index", "/INDEX":
		// if r.URL.Path == "/index" || r.URL.Path == "/Index" || r.URL.Path == "INDEX"  {
		// 	http.Redirect(w, r, "/index", http.StatusMovedPermanently)
		// }
		if r.Method == http.MethodGet {
			serveTemplate(w, "templates/index.html")
		}
	case "/", "/Home", "/home", "/HOME", "/templates/home.html":
		if r.URL.Path == "/" || r.URL.Path == "/Home" || r.URL.Path == "/HOME" {
			http.Redirect(w, r, "/home", http.StatusMovedPermanently)
		}
		if r.Method == http.MethodGet {
			serveTemplate(w, "templates/home.html")
		}
	case "/make_referral", "/Make_referral", "/templates/make_referral.html":
		// if r.URL.Path == "/index" || r.URL.Path == "/Index" || r.URL.Path == "INDEX"  {
		// 	http.Redirect(w, r, "/index", http.StatusMovedPermanently)
		// }
		if r.Method == http.MethodGet {
			serveTemplate(w, "templates/make_referral.html")
		}
	case "/check_referral", "/Check_referral", "/templates/check_referral.html":
		// if r.URL.Path == "/index" || r.URL.Path == "/Index" || r.URL.Path == "INDEX"  {
		// 	http.Redirect(w, r, "/index", http.StatusMovedPermanently)
		// }
		if r.Method == http.MethodGet {
			serveTemplate(w, "templates/check_referral.html")
		}
	case "/create_reminder", "create_reminder", "/templates/create_reminder.html":
		// if r.URL.Path == "/index" || r.URL.Path == "/Index" || r.URL.Path == "INDEX"  {
		// 	http.Redirect(w, r, "/index", http.StatusMovedPermanently)
		// }
		if r.Method == http.MethodGet {
			serveTemplate(w, "templates/create_reminder.html")
		}
	}
	
}

	

// serveTemplate loads and executes a template file
func serveTemplate(w http.ResponseWriter, filename string) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		http.Error(w, "404 Page not Found", http.StatusNotFound)
		return
	}
	tmpl := template.Must(template.ParseFiles(filename))
	errr := tmpl.Execute(w, nil)
	if errr != nil {
		log.Println("500 Internal Server Error")
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Function to check for illegal characters in input
func containsIllegalCharacters(input string) bool {
	// Regular expression to match non-printable ASCII characters
	illegalCharRegex := regexp.MustCompile(`[^\x00-\x7F]`)
	return illegalCharRegex.MatchString(input)
}
