package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
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

	//_useDb     = true       //Flag to see if we are able to use a database
	_startTime = time.Now() //Last restart
	secretKey  *rsa.PrivateKey
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
	secretKey = generatePrivateRSAKey()
	db = connectToDatabase()
	go commandLineInterface()
	fmt.Println("Server is running!")
	fmt.Println("Listening on PORT: " + port)

	//Setup client interface
	http.HandleFunc("/api/login", login)
	http.HandleFunc("/api/logout", logout)
	http.HandleFunc("/api/register", register)
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
	allowed := authenticate(user, passString)
	if allowed {
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
	} else {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Incorrect email or password"))
	}
}

func logout(w http.ResponseWriter, r *http.Request) {

}

//Encrypts the users password and registers it in the database
func register(w http.ResponseWriter, r *http.Request) {
	user := getClientInfo(w, r)
	user.Salt = string(generateSalt())
	passwordHash, _ := scrypt.Key([]byte(user.Password), []byte(user.Salt), (1 << 14), 8, 1, 128)
	passwordHash64 := scryptauth.EncodeBase64((1 << 14), []byte(passwordHash), []byte(user.Salt))
	user.Password = string(passwordHash64)

	if db != nil {
		err := db.AddUser(user)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("User already registered"))
		}
	}
}

func authenticate(user *User, password string) bool {
	passwordHash, _ := scrypt.Key([]byte(password), []byte(user.Salt), (1 << 14), 8, 1, 128)
	passwordHash64 := scryptauth.EncodeBase64((1 << 14), []byte(passwordHash), []byte(user.Salt))
	fmt.Println(passwordHash64)
	fmt.Println(user.Password)
	return (string(passwordHash64) == user.Password)

}

//Uses the jwt-library and the secretKey to generate a signed jwt
func generateToken(user *User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["iss"] = user.Email
	token.Claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
	key, err := x509.MarshalPKIXPublicKey(secretKey.PublicKey)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

//Returns a random byte slice of at least 100b in size
func generateSalt() []byte {
	salt := make([]byte, 128)
	rand.Read(salt)

	return []byte(base64.URLEncoding.EncodeToString(salt))
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

//Generates a RSA key pair (used for signing web tokens)
func generatePrivateRSAKey() *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, (2 << 9))
	if err != nil {
		fmt.Println("Unable to obtain private key")
		fmt.Println("Continued usage is dicurraged")
	}
	return key
}

//Takes care of closing operations
func closeServer() {
	fmt.Println("Bye!")
	if db != nil {
		db.CloseConnection()
	}
	os.Exit(0)
}
