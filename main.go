package main

import (
    "net/http"
    "log"
)

func main() {
    // Define the directory to serve static files from
    fs := http.FileServer(http.Dir("."))

    // Handle routes
    http.Handle("/", fs) // Handle the root route to serve static files

    // Specify the port for the server
    port := ":3000"
    log.Printf("Server is running at http://localhost%s\n", port)
    
    // Start the server
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
