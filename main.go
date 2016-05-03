package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	pseudoRand "math/rand"
	"net/http"
	"net/mail"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gebi/scryptauth"
	"github.com/kennygrant/sanitize"
	"github.com/nytimes/gziphandler" //We might need some sort of license for this
	"golang.org/x/crypto/scrypt"
)

const (
	_version = 0.3
	port     = "8080"
	//Since the code will be run by a raspberry pi, 65536 is the best
	//we can do when it comes to cost for our key. Should be increased
	//to 1048576 (1 << 20) when migrating to a more high end system.
	_passwordCost = 1 << 12
)

var (
	//since opening and closing the database is considered
	//an expensive operation we keep this global to prevent
	//unneccesairy calls to the sql api
	db *DatabaseInterface

	_startTime = time.Now() //Last restart
	quit       = make(chan bool)
	secretKey  string
)

func main() {

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
	go commandLineInterface(quit)
	go SessionCleaner(quit)
	go ImageCleaner(quit)
	fmt.Println("Server is running!")
	fmt.Println("Listening on PORT: " + port)

	//Setup client interface
	http.HandleFunc("/api/login", login)
	http.HandleFunc("/api/logout", logout)
	http.HandleFunc("/api/register", register)
	http.HandleFunc("/api/refreshtoken", refreshToken)
	http.HandleFunc("/api/profile/save", saveProfile)
	http.HandleFunc("/api/profile/get-edit", getProfileEdit)
	http.HandleFunc("/api/profile/get-view/", getProfileView)

	http.HandleFunc("/api/upload/pdf", receiveUploadPDF)
	http.HandleFunc("/api/upload/profile-header", receiveUploadHeader)
	http.HandleFunc("/api/upload/profile-icon", receiveUploadIcon)

	//Setup gzip for everything
	fs := http.FileServer(http.Dir("www"))
	withoutGz := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})
	withGz := gziphandler.GzipHandler(withoutGz)
	http.Handle("/", withGz)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println(err)
	}
}

//TODO: These three functions could be merged if we got some more info from client
func receiveUploadHeader(w http.ResponseWriter, r *http.Request) {

	path, err := saveFile("img/profile-headers/", r)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("Unable to upload file"))
		return
	}
	w.Write([]byte(path))
}
func receiveUploadIcon(w http.ResponseWriter, r *http.Request) {
	path, err := saveFile("img/profile-icons/", r)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("Unable to upload file"))
		return
	}
	w.Write([]byte(path))
}

func receiveUploadPDF(w http.ResponseWriter, r *http.Request) {
	path, err := saveFile("pdf/", r)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("Unable to upload file"))
		return
	}
	w.Write([]byte(path))
}

func saveFile(folder string, r *http.Request) (string, error) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer file.Close()
	handler.Filename = sanitizeUploadFileName(handler.Filename, handler.Filename[(len(handler.Filename)-4):])
	path := folder + handler.Filename

	f, err := os.OpenFile("www/"+path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()
	io.Copy(f, file)

	return path, nil
}

func sanitizeUploadFileName(name, extension string) string {
	if db != nil {
		nameInDatabase, err := db.UniversalLookup(name)
		if nameInDatabase || err != nil {
			name = randBase64String(50)
		}
	}
	if len(name) >= 150 {
		var bytes bytes.Buffer
		name = name[:50]
		bytes.WriteString(name)
		bytes.WriteString(extension)
	}

	path := sanitize.Path(name)
	path = path + extension
	return path
}

//Checks the provided credentials and authenticates
//or denies the user.
func login(w http.ResponseWriter, r *http.Request) {
	if !usingDatabase(w) {
		return
	}

	var err error

	user, err := getClientBody(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	passString := user.Password

	user, err = db.LookupUser(user)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("User not found"))
		return
	}

	allowed := authenticatePassword(user, passString)
	if allowed {
		writeNewToken(w, r, user)
		db.InsertUserSession(user)
	} else {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Incorrect email or password"))
	}
}

//This function will evaluate the user token and if valid provide
//the client with a new one. Valid for 5 minutes.
func refreshToken(w http.ResponseWriter, r *http.Request) {
	if !usingDatabase(w) {
		return
	}

	user, err := handleToken(w, r)
	if err != nil {
		return
	}

	writeNewToken(w, r, user)
	err = db.UpdateUserSession(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		w.Write([]byte("Unable to update token value in database"))
		return
	}
}

//Removes the active session for the user in database which will
//make the rest of the code treat the user as not logged in
func logout(w http.ResponseWriter, r *http.Request) {
	if !usingDatabase(w) {
		return
	}

	user, err := handleToken(w, r)
	if err != nil {
		return
	}
	user, err = db.GetUserSession(user)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("Not logged in"))
		return
	}
	err = db.RemoveUserSession(user.Session)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("Not currently logged in"))
		return
	}

}

//Encrypts the users password and registers it in the database
func register(w http.ResponseWriter, r *http.Request) {
	if !usingDatabase(w) {
		return
	}

	user, err := getClientBody(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	//64 and 128 is a result of Database limitations and security recomendations
	user.UserID = randBase64String(64) //TODO: This should include a unique check
	user.Salt = randBase64String(128)

	passwordHash, _ := scrypt.Key([]byte(user.Password), []byte(user.Salt), _passwordCost, 8, 1, 128)
	passwordHash64 := scryptauth.EncodeBase64(_passwordCost, []byte(passwordHash), []byte(user.Salt))
	user.Password = string(passwordHash64)

	err = db.AddUser(user)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("User already registered"))
		return
	}
	writeNewToken(w, r, user)
	db.InsertUserSession(user)
}

