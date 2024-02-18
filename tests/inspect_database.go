//  This Go script displays all tables in a given database
//  along with their respective columns.

package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// checkError helper function to handle errors
func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func showTablesWithColumns() {
	// sql.Open does not return a connection. It just returns a handle to the database.
	// In a real world scenario, those db credentials could be environment variables
	// and we could use a package like github.com/kelseyhightower/envconfig to read them.
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/walletapp")

	// A defer statement pushes a function call onto a list.
	// The list of saved calls is executed after the surrounding function returns.
	// Defer is commonly used to simplify functions that perform various clean-up actions.
	defer db.Close()

	checkError("Error getting a handle to the database", err)

	// Now it's time to validate the Data Source Name (DSN) to check if the connection
	// can be correctly established.
	err = db.Ping()

	checkError("Error establishing a connection to the database", err)

	showTablesQuery, err := db.Query("SHOW TABLES")

	defer showTablesQuery.Close()

	checkError("Error creating the query", err)

	for showTablesQuery.Next() {
		var tableName string

		// Get table name
		err = showTablesQuery.Scan(&tableName)

		checkError("Error querying tables", err)

		selectQuery, err := db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))

		defer selectQuery.Close()

		checkError("Error creating the query", err)

		// Get column names from the given table
		columns, err := selectQuery.Columns()
		if err != nil {
			checkError(fmt.Sprintf("Error getting columns from table %s", tableName), err)
		}

		fmt.Printf("table name: %s -- columns: %v\n", tableName, strings.Join(columns, ", "))
	}
}

func main() {
	showTablesWithColumns()
}
