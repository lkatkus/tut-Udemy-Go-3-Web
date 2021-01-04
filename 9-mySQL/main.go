package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("MYSQL TEST")

	db, err := sql.Open("mysql", "")
	if err != nil {
		log.Fatal("ON sql.Open", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("ON db.Ping", err)
	}

	res, err := db.Query("SELECT aId, aName FROM test02.amigos")
	if err := db.Ping(); err != nil {
		log.Fatal("ON db.Query", err)
	}

	for res.Next() {
		var id, name string
		if err := res.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%v - %v\n", id, name)
	}
}
