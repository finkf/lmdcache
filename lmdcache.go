package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/finkf/lmd/api"
	"github.com/finkf/lmd/api/client"
	"github.com/finkf/qparams"
)

var (
	char3gramsCache = newCache(1024)
	ngramsCache     = newCache(512)
	host            string
	lmd             string
)

func init() {
	flag.StringVar(&host, "host", "localhost:8181", "set listen address")
	flag.StringVar(&lmd, "lmd", "http://localhost:8080", "set address of lmd")
}

func main() {
	flag.Parse()
	http.HandleFunc("/char3grams", handleChar3Grams)
	http.HandleFunc("/ngrams", handleNGrams)
	http.HandleFunc("/", proxy)
	log.Printf("starting server on %s [lmd: %s]", host, lmd)
	http.ListenAndServe(host, nil)
}

func proxy(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s%s", lmd, r.URL)
	log.Printf("proxy: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(w, resp.Body)
}

func handleNGrams(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling %s", r.URL)
	var q api.TrigramRequest
	if err := qparams.Decode(r.URL.Query(), &q); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	v, err := ngramsCache.get(q, func() (interface{}, error) {
		return client.New(client.WithHost(lmd)).Trigram(q.F, q.S, q.T)
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(buf.Bytes())
	log.Printf("handled %s", r.URL)
}

func handleChar3Grams(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling %s", r.URL)
	var q api.CharTrigramRequest
	if err := qparams.Decode(r.URL.Query(), &q); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	v, err := char3gramsCache.get(q, func() (interface{}, error) {
		return client.New(client.WithHost(lmd)).CharTrigram(q.Q, q.Regex)
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(buf.Bytes())
	log.Printf("handled %s", r.URL)
}
