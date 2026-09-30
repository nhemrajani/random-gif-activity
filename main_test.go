package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
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
