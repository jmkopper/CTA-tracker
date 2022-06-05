// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	config := readConfig()

	// create the request handlers
	ctah := ctaReqHandler{
		config: &config,
	}
	mongoh := mongoSearchHandler{config: &config}
	staticFileServer := http.FileServer(http.Dir("static"))

	mux.Handle("/getCTAData", ctah)
	mux.Handle("/search", mongoh)
	mux.Handle("/", staticFileServer)

	listenAt := fmt.Sprintf(":%d", config.port)
	log.Printf("Open the following URL in the browser: http://localhost:%d\n", config.port)
	log.Fatal(http.ListenAndServe(listenAt, mux))
}
