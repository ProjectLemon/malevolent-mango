package main

import (
	"encoding/json"
	"time"
)

//User provides a struct for gathering the user information
//for the current user.
type User struct {
	Email    string
	UserID   string
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
	UserID        string
	FullName      string //Max 70 characters as suggested by: http://webarchive.nationalarchives.gov.uk/20100407120701/http://cabinetoffice.gov.uk/govtalk/schemasstandards/e-gif/datastandards.aspx
	Phone         string
	EMail         string
	ProfileIcon   string
	ProfileHeader string
	Description   string
	PDFs          []PDF
}

//PDF represents a pdf file. Containing a Title and a search path
type PDF struct {
	Title string
	Path  string
}

func (pdf *PDF) String() string {
	str, _ := json.Marshal(pdf)
	return string(str)
}
