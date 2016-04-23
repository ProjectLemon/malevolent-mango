package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
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

//UniversalLookup will query the database for the given string. It will
//search through all the rows and columns in every table in the database.
//This function should be considered expensive and be used with caution.
func (dbi *DatabaseInterface) UniversalLookup(phrase string) (bool, error) {
	//Were hardcoding the database structure into a map
	// in order to boost performance by reducing sql-queries
	tables := make(map[string][]string)
	tables["Users"] = []string{"EMail", "UserId", "Password", "PasswordSalt"}
	tables["UserSession"] = []string{"SessionKey", "UserId", "LoginTime", "LastSeenTime"}
	tables["UserContent"] = []string{"UserId", "FullName", "Phone", "EMail", "ProfileIcon", "ProfileHeader", "Description", "PDFs"}

	for tableName, columnNames := range tables {
		for i := 0; i < len(columnNames); i++ {
			rows, err := dbi.DB.Query("SELECT * FROM "+tableName+" WHERE "+columnNames[i]+"=?", phrase)
			if err != nil {
				return false, err
			}
			for rows.Next() {
				rows.Close()
				return true, nil
			}
			rows.Close()
		}
	}
	return false, nil
}

//LookupUser sends a query to the database for the specified
//username and password hash. Returns error if query failed
func (dbi *DatabaseInterface) LookupUser(user *User) (*User, error) {
	rows, err := dbi.DB.Query("SELECT * FROM Users WHERE EMail = ?", user.Email)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&user.Email,
			&user.UserID,
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
		user.UserID,
		user.Password,
		user.Salt)

	_, err = dbi.DB.Exec(
		"INSERT INTO `UserContent` (`UserId`, `FullName`, `Phone`, `EMail`, `ProfileIcon`, `ProfileHeader`, `Description`, `PDFs`) VALUES (?, ?, ?, ?, ?, ?,? ,?)",
		user.UserID,
		"",
		"",
		"",
		"",
		"",
		"",
		"")
	return err
}

//GetUserContents looks up, and return, user content in database
func (dbi *DatabaseInterface) GetUserContents(uid string, userContent *UserContents) (*UserContents, error) {
	rows, err := dbi.DB.Query("SELECT * FROM UserContent WHERE UserId=?", uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jsonField []uint8

	for rows.Next() {
		err := rows.Scan(
			&userContent.UserID,
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

//UpdateUserContent inserts the specified UserContent
//for the specified UserId into the database
func (dbi *DatabaseInterface) UpdateUserContent(uid string, uc *UserContents) error {
	var buffer bytes.Buffer

	invalidContent := validateUserContent(uc)
	if invalidContent {
		return errors.New("Invalid content")
	}
	buffer.WriteRune('[')
	for i := 0; i < len(uc.PDFs); i++ {
		buffer.WriteString(uc.PDFs[i].String())
		buffer.WriteRune(',')
	}
	buffer.WriteRune(']')

	_, err := dbi.DB.Exec("UPDATE UserContent set UserId=?, FullName=?, Phone=?, EMail=?, ProfileIcon=?, ProfileHeader=?, Description=?, PDFs=? WHERE UserId=?;",
		uid,
		uc.FullName,
		uc.Phone,
		uc.EMail,
		uc.ProfileIcon,
		uc.ProfileHeader,
		uc.Description,
		buffer.String(),
		uid)
	return err
}

//InsertUserSession creates a new row in the database for a user session
func (dbi *DatabaseInterface) InsertUserSession(user *User) error {
	_, err := dbi.DB.Exec(
		"INSERT INTO UserSession (SessionKey, UserId, LoginTime, LastSeenTime) VALUES (?,?,?,?)",
		user.Token,
		user.UserID,
		time.Now().Format(time.RFC3339),
		time.Now().Format(time.RFC3339))
	return err
}

//UpdateUserSession overwrites the current token value and last seen time in database
func (dbi *DatabaseInterface) UpdateUserSession(user *User) error {
	_, err := dbi.DB.Exec("UPDATE UserSession set SessionKey=?, LastSeenTime=? WHERE UserId =?;", user.Token, time.Now().Format(time.RFC3339), user.UserID)
	return err
}

//GetUserSession reads the user session for the specified user
//into the user session field of the struct
func (dbi *DatabaseInterface) GetUserSession(user *User) (*User, error) {
	rows, err := dbi.DB.Query("SELECT * FROM UserSession WHERE SessionKey=?", user.Token)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user.Session = new(UserSession)
	for rows.Next() {
		err := rows.Scan(
			&user.Session.SessionKey,
			&user.UserID,
			&user.Session.LoginTime,
			&user.Session.LastSeen)
		if err != nil {
			fmt.Println(err)
		}
	}

	if user.UserID != "" {
		return user, nil
	}

	return nil, ErrNoActiveSession
}

//RemoveUserSession removes the session entry in database
//that has the provided session key
func (dbi *DatabaseInterface) RemoveUserSession(session *UserSession) error {
	_, err := dbi.DB.Exec("DELETE FROM UserSession WHERE SessionKey=?", session.SessionKey)
	return err
}

//CleanUserSession gets the user sessions that has a larger last seen time
//than 10 minutes
func (dbi *DatabaseInterface) CleanUserSession() error {
	rows, err := dbi.DB.Query("SELECT * FROM UserSession WHERE LastSeenTime <= ?", (time.Now().Add(time.Minute * -10)).Format(time.RFC3339))
	if err != nil {
		return err
	}
	defer rows.Close()

	userSession := new(UserSession)
	for rows.Next() {
		err := rows.Scan(
			&userSession.SessionKey,
			new(string),
			&userSession.LoginTime,
			&userSession.LastSeen)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = dbi.RemoveUserSession(userSession)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return err
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
func getStringArray(arr []uint8) []PDF {
	str := string(arr)
	var pdfs []PDF
	json.Unmarshal([]byte(str), &pdfs)
	return pdfs
}

func validateUserContent(uc *UserContents) bool {
	return !(len(uc.FullName) < 70 &&
		len(uc.Phone) < 50 &&
		len(uc.EMail) < 80 &&
		len(uc.ProfileIcon) < 150 &&
		len(uc.ProfileHeader) < 150 &&
		len(uc.Description) < 360 &&
		len(uc.PDFs) < 21844)
}
