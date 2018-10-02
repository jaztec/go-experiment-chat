package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/jaztec/go-experiment-chat/errors"
	"github.com/jaztec/go-experiment-chat/network"
)

func printMicros(since int64) {
	fmt.Print("\nRuntime ")
	fmt.Print(time.Now().UnixNano()/int64(time.Microsecond) - since)
	fmt.Print(" micros\n")
}

func startStampMicros() int64 {
	return time.Now().UnixNano() / int64(time.Microsecond)
}

func initializeServer() (network.ServerInterface, chan network.Message, error) {
	s := network.NewServer(network.ServerConfig{Port: 12356, MessageBufferSize: 1024})
	msgs, err := s.Listen()
	if errors.HasError(err) {
		return s, nil, err
	}
	// Listen to messages
	go func() {
		for {
			select {
			case msg := <-msgs:
				s.Broadcast(msg)
			}
		}
	}()
	return s, msgs, nil
}

func initializeClient() (network.ClientInterface, chan network.Message, error) {
	c, err := network.NewClient("tcp", "localhost:12356")
	if errors.HasError(err) {
		return c, nil, err
	}
	conn, err := c.Dial()
	if errors.HasError(err) {
		c.Close()
		return nil, nil, err
	}

	return c, conn, nil
}

func usage() {
	flag.PrintDefaults()
}

func main() {
	start := startStampMicros()

	mode := flag.String("mode", "client", "Define program mode, provide 'server' or 'client' mode")
	help := flag.Bool("h", false, "Print this help message")

	flag.Parse()

	if *help {
		usage()
		os.Exit(0)
	}

	if *mode == "server" {
		s, _, err := initializeServer()
		defer s.Close()
		if errors.HasError(err) {
			fmt.Printf("\nServer not able to initialize: %v\n", err)
			os.Exit(2)
		}
	} else if *mode == "client" {
		c, _, err := initializeClient()
		defer c.Close()
		if errors.HasError(err) {
			fmt.Printf("\nClient not able to initialize: %v\n", err)
			os.Exit(2)
		}
		// Capture input
		go func() {
			reader := bufio.NewReader(os.Stdin)
			for {
				fmt.Print("Enter text: ")
				text, _ := reader.ReadString('\n')
				c.Send(c.CreateMessage([]byte(text)))
			}
		}()
	} else {
		fmt.Printf("\nNo valid mode provided: %s\n", *mode)
		os.Exit(2)
	}

	// Exit stategy
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)
	// Block until a signal is received
	<-sigc

	printMicros(start)

	os.Exit(0)
}
