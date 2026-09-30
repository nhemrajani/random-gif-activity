package main

import (
	_ "embed"
	"html/template"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strings"
)

//go:embed gifs.txt
var gifList string

// allGIFs holds every URL in gifs.txt, with duplicates removed.
var allGIFs = parseGIFs(gifList)

func parseGIFs(list string) []string {
	var gifs []string
	seen := map[string]bool{}
	for _, line := range strings.Split(list, "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !seen[line] {
			seen[line] = true
			gifs = append(gifs, line)
		}
	}
	return gifs
}

func randomGIF(gifs []string) string {
	return gifs[rand.IntN(len(gifs))]
}

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
	if err := page.Execute(w, pageData{GIF: randomGIF(allGIFs)}); err != nil {
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
