package main

// go install ./... && systemctl restart web-server && journalctl -f -u web-server

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strings"
)

const (
	gbbrCertFile     = "/etc/letsencrypt/live/gbbr.io/fullchain.pem"
	gbbrKeyFile      = "/etc/letsencrypt/live/gbbr.io/privkey.pem"
	birdnestCertFile = "/etc/letsencrypt/live/birdnest.io/fullchain.pem"
	birdnestKeyFile  = "/etc/letsencrypt/live/birdnest.io/privkey.pem"
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

// repoFromRequest returns the first token in the path, if it matches against
// validRepos. Otherwise it returns an empty string.
func repoFromRequest(r *http.Request) string {
	if r.Host != "gbbr.io" {
		return ""
	}
	parts := strings.Split(path.Clean("/"+r.URL.Path), "/")
	if len(parts) < 2 || parts[1] == "" {
		return ""
	}
	if _, ok := validRepos[parts[1]]; ok {
		return parts[1]
	}
	return ""
}

// notLookingForRepo will return true if the user has been redirected as a result
// of accessing a URL which doesn't point to a valid repository.
func notLookingForRepo(w http.ResponseWriter, r *http.Request) bool {
	repo := repoFromRequest(r)
	if repo == "" {
		http.Redirect(w, r, "https://github.com/gbbr", http.StatusSeeOther)
	}
	return repo == ""
}

// redirect redirects the user from HTTP to HTTPS
func redirect(w http.ResponseWriter, r *http.Request) {
	switch r.Host {
	case "birdnest.io":
		http.Redirect(w, r, "https://birdnest.io"+r.RequestURI, http.StatusSeeOther)
	default:
		didRedirect := notLookingForRepo(w, r)
		if !didRedirect {
			http.Redirect(w, r, "https://gbbr.io"+r.RequestURI, http.StatusSeeOther)
		}
	}
}

func serveRepo(w http.ResponseWriter, r *http.Request) {
	if notLookingForRepo(w, r) {
		return
	}
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	writeHTML(w, repoFromRequest(r))
}

type mux struct {
	gbbr, birdnest http.ServeMux
}

func (m *mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Host {
	case "birdnest.io":
		m.birdnest.ServeHTTP(w, r)
	default:
		m.gbbr.ServeHTTP(w, r)
	}
}

func newMux() http.Handler {
	var bird, gbbr http.ServeMux

	gbbr.HandleFunc("/", serveRepo)

	bird.Handle("/", http.FileServer(http.Dir("/home/www/web-root")))

	return &mux{gbbr, bird}
}

func main() {
	gbbrCert, err := tls.LoadX509KeyPair(gbbrCertFile, gbbrKeyFile)
	if err != nil {
		log.Fatal(err)
	}
	birdnestCert, err := tls.LoadX509KeyPair(birdnestCertFile, birdnestKeyFile)
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{Certificates: []tls.Certificate{gbbrCert, birdnestCert}}
	tlsConfig.BuildNameToCertificate()

	ln, err := tls.Listen("tcp", ":443", tlsConfig)
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{TLSConfig: tlsConfig, Handler: newMux()}

	go http.ListenAndServe(":80", http.HandlerFunc(redirect))
	if err := server.Serve(ln); err != nil {
		log.Fatal(err)
	}
}
