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
//and sets the respective fields in the DatabaseInterface
func (dbi *DatabaseInterface)SetConfigurations(f *os.File) {
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