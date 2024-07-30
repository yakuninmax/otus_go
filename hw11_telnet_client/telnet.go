package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

// Client struct.
type Client struct {
	address    string
	timeout    time.Duration
	connection net.Conn
	in         io.ReadCloser
	out        io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *Client) Connect() error {
	// Get connection.
	connection, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}
	c.connection = connection

	return nil
}

func (c *Client) Close() error {
	// Close connection, and return error.
	return c.connection.Close()
}

func (c *Client) Send() error {
	// Send.
	_, err := io.Copy(c.connection, c.in)
	if err != nil {
		return fmt.Errorf("sending error: %w", err)
	}

	return nil
}

func (c *Client) Receive() error {
	// Receive.
	_, err := io.Copy(c.out, c.connection)
	if err != nil {
		return fmt.Errorf("receiving error: %w", err)
	}

	return nil
}
