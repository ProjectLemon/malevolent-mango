package main

import (
	"fmt"
	"github.com/kennygrant/sanitize"
	"testing"
)

//Correctness tests
func TestValidateEmailWorking(t *testing.T) {
	email := "test@exemple.com"
	err := validateEmail(email)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestValidateEmailNonWorking(t *testing.T) {
	email := "test@"
	err := validateEmail(email)
	if err == nil {
		t.Fatalf("Should not be allowed to omit domain")
	}
}

func TestValidateEmailEmptyString(t *testing.T) {
	email := ""
	err := validateEmail(email)
	if err == nil {
		t.Fatalf("Email strings should not be allowed to be empty")
	}
}

func TestValidateEmailAgainstSQLInject(t *testing.T) {
	email := "test@exemple.com'"
	err := validateEmail(email)
	if err == nil {
		t.Fatalf("Emails should not be allowed to end on a '")
	}
}

func TestValidateEmailSQLInject(t *testing.T) {
	email := "test@user.com'select * from Users;"
	err := validateEmail(email)
	if err == nil {
		t.Fatalf("We really should not allow sql injections")
	}
}

func TestGenerateSalt(t *testing.T) {
	bytes := randBase64String(128)
	if len(bytes) < 100 {
		t.Fatal("Should have read 100b, only got: ", len(bytes))
	}
}

func TestValidateToken(t *testing.T) {
	UserID := randBase64String(64)
	token, _ := generateToken(UserID)
	session := UserSession{SessionKey: token}
	user := User{UserID: UserID, Session: &session}
	valid, _ := validateToken(&user)
	if !valid {
		t.Fatalf("Token should be valid")
	}
}

func TestEmptyToken(t *testing.T) {
	UserID := randBase64String(64)
	session := UserSession{SessionKey: ""}
	user := User{UserID: UserID, Session: &session}
	valid, _ := validateToken(&user)
	if valid {
		t.Fatalf("Empty Token should not be valid")
	}
}

func TestSanitizeFileName(t *testing.T) {
	path := "File Name!"
	newPath := sanitize.Path(path)
	if newPath == path {
		t.Fatalf("Sanitized string should not be equal to old string")
	}
}

func TestUploadSanitizer(t *testing.T) {
	path := "file-name.pdf"
	newPath := sanitizeUploadFileName(path, path[(len(path)-4):])
	if newPath != path {
		t.Fatalf("ERROR: Strings should be equal")
	}
}

func TestDatabaseContains(t *testing.T) {
	db := connectToDatabase()
	phrase := "hola"
	contains, err := db.UniversalLookup(phrase)
	if err != nil {
		fmt.Println(err)
		t.Fatalf("Recieved an unexpected error")
	}
	if !contains {
		t.Fatalf("String should be present in database")
	}
}

func TestDatabaseContainsFalse(t *testing.T) {
	db := connectToDatabase()
	phrase := "test@nonContain"
	contains, err := db.UniversalLookup(phrase)
	if err != nil {
		fmt.Println(err)
		t.Fatalf("Recieved an unexpected error")
	}
	if contains {
		t.Fatalf("String should not be present in database")
	}
}

//Benchmark tests
func BenchmarkGenerateToken(b *testing.B) {
	UserID := randBase64String(64)
	for n := 0; n < b.N; n++ {
		generateToken(UserID)
	}
}

func BenchmarkRandString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		randBase64String(128)
	}
}

func BenchmarkValidateEmail(b *testing.B) {
	for n := 0; n < b.N; n++ {
		validateEmail("linus.lagerhjelm@gmail.com")
	}
}

func BenchmarkValidateToken(b *testing.B) {
	UserID := randBase64String(64)
	token, _ := generateToken(UserID)
	session := UserSession{SessionKey: token}
	user := User{UserID: UserID, Session: &session}
	for n := 0; n < b.N; n++ {
		validateToken(&user)
	}
}

func BenchmarkUniversalLookup(b *testing.B) {
	db := connectToDatabase()
	phrase := "test@nonContain"
	for n := 0; n < b.N; n++ {
		db.UniversalLookup(phrase)
	}
}
