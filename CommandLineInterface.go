package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func commandLineInterface() {
	printWelcome()
	catchCtrlC()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("mango> ")
		text, _ := reader.ReadString('\n')
		text = strings.ToLower(text)
		text = strings.Replace(text, "\n", "", -1)
		handle(text)
	}
}

func handle(input string) {
	switch input {
	case "help":
		printCommands()
		break
	case "quit":
		closeServer()
		break
	case "exit":
		closeServer()
		break
	case "version":
		fmt.Println("Current version: ", _version)
		break
	case "uptime":
		printUpTime()
	default:
		break
	}
}

func printWelcome() {
	fmt.Println(" ")
	fmt.Println("MALICIOUS MANGO ", _version)
	fmt.Println("Welcome to Malicious-Mango webserver command line interface")
	fmt.Println("For a list of available commands, type 'help' ")
	fmt.Println(" ")
}

func printCommands() {
	fmt.Println("Available commands: ")
	fmt.Println("\t help - print this help")
	fmt.Println("\t version - show the current server version")
	fmt.Println("\t uptime - show uptime for server")
	fmt.Println("\t quit/exit - close the server")
}

func printUpTime() {
	uptime := time.Since(_startTime)
	fmt.Println(" ")
	fmt.Println("Current uptime: ", uptime)
	fmt.Println("Running since: ", _startTime)
	fmt.Println(" ")
}

func catchCtrlC() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println(" ")
		closeServer()
	}()
}
