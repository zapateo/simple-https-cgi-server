package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/cgi"
	"os"
	"strings"

	"github.com/fatih/color"
)

const version = "v0.1.0"

func main() {

	var (
		addr         string
		certFile     string
		keyFile      string
		printVersion bool
	)
	flag.StringVar(&addr, "address", "localhost:7070", "")
	flag.StringVar(&certFile, "cert-file", "cert.pem", "")
	flag.StringVar(&keyFile, "key-file", "key.pem", "")
	flag.BoolVar(&printVersion, "version", false, "")
	flag.Parse()

	if printVersion {
		fmt.Printf("%s %s\n", os.Args[0], version)
		os.Exit(0)
	}

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

	color.HiGreen("Listening on %s", addr)
	color.HiGreen("I'm ready 🙂")

	err = s.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal(err)
	}

}
