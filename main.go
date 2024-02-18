package main

import (
	"wallet-service-gin/api"
	"wallet-service-gin/dbhelper"
)

// main function where all the code starts
func main() {
	// setting up the DB
	DB, _ := dbhelper.SetupDB()

	// running the API in the machine
	api.API(DB)
}
