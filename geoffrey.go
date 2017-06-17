package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

func init() {
	// Output to stderr
	log.SetOutput(os.Stderr)

	// Set the log level to debug
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Infoln("Geoffrey is now running...")
}
