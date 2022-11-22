package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"net/http/cgi"
	"strings"

	"github.com/fatih/color"
)

func main() {

	var (
		addr     string
		certFile string
		keyFile  string
	)
	flag.StringVar(&addr, "address", "localhost:7070", "")
	flag.StringVar(&certFile, "cert-file", "cert.pem", "")
	flag.StringVar(&keyFile, "key-file", "key.pem", "")
	flag.Parse()

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}

	s := &http.Server{
		Addr:    addr,
		Handler: nil,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	http.HandleFunc("/cgi-bin/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.SplitN(r.URL.Path, ".cgi", 2)
		cgiPath := parts[0][1:] + ".cgi"
		color.HiCyan("%s %s %#v (CGI script: %s)", r.Method, r.URL.Path, r.URL.Query(), cgiPath)
		(&cgi.Handler{
			Path: cgiPath,
			Dir:  "./",
		}).ServeHTTP(w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		color.HiCyan("%s %s %#v", r.Method, r.URL.Path, r.URL.Query())
		http.FileServer(http.Dir("./")).ServeHTTP(w, r)
	})

	color.HiGreen("I'm ready 🙂")

	err = s.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal(err)
	}

}
