package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// Get executable name.
	executable := cmd[0]

	// Get args.
	args := cmd[1:]

	// Create a command struct.
	command := exec.Command(executable, args...)
	command.Env = os.Environ()
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout

	// Set environment vars for command.
	for envKey, envVar := range env {
		if envVar.NeedRemove {
			command.Env = append(command.Env, fmt.Sprintf("%s=", envKey))
		} else {
			command.Env = append(command.Env, fmt.Sprintf("%s=%s", envKey, envVar.Value))
		}
	}

	// Run command.
	err := command.Run()
	// Get exit code.
	if err != nil {
		exitError := &exec.ExitError{}
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
	}

	return 0
}
