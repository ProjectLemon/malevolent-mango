package main

import (
    "net/http"
    "fmt"
    "os"
    "runtime"
    "os/exec"
)

var _useDb = true
var db *DatabaseInterface //in order to close on exit

func main() {
    port := "8080"
    
    if runtime.GOOS == "windows" {
        c := exec.Command("cls")
        c.Stdout = os.Stdout
        c.Run()
    } else {
        c := exec.Command("clear")
        c.Stdout = os.Stdout
        c.Run()
    }
    
    //Setup back-end
    db = connectToDatabase()
    go commandLineInterface()
    fmt.Println("Server is running!")
    fmt.Println("Listening on PORT: "+port)

    //Setup client interface    
    http.HandleFunc("/login", login)
    http.HandleFunc("/logout", logout)
    http.HandleFunc("/register", register)
	http.ListenAndServe(":"+port, http.FileServer(http.Dir("www")))
}

func login(w http.ResponseWriter, r *http.Request)  {
    //Get posted info
    //Lookup in database
    //Log in or deny
}

func logout(w http.ResponseWriter, r *http.Request)  {
    //Validate posted info in database
    //Clear session
}

func register(w http.ResponseWriter, r *http.Request)  {
    //Validate input
    //Add to database 
    //Login
}

func connectToDatabase() *DatabaseInterface {
    conf, err := os.Open(".db_cnf")
    if err != nil {
        _useDb = false
        fmt.Println("No database config file detected")
        fmt.Println("Continuing without database")
        return nil
    }
    defer conf.Close()
    
    db := new(DatabaseInterface)
    db.SetConfigurations(conf)
    err = db.OpenConnection()
    if err != nil {
        fmt.Println("Failed to connect to database with error:")
        fmt.Println(err)
        fmt.Println("Continuing without database")
        _useDb = false
        return nil
    }
    fmt.Println("Successfully connected to database")
    return db
}

func closeServer()  {
    fmt.Println("Bye!")
    if _useDb {
        db.CloseConnection()   
    }
    os.Exit(0)
}