package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dewkul/skillkamp-shop-proxy/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	listenAddr := flag.String("listen", ":3030", "Listen address")
	flag.Parse()

	logLevel := os.Getenv("LOG_LEVEL")
	version := os.Getenv("VERSION")
	origins := os.Getenv("ALLOW_ORIGINS")
	serverUrl := os.Getenv("SERVER")
	if origins == "" {
		origins = "http://localhost:5173"
	}

	if serverUrl == "" {
		serverUrl = "http://localhost:3000"
	}

	switch level := strings.ToUpper(logLevel); level {
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}

	server := api.NewServer(*listenAddr, serverUrl, version, origins)
	fmt.Println("server is listening on port ", *listenAddr)
	log.Fatal().Err(server.Start())
}
