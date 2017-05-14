package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// image to push
var image []byte

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// preparing image
func init() {
	fmt.Println("Initializing logger")
	Log(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	Info.Println("Reading image")
	var err error
	image, err = ioutil.ReadFile("./image.png")
	if err != nil {
		panic(err)
	}
}

// Send HTML and push image
func handlerHtml(w http.ResponseWriter, r *http.Request) {
	pusher, ok := w.(http.Pusher)
	if ok {
		Info.Println("Pushing image")
		pusher.Push("/image", nil)
	}
	Info.Println("Writing HTML")
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, `<html><body><img src="/image"></body></html>`)
}

// Send image as usual HTTP request
func handlerImage(w http.ResponseWriter, r *http.Request) {
	Info.Println("Writing Image")
	w.Header().Set("Content-Type", "image/png")
	w.Write(image)
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	if "" == w.Header().Get("Content-Type") {
		Info.Println("Detecting ContentType")
		// If no content type, apply sniffing algorithm to un-gzipped body.
		w.Header().Set("Content-Type", http.DetectContentType(b))
	}
	Info.Println("Writing gzip")
	return w.Writer.Write(b)
}

func makeGzipHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			Info.Println("Payload is already gzip")
			fn(w, r)
			return
		}
		Info.Println("Creating gzip writer")
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r)
	}
}

func Server() {
	http.HandleFunc("/", makeGzipHandler(handlerHtml))
	http.HandleFunc("/image", makeGzipHandler(handlerImage))
	Info.Println("start http listening :1337")
	err := http.ListenAndServeTLS(":1337", "server.pem", "server.key", nil)
	Error.Println(err)
}
