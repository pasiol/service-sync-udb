package main

import (
	"log"
	"time"

	"github.com/pasiol/service-sync-udb/cmd"
)

var (
	// Version for build
	Version string
	// Build for build
	Build string
)

func start() {
	start := time.Now()
	log.Printf("Starttime: %v", start.Format(time.RFC3339))
	log.Println("Version: ", Version)
	log.Println("Build Time: ", Build)
}

func main() {
	start()
	cmd.Execute()
}
