package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"
)

func get(t *testing.T, target string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	handler(rec, httptest.NewRequest(http.MethodGet, target, nil))
	return rec
}

func TestHomeIsHTML(t *testing.T) {
	rec := get(t, "/")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q", ct)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "<h1>") {
		t.Error("page has no <h1> heading")
	}
	if !strings.Contains(body, "background-color") {
		t.Error("page has no background color")
	}
}

func TestUnknownPathIs404(t *testing.T) {
	if rec := get(t, "/nope"); rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rec.Code)
	}
}

var imgSrc = regexp.MustCompile(`<img src="([^"]+)"`)

// gifFromPage returns the image URL on the page, failing if there is none.
func gifFromPage(t *testing.T, body string) string {
	t.Helper()
	m := imgSrc.FindStringSubmatch(body)
	if m == nil {
		t.Fatalf("page has no <img>:\n%s", body)
	}
	return m[1]
}

func TestPageShowsGIFFromList(t *testing.T) {
	list, err := os.ReadFile("gifs.txt")
	if err != nil {
		t.Fatal(err)
	}
	src := gifFromPage(t, get(t, "/").Body.String())
	if !strings.Contains(string(list), src+"\n") {
		t.Errorf("GIF %q is not in gifs.txt", src)
	}
	if !strings.Contains(get(t, "/").Body.String(), `alt="`) {
		t.Error("image has no alt text")
	}
}

func TestReloadChangesGIF(t *testing.T) {
	if len(allGIFs) < 2 {
		t.Fatalf("only %d GIFs loaded", len(allGIFs))
	}
	seen := map[string]bool{}
	for i := 0; i < 50; i++ {
		seen[gifFromPage(t, get(t, "/").Body.String())] = true
	}
	if len(seen) < 2 {
		t.Errorf("50 reloads showed only %d distinct GIF(s)", len(seen))
	}
}

func TestParseGIFsDropsBlanksAndDuplicates(t *testing.T) {
	got := parseGIFs("a\n\nb\na\n")
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Errorf("parseGIFs = %q, want [a b]", got)
	}
}

func TestSearchWordsUseListedGIFs(t *testing.T) {
	for word, gifs := range gifsByWord {
		for _, g := range gifs {
			if !slices.Contains(allGIFs, g) {
				t.Errorf("%s GIF %q is not in gifs.txt", word, g)
			}
		}
	}
}

func TestSearchCatAndPuppy(t *testing.T) {
	for _, q := range []string{"cat", "puppy", "CAT", "%20puppy%20"} {
		word := strings.ToLower(strings.Trim(q, "%20"))
		for i := 0; i < 20; i++ {
			src := gifFromPage(t, get(t, "/?q="+q).Body.String())
			if !slices.Contains(gifsByWord[word], src) {
				t.Fatalf("q=%s showed %q, not a %s GIF", q, src, word)
			}
		}
	}
}

func TestSearchUnknownOrEmpty(t *testing.T) {
	rec := get(t, "/?q=zebra")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	gifFromPage(t, body)
	if !strings.Contains(body, "Try cat or puppy") {
		t.Error("unknown query shows no helpful message")
	}

	body = get(t, "/?q=").Body.String()
	gifFromPage(t, body)
	if strings.Contains(body, "Try cat or puppy") {
		t.Error("empty query should not show the no-results message")
	}
}

func TestSearchFormAndEscaping(t *testing.T) {
	body := get(t, "/?q=%3Cscript%3E").Body.String()
	for _, want := range []string{`method="get"`, `action="/"`, `name="q"`} {
		if !strings.Contains(body, want) {
			t.Errorf("form is missing %s", want)
		}
	}
	if strings.Contains(body, "<script>") {
		t.Error("query was not HTML-escaped")
	}
}