//Generates a KDF from the provided password and user salt and compares them
//Returns true/false if hashes match
func authenticatePassword(user *User, password string) bool {
	passwordHash, _ := scrypt.Key([]byte(password), []byte(user.Salt), _passwordCost, 8, 1, 128)
	passwordHash64 := scryptauth.EncodeBase64(_passwordCost, []byte(passwordHash), []byte(user.Salt))
	return (string(passwordHash64) == user.Password)

}

//Returns a profile to the client
func getProfileView(w http.ResponseWriter, r *http.Request) {
	if !usingDatabase(w) {
		return
	}

	requestURLParts := strings.Split(r.RequestURI, "/")
	if len(requestURLParts) < 2 {
		return
	}

	publicName := requestURLParts[len(requestURLParts)-1]
	uid, err := db.GetUserIDFromPublicName(publicName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	user := new(User)
	user.UserID = uid
	user, err = db.LookupUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	writeUserContentToClient(w, r, user)
}

//Validates token and returns a profile to client for edit
func getProfileEdit(w http.ResponseWriter, r *http.Request) {
	if !usingDatabase(w) {
		return
	}

	user, err := handleToken(w, r) //feedback to client happens inside function
	if err != nil {
		return
	}
	writeUserContentToClient(w, r, user)
}

//Writes UserContent from database to client
func writeUserContentToClient(w http.ResponseWriter, r *http.Request, user *User) {
	if !usingDatabase(w) {
		return
	}

	userContent := new(UserContents)
	userContent, err := db.GetUserContents(user.UserID, userContent)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("No content for the specified user"))
		return
	}
	userContent.UserID = "" //Since it potentially could be exploited if we sent uid to client
	JSON, err := json.Marshal(userContent)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to send content"))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(JSON)
}

//Saves the profile into database
func saveProfile(w http.ResponseWriter, r *http.Request) {
	if !usingDatabase(w) {
		return
	}

	_, err := handleToken(w, r) //feedback to client happens inside function
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not read content"))
		return
	}
	userContent := new(UserContents)

	err = json.Unmarshal(body, userContent)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unexpected end of json-input"))
		return
	}
	user := new(User)
	user.Token = strings.Split(r.Header.Get("Authorization"), " ")[1]
	_, err = db.GetUserSession(user)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("No active session"))
		return
	}
	err = db.UpdateUserContent(user.UserID, userContent)
	if err != nil {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		fmt.Println(err)
		w.Write([]byte(err.Error()))
		return
	}

	//We are using the users public name as part of their URL
	//therefore we have to make sure it's unique or else make it unique
	publicName := strings.ToLower(userContent.FullName)
	publicName = strings.Replace(publicName, " ", "", -1)
	nameInDb, _ := db.LookupPublicName(publicName)
	if nameInDb {
		//Since we want a 4 digit long number we have to do this
		//somewhat complicated conversion from []int to string using a []byte
		pseudoRand.Seed(time.Now().UTC().UnixNano())
		nrExtension := pseudoRand.Perm(4)
		temp := []byte{}
		for i := 0; i < len(nrExtension); i++ {
			temp = strconv.AppendInt(temp, int64(nrExtension[i]), 10)
		}
		strNrExtension := string(temp)

		publicName += strNrExtension
		userContent.PublicName = publicName
		db.UpdatePublicName(userContent, user)

		//cat and insert to database
	} else {
		userContent.PublicName = publicName
		db.UpdatePublicName(userContent, user)
	}
} // End saveProfile

//Uses the jwt-library and the secretKey to generate a signed jwt
func generateToken(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["uid"] = userID
	token.Claims["exp"] = time.Now().Add(time.Minute * 2).Unix()
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
func getClientBody(w http.ResponseWriter, r *http.Request) (*User, error) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}
	user := new(User)

	err = json.Unmarshal(body, user)
	if err != nil {
		return nil, err
	}
	err = validateEmail(user.Email)
	if err != nil {
		user.Email = ""
	}
	return user, nil
}

//Validates an email provided by the user
func validateEmail(email string) error {
	if len(email) >= 80 {
		return errors.New("Email can be no more than 80 chars in length")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return err
	} else if strings.Contains(email, "'") {
		return errors.New("Email cannot contain '")
	}
	return nil
}

//handleToken takes care of reading and validating the token provided by the client
//returns a user containing token and user session
func handleToken(w http.ResponseWriter, r *http.Request) (*User, error) {
	user := new(User)
	providedTokens := strings.Split(r.Header.Get("Authorization"), " ")
	if len(providedTokens) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invlaid number of tokens provided"))
		return nil, errors.New("Invlaid number of tokens provided")
	}
	user.Token = providedTokens[1]

	user, err := db.GetUserSession(user)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return nil, err
	}

	valid, _ := validateToken(user)
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid web token"))
		return nil, errors.New("Invalid web token")
	}
	return user, nil
}

//Validates a token's signing method, userID and expiration date
func validateToken(user *User) (bool, *jwt.Token) {
	token, err := jwt.Parse(user.Session.SessionKey, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	if token.Claims["uid"] != user.UserID {
		token.Valid = false
	} else if token.Claims["exp"].(float64) <= float64(time.Now().Unix()) {
		token.Valid = false
	}

	if token.Valid {
		return true, token
	}
	fmt.Println(err)
	return false, token
}

//Generates a token and writes it to the client
func writeNewToken(w http.ResponseWriter, r *http.Request, user *User) {
	token, err := generateToken(user.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to provide web token"))
		return
	}
	user.Token = token

	JSON, err := json.Marshal(Response{token})
	if err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(JSON)
}

func usingDatabase(w http.ResponseWriter) bool {
	if db == nil {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("No database associated"))
		return false
	}
	return true
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
	close(quit) //Exits all runnign go routines
	os.Exit(0)
}
