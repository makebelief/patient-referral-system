package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"digital-referral-system/handlers"
)

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Printf("Failed to open browser: %v", err)
	}
}

func main() {
	if len(os.Args) != 1 {
		fmt.Println("usage: go run .")
		return
	}
	
	// Define HTTP handlers
	http.HandleFunc("/", handlers.Index)
    http.HandleFunc("/index", handlers.Index)

    

//	http.HandleFunc("/ascii-art", handlers.HandleASCIIArt)

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	
	// Defining server protocol and http port
	url := "http://localhost:8080"
	log.Println("Server is running on", url)
	openBrowser(url)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
