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

const gifBase = "https://s3.amazonaws.com/files.656.mba/mgt656/fall-2026/random-gifs/adorbs/"

// gifsByWord maps a search word to the GIFs that match it.
var gifsByWord = map[string][]string{
	"cat": {
		gifBase + "CatPat.gif",
		gifBase + "CatBowl.gif",
		gifBase + "FightsAndHunts/LaserCat.gif",
		gifBase + "TrainedCats.gif",
		gifBase + "cat-red-dot.gif",
		gifBase + "kitten-eep.gif",
		gifBase + "kitten_toes.gif",
		gifBase + "FightsAndHunts/water-balloon-kitten.gif",
	},
	"puppy": {
		gifBase + "PuppyTongue.gif",
		gifBase + "Puppyruettes.gif",
		gifBase + "Oops/PupperJumpMiss.gif",
		gifBase + "WagThePup.gif",
		gifBase + "puppies.gif",
		gifBase + "Oops/PugsCantJump.gif",
		gifBase + "puglove.gif",
		gifBase + "derpy-dog.gif",
	},
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
  <form method="get" action="/">
    <input type="search" name="q" value="{{.Query}}" placeholder="cat or puppy">
    <button type="submit">Search</button>
  </form>
  {{if .Message}}<p>{{.Message}}</p>{{end}}
  <img src="{{.GIF}}" alt="{{.Alt}}">
</body>
</html>`))

type pageData struct {
	Query   string
	Message string
	GIF     string
	Alt     string
}

// choose picks a GIF for the search query q.
func choose(q string) pageData {
	word := strings.ToLower(strings.TrimSpace(q))
	if gifs, ok := gifsByWord[word]; ok {
		return pageData{Query: q, GIF: randomGIF(gifs), Alt: "A cute " + word}
	}
	data := pageData{Query: q, GIF: randomGIF(allGIFs), Alt: "A cute animal"}
	if word != "" {
		data.Message = "No GIFs for \"" + q + "\". Try cat or puppy! Here's a random one instead."
	}
	return data
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := page.Execute(w, choose(r.URL.Query().Get("q"))); err != nil {
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
