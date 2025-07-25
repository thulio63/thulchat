package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"

	"github.com/chzyer/readline"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/thulio63/thulchat/internal/database"
)

type config struct {
	db *database.Queries
	UID uuid.UUID
}

func main() {
	input := bufio.NewScanner(os.Stdin)
	dbURL := "postgres://andrewthul:@localhost:5432/thulchat?sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Error opening database:", err)
	}
	defer db.Close()
	dbQueries := database.New(db)
	config := config{db: dbQueries}

	
	fmt.Println("Hello! Welcome to ThulChat")
	
	//readline for repeated repl segments. necessary?
	rl, err := readline.New("> ")
	if err != nil {
		//change error handling
		fmt.Println("Error reading previous line")
	}
	defer rl.Close()


	id := config.login(*input)	
	if id == uuid.Nil {
		fmt.Println("User not found. Would you like to make an account?")
	}
	
}