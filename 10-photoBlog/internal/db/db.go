package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"photoBlog/internal/auth"
	"strings"
	"time"

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
	ticker := time.NewTicker(5 * time.Second)

	conn, err := sql.Open("mysql", "testUser:password@tcp(localhost:3306)/photoBlog")
	if err != nil {
		panic(err.Error())
	}

	if err := conn.Ping(); err != nil {
		log.Fatal("ON conn.Ping", err)
	}

	DB = conn

	go checkSessions(ticker)
}

// Close - todo
func Close() {
	DB.Close()
}

func checkSessions(ticker *time.Ticker) {
	q := fmt.Sprint("SELECT sID, sActivity FROM sessions;")

	for {
		select {
		case <-ticker.C:
			res, err := DB.Query(q)
			if err != nil {
				log.Fatal("Internal server error")
			}

			var oldSessions []string

			for res.Next() {
				var sID string
				var sActivity int
				if err := res.Scan(&sID, &sActivity); err != nil {
					log.Fatal(err)
				}

				if int(time.Now().Unix())-sActivity > 60 {
					oldSessions = append(oldSessions, sID)
				}
			}

			if len(oldSessions) > 0 {
				dq := fmt.Sprintf("DELETE from sessions WHERE sID IN ('%v');", strings.Join(oldSessions, "', '"))
				_, err := DB.Query(dq)
				fmt.Println(dq)
				if err != nil {
					fmt.Println(err)
					log.Fatal("Internal server error")
				}
			}
		}
	}
}

// CreateUser - todo
func CreateUser(u User) (int, error) {
	un := u.UserName
	up, err := auth.GetHash(u.Password)
	if err != nil {
		return 0, err
	}

	q := fmt.Sprintf("INSERT INTO users (`uName`, `uPassword`) VALUES ('%v', '%v')", un, up)

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
	q := fmt.Sprintf("SELECT uId, uPassword FROM users WHERE uName = '%v';", u.UserName)
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
	uuid, err := uuid.NewV4()
	if err != nil {
		return "", errors.New("Internal server error")
	}

	sID := uuid.String()
	uID := id
	t := time.Now().Unix()

	q := fmt.Sprintf("INSERT INTO sessions (`sID`, `uID`, `sActivity`) VALUES ('%v', '%v', '%v')", sID, uID, t)
	_, err = DB.Query(q)
	if err != nil {
		return "", errors.New("Internal server error")
	}

	return sID, nil
}

// CheckSession - todo
func CheckSession(sID string) error {
	q := fmt.Sprintf("SELECT sActivity FROM sessions WHERE sID = '%v';", sID)
	res := DB.QueryRow(q)

	var sActivity string
	err := res.Scan(&sActivity)
	if err != nil {
		return errors.New("Session not available")
	}

	RefreshSession(sID)

	return nil
}

// RefreshSession - todo
func RefreshSession(sID string) error {
	nt := time.Now().Unix()
	q := fmt.Sprintf("UPDATE sessions SET sActivity = %v WHERE sID = '%v';", nt, sID)

	res := DB.QueryRow(q)
	if err := res.Scan(); err != nil {
		return errors.New("Internal server error")
	}

	return nil
}
