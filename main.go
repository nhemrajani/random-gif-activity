package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

const gifBase = "https://s3.amazonaws.com/files.656.mba/mgt656/fall-2026/random-gifs/adorbs/"

var page = template.Must(template.New("page").Parse(`<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>Random GIF</title>
</head>
<body style="background-color: #fde2e4; font-family: sans-serif; text-align: center;">
  <h1>Random GIF</h1>
  <img src="{{.GIF}}" alt="A cute animal">
</body>
</html>`))

type pageData struct {
	GIF string
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := page.Execute(w, pageData{GIF: gifBase + "CatPat.gif"}); err != nil {
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
