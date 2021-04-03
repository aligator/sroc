package main

// Based on https://github.com/nicolaspearson/gogo-cors-proxy/blob/master/proxy.go

import (
	"crypto/tls"
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"strings"
)

var (
	ErrInvalidOrigin   = errors.New("invalid origin tried to use this service")
	ErrBadTarget       = errors.New("the target is wrong")
	ErrFetchingTarget  = errors.New("could not fetch the target")
	ErrSendingResponse = errors.New("could not send response")
)

type handler struct {
	origin      string
	credentials bool
	methods     string
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	originHeader := r.Header.Get("origin")

	if originHeader != h.origin {
		log.Printf("%v: %v", ErrInvalidOrigin, originHeader)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	targets, ok := r.URL.Query()["target"]
	if !ok || len(targets) != 1 {
		log.Println(ErrBadTarget)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	target := targets[0]

	log.Printf("Proxy %v -> %v", h.origin, target)

	w.Header().Add("Access-Control-Allow-Origin", originHeader)

	if h.credentials {
		w.Header().Add("Access-Control-Allow-Credentials", "true")
	}

	w.Header().Add("Access-Control-Allow-Methods", h.methods)

	if r.Method == "OPTIONS" {
		for n, h := range r.Header {
			if strings.Contains(n, "Access-Control-Request") {
				for _, h := range h {
					k := strings.Replace(n, "Request", "Allow", 1)
					w.Header().Add(k, h)
				}
			}
		}
		return
	}

	// Fetch the target resource.
	req, err := http.NewRequest(r.Method, target, r.Body)
	if !ok || len(targets) != 1 {
		log.Println(ErrFetchingTarget)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Setup headers for the request.
	for n, h := range r.Header {
		for _, h := range h {
			if n != "Access-Control-Allow-Origin" && n != "Access-Control-Allow-Credentials" && n != "Access-Control-Allow-Methods" {
				req.Header.Add(n, h)
			}
		}
	}

	client := http.Client{}
	if r.TLS != nil {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("%v: %v\n", ErrFetchingTarget, err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	// Copy the response from the server to the connected client request.
	for h, v := range resp.Header {
		for _, v := range v {
			w.Header().Add(h, v)
		}
	}
	w.WriteHeader(resp.StatusCode)

	n, err := io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("%v: already sent: %v, error: %v\n", ErrSendingResponse, n, err)
	}

	log.Printf("Sent %v <- %v", h.origin, target)
}

func main() {
	h := handler{
		origin:      "http://localhost:3000",
		credentials: false,
		methods:     "GET, PUT, POST, HEAD, TRACE, DELETE, PATCH, COPY, HEAD, LINK, OPTIONS",
	}
	listen := "localhost:4242"
	flag.StringVar(&listen, "listen", listen, "host:port to listen on")
	flag.StringVar(&h.origin, "origin", h.origin, "the only allowed origin")
	flag.BoolVar(&h.credentials, "credentials", h.credentials, "add Access-Control-Allow-Credentials header")
	flag.StringVar(&h.methods, "methods", h.methods, "restrict to only a few methods. (e.g. 'GET, POST')")
	flag.Parse()

	log.Fatal(http.ListenAndServe(listen, &h))
}
