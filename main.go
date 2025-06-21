package main

import (
	"crawler-bagu/internal"
	"github.com/rs/zerolog/log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatal().Err(err).Msg("http.ListenAndServe error")
		}
	}()
	app := internal.InitializeApp()
	app.Run()
}
