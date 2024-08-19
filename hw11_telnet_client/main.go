package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Set timeout flag.
	timeout := flag.Duration("timeout", 10*time.Second, "Connection timeout")
	flag.Parse()

	// Check arguments count.
	if flag.NArg() != 2 {
		log.Fatal("not enough arguments")
	}

	// Get target address from args.
	address := net.JoinHostPort(flag.Arg(0), flag.Arg(1))

	// Create new client.
	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	// Connect client.
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())

	// Sending goroutine.
	go func() {
		defer cancel()

		err := client.Send()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	// Receiving goroutine.
	go func() {
		defer cancel()

		err := client.Receive()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	// Create OS signals channel.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Wait for signals.
	select {
	case <-signals:
	case <-ctx.Done():
		close(signals)
	}
}
