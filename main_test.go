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

func TestBunchOfValidEmails(t *testing.T) {
	validEmails := []string{"email@domain.com", "firstname.lastname@domain.com", "email@subdomain.domain.com", "firstname+lastname@domain.com", "email@domain-one.com", "email@domain.name", "email@domain.co.jp", "firstname-lastname@domain.com", "email@domain"}
	for i := range validEmails {
		err := validateEmail(validEmails[i])
		if err != nil {
			fmt.Println(err)
			t.Fatalf("Email: " + validEmails[i] + " should be considered valid")
		}
	}
}

func TestBunchOfInValidEmails(t *testing.T) {
	invalidEmails := []string{"plainaddress", "#@%^%#$@#$@#.com", "@domain.com", "Joe Smith <email@domain.com>", "email.domain.com", "email@domain@domain.com", ".email@domain.com", "email.@domain.com", "email..email@domain.com", "あいうえお@domain.com", "email@domain.com (Joe Smith)", "email@111.222.333.44444", "email@domain..com", "1234567890@domain.com", "_______@domain.com", "root@domain.com", "localhost@domain.com", "email@123.123.123.123", "root@0.0.0.0", "root@127.0.0.1", "\"email\"@domain.com"}
	for i := range invalidEmails {
		err := validateEmail(invalidEmails[i])
		if err == nil {
			t.Fatalf("Email: " + invalidEmails[i] + " should be considered invalid")
		}
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
		fmt.Println("path: " + path + " newPath: " + newPath)
		t.Fatalf("ERROR: Strings should be equal")
	}
}

func TestDatabaseContains(t *testing.T) {
	db := connectToDatabase()
	phrase := "test@test"
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
