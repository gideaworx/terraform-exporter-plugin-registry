package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

//go:embed registry-site
var site embed.FS

func main() {
	var port int

	flag.IntVar(&port, "port", 3000, "The port to listen on")
	flag.Parse()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), http.FileServer(getSiteFS())))
}

func getSiteFS() http.FileSystem {
	subDir, err := fs.Sub(site, "registry-site")
	if err != nil {
		log.Fatal(err)
	}

	return http.FS(subDir)
}
