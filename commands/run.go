package commands

import (
	"embed"
	"fmt"
	"net/http"
)

type ServeCommand struct {
	Port int `short:"p" default:"3000" env:"PORT" help:"The port to listen on"`
}

func (s *ServeCommand) Run(site embed.FS) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), http.FileServer(http.FS(site)))
}
