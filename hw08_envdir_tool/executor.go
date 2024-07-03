package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// Create a command struct
	command := exec.Command(cmd[0], cmd[1:]...)

	// Set environment vars for command
	for key, envVar := range env {
		command.Env = append(os.Environ(), fmt.Printf("%s=%s", key, envVar.Value))
	}

	return 0
}
