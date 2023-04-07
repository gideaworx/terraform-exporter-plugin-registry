package commands

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"sync"
)

type ServeCommand struct {
	Port int `short:"p" default:"3000" env:"PORT" help:"The port to listen on"`
}

func (s *ServeCommand) Run(site embed.FS) error {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		http.ListenAndServe(fmt.Sprintf(":%d", s.Port), http.FileServer(getSiteFS(site)))
	}()

	log.Printf("listening on port %d...", s.Port)
	wg.Wait()
	return nil
}

func getSiteFS(site embed.FS) http.FileSystem {
	subDir, err := fs.Sub(site, "build/web")
	if err != nil {
		log.Fatal(err)
	}

	return http.FS(subDir)
}
