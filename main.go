package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"golang.org/x/crypto/scrypt"
	//"github.com/dgrijalva/jwt-go"
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

	//_useDb     = true       //Flag to see if we are able to use a database
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

	//Setup client interface
	http.HandleFunc("/api/login", login)
	http.HandleFunc("/api/logout", logout)
	http.HandleFunc("/api/register", register)
	http.Handle("/", http.FileServer(http.Dir("www")))
	http.ListenAndServe(":"+port, nil)
}

//Checks the provided credentials and authenticates
//or denies the user.
func login(w http.ResponseWriter, r *http.Request) {
	var err error

	user := getClientInfo(w, r)

	//See if user is in database if we're using one
	if db != nil {
		user, err = db.LookupUser(user)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No database present"))
		return
	}

	allowed := authenticate(user)
	if allowed {
		w.Write([]byte("{\"token\":\"Token\"}"))
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	//Validate posted info in database
	//Clear session
}

//Entrypts the users password and registers it in the database
func register(w http.ResponseWriter, r *http.Request) {
	user := getClientInfo(w, r)
	user.Salt = string(generateSalt())
	passwordHash, _ := scrypt.Key([]byte(user.Password), []byte(user.Salt), (2 << 16), 8, 1, 32)
	user.Password = string(passwordHash)

	if db != nil {
		db.AddUser(user)
	}
	login(w, r)
}

func authenticate(user *User) bool {
	return true
}

//Returns a random byte slice of at least 100b in since
func generateSalt() []byte {
	salt := make([]byte, 128)
	n, _ := rand.Read(salt)
	for n >= 100 {
		salt = make([]byte, 128)
		n, _ = rand.Read(salt)
	}
	return salt
}

//Takes the request from the client and parses the json
//inside into a user struct which is being returned
func getClientInfo(w http.ResponseWriter, r *http.Request) *User {
	//Read json from client
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.FileServer(http.Dir("www"))
		return nil
	}
	user := new(User)

	//Parse json into User struct
	err = json.Unmarshal(body, user)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return nil
	}
	return user
}

//Tries to open a connection to the database
//On succes: global value db is asigned a DatabaseInterface-struct
//On failure: global value db is asigned nil
func connectToDatabase() *DatabaseInterface {
	conf, err := os.Open(".db_cnf")
	if err != nil {
		db = nil
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
		db = nil
		return nil
	}
	fmt.Println("Successfully connected to database")
	return db
}

//Takes care of closing operations
func closeServer() {
	fmt.Println("Bye!")
	if db != nil {
		db.CloseConnection()
	}
	os.Exit(0)
}
