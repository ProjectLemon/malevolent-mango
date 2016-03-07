package main

import (
    "os"
    "database/sql"
    "bufio"
    "strings"
    _ "github.com/go-sql-driver/mysql"
)

//DatabaseInterface represent a configuration object, containing configurations
// for the current database
type DatabaseInterface struct {
    User string
    Password string
    DriverName string
    DataSourceName string
    DB *sql.DB
}

//SetConfigurations reads the specified config file
//and sets the respective fields
func (dbi *DatabaseInterface)SetConfigurations(f *os.File) {
    cnf := make(map[string]string)
    
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        insert := strings.SplitAfter(scanner.Text(), " ")
        cnf[insert[0]] = insert[1]
    }
    dbi.User = cnf["User "]
    dbi.Password = cnf["Password "]
    dbi.DriverName = cnf["DriverName "]
    dbi.DataSourceName = cnf["DataSourceName "]
}

//OpenConnection connects to a database using the DatabaseInterface
func (dbi *DatabaseInterface)OpenConnection() error {
    db, err := sql.Open(dbi.DriverName, dbi.getConnectionString())
    if err != nil {
        return err
    }
    err = db.Ping()
    if err != nil {
        return err
    }
    dbi.DB = db
    return nil;
}

//getConnectionString returns the connection details as a formated dataSourceName
func (dbi *DatabaseInterface)getConnectionString() string {
    return dbi.User+":"+dbi.Password+"@"+dbi.DataSourceName
}