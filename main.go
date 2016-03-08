package main

import (
    "net/http"
    "fmt"
    "os"
)

var _useDb = true

func main() {
    connectToDatabase()
    http.HandleFunc("/login", login)
    http.HandleFunc("/logout", logout)
    http.HandleFunc("/register", register)
	http.ListenAndServe(":8080", http.FileServer(http.Dir("www")))
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

func connectToDatabase()  {
    conf, err := os.Open(".db_cnf")
    if err != nil {
        _useDb = false
        fmt.Println("No database config file detected")
        fmt.Println("Continuing without database")
        return
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
        return
    }
    fmt.Println("Successfully connected to database")
}
