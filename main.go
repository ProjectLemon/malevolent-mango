package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
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
	fmt.Println("Listening on PORT: " + port)

	//Setup client interface
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/register", register)
	http.Handle("/", http.FileServer(http.Dir("www")))
	http.ListenAndServe(":"+port, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	//Get posted info
	//Lookup in database
	//Log in or deny

	/* I'm just testing the login on client, remove this
	   when you actually implement the login /Fredrik */

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}

	type User struct {
		Email    string
		Password string
	}
	var user User
	err = json.Unmarshal(body, &user)

	//fmt.Println("Header: ", r.Header)
	//fmt.Println("Body: ", user)

	if user.Email == "testing@example.com" && user.Password == "supersecret" {
		w.Write([]byte("{\"token\":\"Token\"}"))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Wrong email or password"))
	}
	/* End of testing */
}

func logout(w http.ResponseWriter, r *http.Request) {
	//Validate posted info in database
	//Clear session
}

func register(w http.ResponseWriter, r *http.Request) {
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

func closeServer() {
	fmt.Println("Bye!")
	if _useDb {
		db.CloseConnection()
	}
	os.Exit(0)
}
