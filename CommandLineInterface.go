package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
    "syscall"
    "os/signal"
)

func commandLineInterface() {
    printWelcome()
    catchCtrlC()
    
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("mango> ")
        text, _ := reader.ReadString('\n')
        text = strings.ToLower(text)
        handle(text)
    }
}

func handle(input string)  {
    switch input {
    case "help\n":
        printCommands()
        break
    case "quit\n":
        closeServer()
        break
    case "exit\n":
        closeServer()
        break
    default:
        break
    }
}

func printWelcome() {
    fmt.Println(" ")
    fmt.Println("MALICIOUS MANGO 0.5")
    fmt.Println("Welcome to Malicious-Mango webserver command line interface")
    fmt.Println("For a list of available commands, type 'help' ")
    fmt.Println(" ")
}

func printCommands() {
    fmt.Println("Available commands: ")
    fmt.Println("\t help - print this help")
    fmt.Println("\t quit/exit - close th server")
}

func catchCtrlC()  {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    signal.Notify(c, syscall.SIGTERM)
    go func() {
        <-c
        fmt.Println(" ")
        closeServer()
    }()
}
