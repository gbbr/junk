package main

// go install ./... && systemctl restart web-server && journalctl -f -u web-server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strings"
)

const (
	certFile = "/etc/letsencrypt/live/gbbr.io/fullchain.pem"
	keyFile  = "/etc/letsencrypt/live/gbbr.io/privkey.pem"
)

// validRepos holds the valid repos that may be accessed via a path,
// ie. gbbr.io/repo
var validRepos = map[string]struct{}{
	"flippo":  struct{}{},
	"retreat": struct{}{},
	"hue":     struct{}{},
	"ev":      struct{}{},
}

func writeHTML(w io.Writer, repo string) {
	w.Write([]byte(`<!DOCTYPE html><html><head><meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>`))
	w.Write([]byte(fmt.Sprintf(`<meta name="go-import" content="gbbr.io/%s git https://github.com/gbbr/%s">`, repo, repo)))
	w.Write([]byte(fmt.Sprintf(`<meta http-equiv="refresh" content="0; url=https://godoc.org/gbbr.io/%s">`, repo)))
	w.Write([]byte("</head><body>"))
	w.Write([]byte(fmt.Sprintf(`Redirecting to documentation at <a href="https://godoc.org/gbbr.io/%s">godoc.org/gbbr.io/%s</a>...`, repo, repo)))
	w.Write([]byte("</body></html>"))
}

// redirect redirects user to the correct URL.
func redirect(w http.ResponseWriter, r *http.Request) {
	if notLookingForRepo(w, r) {
		return
	}
	http.Redirect(w, r, "https://gbbr.io"+r.RequestURI, http.StatusSeeOther)
}

// repoFromRequest returns the first token in the path, if it matches against
// validRepos. Otherwise it returns an empty string.
func repoFromRequest(r *http.Request) string {
	parts := strings.Split(path.Clean("/"+r.URL.Path), "/")
	if len(parts) < 2 || parts[1] == "" {
		return ""
	}
	if _, ok := validRepos[parts[1]]; ok {
		return parts[1]
	}
	return ""
}

// notLookingForRepo will return true if the user has been redirected.
func notLookingForRepo(w http.ResponseWriter, r *http.Request) bool {
	repo := repoFromRequest(r)
	if repo == "" {
		http.Redirect(w, r, "https://github.com/gbbr", http.StatusSeeOther)
	}
	return repo == ""
}

func serveRepo(w http.ResponseWriter, r *http.Request) {
	if notLookingForRepo(w, r) {
		return
	}
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	writeHTML(w, repoFromRequest(r))
}

func main() {
	go http.ListenAndServe(":80", http.HandlerFunc(redirect))
	log.Fatal(http.ListenAndServeTLS(":443", certFile, keyFile, http.HandlerFunc(serveRepo)))
}
