package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const page = `<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>Random GIF</title>
</head>
<body style="background-color: #fde2e4; font-family: sans-serif; text-align: center;">
  <h1>Random GIF</h1>
</body>
</html>`

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := fmt.Fprint(w, page); err != nil {
		log.Printf("write response: %v", err)
	}
}

func main() {
	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Open http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
