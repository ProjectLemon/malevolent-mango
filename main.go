package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

//Global constants
const (
	_version = 0.1
)

var (
	//since opening and closing the database is considered
	//an expensive operation we keep this global to prevent
	//unneccesairy calls to the sql api
	db *DatabaseInterface

	_useDb     = true       //Flag to see if we are able to use a database
	_startTime = time.Now() //Last restart
)

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

	//New user for exemple
	user := NewUser("linus.lagerhjelm@gmail.com")
	db.LookupUser(user)

	//Setup client interface
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/register", register)
	http.Handle("/", http.FileServer(http.Dir("www")))
	http.ListenAndServe(":"+port, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	credentials := NewUser(r.PostFormValue("email"))
	user, err := db.LookupUser(credentials)
	if err != nil {
		http.FileServer(http.Dir("www"))
		return
	}
	allowed := authenticate(user)
	if allowed {
		http.FileServer(http.Dir("www")) //temporary line
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	//Validate posted info in database
	//Clear session
}

func register(w http.ResponseWriter, r *http.Request) {
	//Validate input
	//Add to database
	//Login

	//This is examplecode and should be rewritten to suit our needs
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["foo"] = "bar"
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	// Sign and get the complete encoded token as a string
	_, err := token.SignedString(nil)
	if err != nil {
		fmt.Println("Log")
	}
}

func authenticate(user *User) bool {
	return true
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
