package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// image to push
var image []byte
var port int
var serverPem string
var serverKey string

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// preparing image
func init() {
	fmt.Println("Initializing logger")
	Log(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	Info.Println("parsing cli args")
	var imagename string
	flag.StringVar(&imagename, "image", "./image.png", "The path to the image")
	flag.StringVar(&serverPem, "pem", "./server.pem", "The path to the pem")
	flag.StringVar(&serverKey, "key", "./server.key", "The path to the key")
	flag.IntVar(&port, "port", 1337, "Port to bind to")
	flag.Parse()
	Info.Println("Reading image")
	var err error
	image, err = ioutil.ReadFile(imagename)
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

// Server function for starting http server
func Server() {
	http.HandleFunc("/", makeGzipHandler(handlerHtml))
	http.HandleFunc("/image", makeGzipHandler(handlerImage))
	Info.Println("start http listening " + strconv.Itoa(port))
	err := http.ListenAndServeTLS(":"+strconv.Itoa(port), serverPem, serverKey, nil)
	Error.Println(err)
}
