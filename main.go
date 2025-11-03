package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/HikaruEgashira/osquery-extensions-node-packages/pkg/scanner"
	"github.com/osquery/osquery-go"
	"github.com/osquery/osquery-go/plugin/table"
)

func main() {
	var (
		socket   = flag.String("socket", "", "Path to osquery socket")
		timeout  = flag.Int("timeout", 3, "Timeout in seconds")
		interval = flag.Int("interval", 3, "Interval in seconds")
		verbose  = flag.Bool("verbose", false, "Enable verbose logging")
	)
	flag.Parse()

	// Debug: log all possible socket sources (only if verbose)
	if *verbose {
		log.Printf("Debug: Args=%v, OSQUERY_SOCKET=%s, socket flag=%s",
			flag.Args(), os.Getenv("OSQUERY_SOCKET"), *socket)
	}

	if *socket == "" {
		// Try to get socket from OSQUERY_SOCKET environment variable
		if envSocket := os.Getenv("OSQUERY_SOCKET"); envSocket != "" {
			*socket = envSocket
			if *verbose {
				log.Printf("Using socket from OSQUERY_SOCKET: %s", *socket)
			}
		}
	}

	// When launched by osqueryi --extension, socket is passed as first positional argument
	if *socket == "" && len(flag.Args()) > 0 {
		*socket = flag.Args()[0]
		if *verbose {
			log.Printf("Using socket from positional arg: %s", *socket)
		}
	}

	if *socket == "" {
		log.Fatalln("Usage: node_packages_extension --socket <path>")
	}

	if *verbose {
		log.Printf("Using socket: %s", *socket)
	}

	server, err := osquery.NewExtensionManagerServer(
		"node_packages_extension",
		*socket,
		osquery.ServerTimeout(time.Duration(*timeout)*time.Second),
		osquery.ServerPingInterval(time.Duration(*interval)*time.Second),
	)
	if err != nil {
		log.Fatalf("Error creating extension: %v", err)
	}

	columns := []table.ColumnDefinition{
		table.TextColumn("name"),
		table.TextColumn("version"),
		table.TextColumn("manager"),
		table.TextColumn("cache_path"),
	}

	server.RegisterPlugin(table.NewPlugin("node_packages", columns, func(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
		packages, err := scanner.ScanAllManagers()
		if err != nil {
			log.Printf("Error scanning packages: %v", err)
			return []map[string]string{}, nil
		}

		var results []map[string]string
		for _, pkg := range packages {
			results = append(results, map[string]string{
				"name":       pkg.Name,
				"version":    pkg.Version,
				"manager":    pkg.Manager,
				"cache_path": pkg.CachePath,
			})
		}
		return results, nil
	}))

	log.Printf("Starting node_packages extension...")
	if err := server.Run(); err != nil {
		log.Fatalf("Error running extension: %v", err)
	}
}
