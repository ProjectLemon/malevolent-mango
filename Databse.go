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
func (dbi *DatabaseInterface) LookupUser(user *User) (bool, error) {
	rows, err := dbi.DB.Query("SELECT * FROM Users WHERE EMail = '" + user.Email.String + "'")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&user.Email,
			&user.FullName,
			&user.PasswordHash,
			&user.Salt)
		if err != nil {
			fmt.Println(err)
		}
	}

	if user.Email.Valid && user.FullName.Valid && user.PasswordHash.Valid && user.Salt.Valid {
		return true, nil
	}
	return false, ErrNoUserFound
}

//CloseConnection closes any active connection to the current database
func (dbi *DatabaseInterface) CloseConnection() {
	dbi.DB.Close()
}

//getConnectionString returns the connection details as a formated dataSourceName
func (dbi *DatabaseInterface) getConnectionString() string {
	return dbi.User + ":" + dbi.Password + "@" + dbi.DataSourceName
}
