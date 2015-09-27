package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func main() {

	// create a connection using a connection string /w details
	db, err := sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		fmt.Printf("Failed to connect: %s\n", err)
		return
	}

	// print the database object
	fmt.Printf("%+v\n", db)

	// insert with create on failure?  (on application load can look for table, create/update)
	res, err := db.Query("INSERT INTO posts (author, message, created) VALUES ($1, $2, $3) RETURNING id", "Casey", "A Message", time.Now().UTC())
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		_, err = db.Query("CREATE TABLE IF NOT EXISTS posts(id BIGSERIAL PRIMARY KEY NOT NULL, author TEXT NOT NULL, message TEXT NOT NULL, created TIMESTAMP)")
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		return
	}

	// print the result(s)
	fmt.Printf("%+v\n", res)
}
