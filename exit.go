package main

import (
	"fmt"
	"os"
)

func exit() {
	fmt.Println("Disconnecting from the database...")
	
	//make sure database is disconnected
	fmt.Println("Closing ThulChat. Goodbye!")

	os.Exit(0)
	//return
}