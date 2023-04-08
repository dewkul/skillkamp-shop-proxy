package main

import (
	"flag"
	"fmt"

	"github.com/dewkul/skillkamp-shop-proxy/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	listenAddr := flag.String("listen", ":3030", "Listen address")
	serverUrl := flag.String("server", "https://skillkamp-api.com", "Upstream server URL")
	flag.Parse()

	server := api.NewServer(*listenAddr, *serverUrl)
	fmt.Println("server is listening on port ", *listenAddr)
	log.Fatal().Err(server.Start())
}
