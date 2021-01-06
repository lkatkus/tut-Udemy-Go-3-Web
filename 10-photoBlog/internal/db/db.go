package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"photoBlog/internal/auth"

	uuid "github.com/satori/go.uuid"
)

// User - todo
type User struct {
	UserName string
	Password string
}

// DB - todo
var DB *sql.DB

// Connect - todo
func Connect() {
	conn, err := sql.Open("mysql", "testUser:password@tcp(localhost:3306)/photoBlog")
	if err != nil {
		panic(err.Error())
	}

	if err := conn.Ping(); err != nil {
		log.Fatal("ON conn.Ping", err)
	}

	DB = conn
}

// Close - todo
func Close() {
	DB.Close()
}

// CreateUser - todo
func CreateUser(u User) (int, error) {
	un := u.UserName
	up, err := auth.GetHash(u.Password)
	if err != nil {
		return 0, err
	}

	q := fmt.Sprintf("INSERT INTO `photoBlog`.`users` (`uName`, `uPassword`) VALUES ('%v', '%v')", un, up)

	res, err := DB.Query(q)
	if err != nil {
		fmt.Println(err)
		return 0, errors.New("Username already in use")
	}

	for res.Next() {
		var name string
		if err := res.Scan(&name); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("name - %v\n", name)
	}

	return 1, nil
}

// LoginUser - todo
func LoginUser(u User) (string, error) {
	q := fmt.Sprintf("SELECT uId, uPassword FROM photoBlog.users WHERE uName = '%v';", u.UserName)
	res := DB.QueryRow(q)

	var id string
	var password string
	if err := res.Scan(&id, &password); err != nil {
		return "", errors.New("Wrong username of password")
	}

	err := auth.CheckHash(password, u.Password)
	if err != nil {
		return "", err
	}

	sID, err := CreateSession(id)
	if err != nil {
		return "", err
	}

	return sID, nil
}

// CreateSession - todo
func CreateSession(id string) (string, error) {
	sID, err := uuid.NewV4()
	if err != nil {
		return "", errors.New("Internal server error")
	}

	q := fmt.Sprintf("INSERT INTO `photoBlog`.`sessions` (`sID`, `uID`) VALUES ('%v', '%v')", sID.String(), id)
	_, err = DB.Query(q)
	if err != nil {
		return "", errors.New("Internal server error")
	}

	return sID.String(), nil
}

// CheckSession - todo
func CheckSession(sID string) error {
	q := fmt.Sprintf("SELECT uId FROM photoBlog.sessions WHERE sID = '%v';", sID)
	res := DB.QueryRow(q)

	var id string
	if err := res.Scan(&id); err != nil {
		return errors.New("Session not available")
	}

	return nil
}
