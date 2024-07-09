package main

import (
	"log"
	"os"
)

func main() {
	// Check arguments.
	if len(os.Args) < 3 {
		log.Fatal("Not enough arguments.")
	}

	// Read env vars.
	envVars, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// Run command.
	exitCode := RunCmd(os.Args[2:], envVars)

	// Exit with code.
	os.Exit(exitCode)
}
