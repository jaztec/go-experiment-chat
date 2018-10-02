package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	appErr "github.com/jaztec/go-experiment-chat/error"
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
	c, err := s.Listen()
	if appErr.HasError(err) {
		return s, nil, err
	}
	// Listen to messages
	go func() {
		for {
			select {
			case msg := <-c:
				fmt.Println(msg)
			}
		}
	}()
	return s, c, nil
}

func initializeClient() (network.ClientInterface, chan network.Message, error) {
	c, err := network.NewClient()
	if appErr.HasError(err) {
		return c, nil, err
	}
	conn, err := c.Dial()
	if appErr.HasError(err) {
		return c, conn, err
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
		if appErr.HasError(err) {
			fmt.Printf("\nServer not able to initialize: %v\n", err)
			os.Exit(2)
		}
	} else if *mode == "client" {
		c, _, err := initializeClient()
		defer c.Close()
		if appErr.HasError(err) {
			fmt.Printf("\nClient not able to initialize: %v\n", err)
			os.Exit(2)
		}
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
}
