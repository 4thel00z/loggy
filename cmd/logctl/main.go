package main

import (
	"flag"
	"fmt"
	"github.com/4thel00z/loggy"
	tea "github.com/charmbracelet/bubbletea"
	"net/http"
	"os"
)

var (
	port    = flag.Int("port", 12345, "Port to listen on")
	host    = flag.String("host", "localhost", "Host to listen on")
	tls     = flag.Bool("tls", false, "Use TLS")
	retries = flag.Int("retries", 5, "Number of retries")
)

func main() {
	flag.Parse()
	fetcher := loggy.GenerateLogEntriesFetcher(http.DefaultClient, *host, *port, *tls, *retries)
	p := tea.NewProgram(loggy.NewModel(fetcher), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
