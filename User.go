package main

import "database/sql"

//User provides a datatype for gathering the user information
//for the current user.
type User struct {
	Email        sql.NullString
	FullName     sql.NullString
	PasswordHash sql.NullString
	Salt         sql.NullString
}
