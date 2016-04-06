package main

import (
	"time"
)

//User provides a struct for gathering the user information
//for the current user.
type User struct {
	Email    string
	UserId   string
	Password string
	Salt     string
	Token    string
	Session  *UserSession
}

//UserSession holds the current session information for a specific user
type UserSession struct {
	SessionKey string
	LoginTime  time.Time
	LastSeen   time.Time
}

//UserContents holds information about users name, phone, email, pdf etc
type UserContents struct {
	UserId        string
	FullName      string //Max 70 characters as suggested by: http://webarchive.nationalarchives.gov.uk/20100407120701/http://cabinetoffice.gov.uk/govtalk/schemasstandards/e-gif/datastandards.aspx
	Phone         string
	EMail         string
	ProfileIcon   string
	ProfileHeader string
	Description   string
	PDFs          []PDF
}

type PDF struct {
	Title string
	Path  string
}
