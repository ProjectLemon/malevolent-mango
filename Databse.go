package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
	"time"
)

var (
	//ErrNoUserFound if the specified email was not in the database
	ErrNoUserFound = errors.New("User was not found in database")

	//ErrNoActiveSession if the session was not found in database
	ErrNoActiveSession = errors.New("Session was not found for the specified user")

	//ErrNoContentInDatabase if no content was found for the user
	ErrNoContentInDatabase = errors.New("No content in database for the specified user")
)

//DatabaseInterface represent a configuration object, containing configurations
// for the current database
type DatabaseInterface struct {
	User           string
	Password       string
	DriverName     string
	DataSourceName string
	DB             *sql.DB
}

//SetConfigurations reads the specified config file
//and sets the respective fields in the DatabaseInterface
func (dbi *DatabaseInterface) SetConfigurations(f *os.File) {
	cnf := make(map[string]string)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if data := scanner.Text(); strings.Contains(data, " ") {
			insert := strings.Split(data, " ")
			cnf[strings.ToUpper(insert[0])] = insert[1]
		}
	}
	dbi.User = cnf["USER"]
	dbi.Password = cnf["PASSWORD"]
	dbi.DriverName = cnf["DRIVERNAME"]
	dbi.DataSourceName = cnf["DATASOURCENAME"]
}

//OpenConnection connects to a database using the information
//provided in DatabaseInterface
func (dbi *DatabaseInterface) OpenConnection() error {
	db, err := sql.Open(dbi.DriverName, dbi.getConnectionString())
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	dbi.DB = db
	return nil
}

//LookupUser sends a query to the database for the specified
//username and password hash. Returns error if query failed
func (dbi *DatabaseInterface) LookupUser(user *User) (*User, error) {
	rows, err := dbi.DB.Query("SELECT * FROM Users WHERE EMail = '" + user.Email + "'")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&user.Email,
			&user.UserId,
			&user.Password,
			&user.Salt)
		if err != nil {
			fmt.Println(err)
		}
	}

	if userInDatabase(user) {
		return user, nil
	}

	return nil, ErrNoUserFound
}

//AddUser inserts the specified user into the database
//returns error where err == nil if everything went okay
func (dbi *DatabaseInterface) AddUser(user *User) error {
	_, err := dbi.DB.Exec(
		"INSERT INTO Users (EMail, UserId, Password, PasswordSalt) VALUES (?,?,?,?)",
		user.Email,
		user.UserId,
		user.Password,
		user.Salt)
	return err
}

//GetUserContents looks up, and return, user content in database
func (dbi *DatabaseInterface) GetUserContents(uid string, userContent *UserContents) (*UserContents, error) {
	rows, err := dbi.DB.Query("SELECT * FROM UserContent WHERE UserId='" + uid + "'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jsonField []uint8

	for rows.Next() {
		err := rows.Scan(
			&userContent.UserId,
			&userContent.FullName,
			&userContent.Phone,
			&userContent.EMail,
			&userContent.ProfileIcon,
			&userContent.ProfileHeader,
			&userContent.Description,
			&jsonField)
		if err != nil {
			fmt.Println(err)
		}
		//Because []string is not supported by the database api in go
		userContent.PDFs = getStringArray(jsonField)
	}
	if contentInDatabase(userContent) {
		return userContent, nil
	}

	return nil, ErrNoContentInDatabase
}

//InsertUserSession creates a new row in the database for a user session
func (dbi *DatabaseInterface) InsertUserSession(user *User) error {
	_, err := dbi.DB.Exec(
		"INSERT INTO UserSession (SessionKey, UserId, LoginTime, LastSeenTime) VALUES (?,?,?,?)",
		user.Token,
		user.UserId,
		time.Now().Format(time.RFC3339),
		time.Now().Format(time.RFC3339))
	return err
}

//GetUserSession reads the user session for the specified user
//into the user session field of the struct
func (dbi *DatabaseInterface) GetUserSession(user *User) (*User, error) {
	rows, err := dbi.DB.Query("SELECT * FROM UserSession WHERE SessionKey='" + user.Token + "'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user.Session = new(UserSession)
	for rows.Next() {
		err := rows.Scan(
			&user.Session.SessionKey,
			&user.UserId,
			&user.Session.LoginTime,
			&user.Session.LastSeen)
		if err != nil {
			fmt.Println(err)
		}
	}

	if user.UserId != "" {
		return user, nil
	}

	return nil, ErrNoActiveSession
}

//CloseConnection closes any active connection to the current database
func (dbi *DatabaseInterface) CloseConnection() {
	dbi.DB.Close()
}

//getConnectionString returns the connection details as a formatted dataSourceName
func (dbi *DatabaseInterface) getConnectionString() string {
	return dbi.User + ":" + dbi.Password + "@" + dbi.DataSourceName
}

//inDatabase checks if the current user was found in the database
func userInDatabase(u *User) bool {
	return (u.Password != "" && u.Salt != "")
	//return (u.FullName != "" && u.Password != "" && u.Salt != "")
}

//contentInDatabase checks if the current UserContents was found in the database
func contentInDatabase(u *UserContents) bool {
	return u.EMail != ""
}

//In order to handle strange behaviour in sql
func getStringArray(arr []uint8) string {
	str := string(arr)
	return str
}
