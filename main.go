package main

import (
	"database/sql"
	"fmt"

	"github.com/chzyer/readline"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/thulio63/thulchat/internal/database"
)

type config struct {
	db *database.Queries
	User User
	command_list map[string]cli_command
}

type cli_command struct {
	name string
	description string
	callback func()
}

func main() {
	//connect to database
	dbURL := "postgres://andrewthul:@localhost:5432/thulchat?sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Error opening database:", err)
	}
	defer db.Close()
	dbQueries := database.New(db)
	//store data on db, user id, available commands
	config := config{db: dbQueries, User: User{}}
	config.command_list = map[string]cli_command{
		"login": {
			name: "login",
			description: "Enter your email and password to log in to your account",
			callback: config.login,
		},
		"signup": {
			name: "signup",
			description: "Create an account with your email and a password",
			callback: config.sign_up,
		},
		"help": {
			name: "help",
			description: "Displays available commands and their descriptions",
			callback: config.help,
		},
		"exit": {
			name: "exit",
			description: "Exits the application",
			callback: exit,
		},
	}

	//open storage for current user info
	//thisUser := User{}


	//greeting
	fmt.Println("\nHello! Welcome to ThulChat")
	fmt.Println("For a list of available commands, type 'help'")
	fmt.Println("")
	
	//readline for repl commands
	prompt := "> "
	if config.User.Nickname != "" {
		prompt = config.User.Nickname + " > "
	} else if config.User.UserID != uuid.Nil {
		prompt = config.User.Username + " > "
	}
	rl, err := readline.New(prompt)
	if err != nil {
		//change error handling
		panic(err)
	}
	//defer fmt.Println("Closing database...")
	defer rl.Close()

	for {
		found := false
		line, err := rl.Readline()
		if err != nil {
			break
		}
		if line == "exit" {
			personal := ""
			if config.User.Username != "" {
				personal = ", " + config.User.Username
			}
			farewell := fmt.Sprintf("Closing ThulChat. Goodbye%s!", personal)
			fmt.Println(farewell)
			return
		}
		for _, command := range config.command_list {
			if line == command.name {
				command.callback()
				found = true
			}
		}
		if !found {
			response := fmt.Sprintf("%s is not a valid command", line)
			fmt.Println(response)
		} else {
			if config.User.Nickname != "" {
				rl.SetPrompt(config.User.Nickname + " > ")
			} else if config.User.UserID != uuid.Nil {
				rl.SetPrompt(config.User.Username + " > ")
			}
		}
	}	
}