package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/mail"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gebi/scryptauth"
	"golang.org/x/crypto/scrypt"
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

	_startTime = time.Now() //Last restart
	secretKey  string
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
	secretKey = randBase64String(128)
	db = connectToDatabase()
	go commandLineInterface()
	fmt.Println("Server is running!")
	fmt.Println("Listening on PORT: " + port)

	//Setup client interface
	http.HandleFunc("/api/login", login)
	http.HandleFunc("/api/logout", logout)
	http.HandleFunc("/api/register", register)
	http.HandleFunc("/api/profile", getProfile)
	http.Handle("/", http.FileServer(http.Dir("www")))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println(err)
	}
}

//Checks the provided credentials and authenticates
//or denies the user.
func login(w http.ResponseWriter, r *http.Request) {
	var err error

	user := getClientInfo(w, r)
	passString := user.Password

	//See if user is in database (if we're using one)
	if db != nil {
		user, err = db.LookupUser(user)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("User not found"))
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No database present"))
		return
	}

	//Authenticate user and provide a jwt
	allowed := authenticatePassword(user, passString)
	if allowed {
		writeToken(w, r, user)
	} else {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Incorrect email or password"))
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement
}

//Encrypts the users password and registers it in the database
func register(w http.ResponseWriter, r *http.Request) {
	user := getClientInfo(w, r)
	user.UserId = randBase64String(64)
	user.Salt = randBase64String(128)

	//Since the code will be run by a raspberry pi, 65536 is the best
	//we can do when it comes to cost for our key. Should be increased
	//to 1048576 (1 << 20) when migrating to a more high end system.
	passwordHash, _ := scrypt.Key([]byte(user.Password), []byte(user.Salt), (1 << 16), 8, 1, 128)
	passwordHash64 := scryptauth.EncodeBase64((1 << 16), []byte(passwordHash), []byte(user.Salt))
	user.Password = string(passwordHash64)

	if db != nil {
		err := db.AddUser(user)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("User already registered"))
		} else {
			writeToken(w, r, user)
		}
	}
}

//Generates a scrypt key from the provided password and user salt and compares them
//Returns true/false if hashes match
func authenticatePassword(user *User, password string) bool {
	passwordHash, _ := scrypt.Key([]byte(password), []byte(user.Salt), (1 << 16), 8, 1, 128)
	passwordHash64 := scryptauth.EncodeBase64((1 << 16), []byte(passwordHash), []byte(user.Salt))
	return (string(passwordHash64) == user.Password)

}

func getProfile(w http.ResponseWriter, r *http.Request) {
	//validToken := validateToken()
}

//Uses the jwt-library and the secretKey to generate a signed jwt
func generateToken(user *User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["iss"] = user.UserId
	token.Claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return tokenString, nil
}

//Reads n crypto random bytes and return them as a base64 encoded string
func randBase64String(n int) string {
	bytes := make([]byte, n)
	rand.Read(bytes)

	return base64.URLEncoding.EncodeToString(bytes)
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
	err = validateEmail(user)
	if err != nil {
		user.Email = ""
	}
	return user
}

//Validates an email provided by the user
func validateEmail(user *User) error {
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return err
	} else if strings.Contains(user.Email, "'") {
		return errors.New("Email cannot contain '")
	}
	return nil
}

//Validates a token
func validateToken(rToken string) bool {
	token, err := jwt.Parse(rToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	fmt.Println("")
	if err == nil && token.Valid {
		fmt.Println(token)
		return true
	} else {
		fmt.Println(err)
		return false
	}
}

//Generates a token and writes it to the client
func writeToken(w http.ResponseWriter, r *http.Request, user *User) {
	token, err := generateToken(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to provide web token"))
		return
	}

	JSON, err := json.Marshal(Response{token})
	if err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(JSON)
}

//Tries to open a connection to the database
//On success: global value db is assigned a DatabaseInterface-struct
//On failure: global value db is assigned nil
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
