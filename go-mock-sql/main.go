package main

import (
	"fmt"

	"database/sql"
)

// @note: the problem is that I want to test in isolation
// I have or can write integration tests to verify the SQL
// all I care about in isolation is that my code is valid
// hence why the lack of sql mocks is a serious issue

// @link: https://github.com/DATA-DOG/go-sqlmock
// the above is the closest to a solution
// however, imo it's super overkill

// We could try to fully stub-out the database system but
// we end up hitting a wall with driverConn being a private type

type DbInterfaceOne interface {
	Begin() (*sql.Tx, error)
}

type DbInterfaceTwo interface {
	Begin() (interface{}, error)
}

func main() {
	fmt.Println("Isolating the standard database library does not appear to be possible.")
	fmt.Println("We can create an interface if we use the exact type `*Tx`, but not another interface.")

	// this will work, because we matched an actual return type
	var TestDbOne DbInterfaceOne
	TestDbOne = &sql.DB{}
	fmt.Printf("%+v\n", TestDbOne)

	// but we can't use an interface, so we can't control the return value
	var TestDbTwo DbInterfaceTwo
	TestDbTwo = &sql.DB{}
	fmt.Printf("%+v\n", TestDbTwo)
}
