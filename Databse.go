package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
)

var (
	//ErrNoUserFound if the specified email was not in the database
	ErrNoUserFound = errors.New("User was not found in database")
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
			&user.FullName,
			&user.Password,
			&user.Salt)
		if err != nil {
			fmt.Println(err)
		}
	}

	if inDatabase(user) {
		return user, nil
	}

	return nil, ErrNoUserFound
}

//AddUser insersts the specified user into the database
//returns error where err == nil if everything went okay
func (dbi *DatabaseInterface) AddUser(user *User) error {
	_, err := dbi.DB.Exec(
		"INSERT INTO Users (EMail, FullName, Password, PasswordSalt) VALUES (?,?,?,?)",
		user.Email,
		user.FullName,
		user.Password,
		user.Salt)
	return err
}

//CloseConnection closes any active connection to the current database
func (dbi *DatabaseInterface) CloseConnection() {
	dbi.DB.Close()
}

//getConnectionString returns the connection details as a formated dataSourceName
func (dbi *DatabaseInterface) getConnectionString() string {
	return dbi.User + ":" + dbi.Password + "@" + dbi.DataSourceName
}

//inDatabase checks is the current user was found in the database
func inDatabase(u *User) bool {
	return (u.FullName != "" && u.Password != "" && u.Salt != "")
}
