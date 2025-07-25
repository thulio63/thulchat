package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chzyer/readline"
	"github.com/google/uuid"
	"github.com/thulio63/thulchat/internal/database"
)

func (cfg *config)sign_up() {
	fmt.Println("Enter your email address:")

	//readline for repl commands
	rl, err := readline.New("> ")
	if err != nil {
		//change error handling
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		
		//implement verification of proper format... regex?
		email := line
		//if input is incorrect, continue loop
		if email == "" {
			fmt.Println("Please enter a valid email address:")
			continue
		}
		//check if email already exists
		info := context.Background()
		uid, err := cfg.db.FindUser(info, email)
		//ensures an empty table doesn't break the request
		if err != nil && !errors.Is(err, sql.ErrNoRows){
			fmt.Println("Error connecting to the user database:", err)
			return 
		}
		if uid != uuid.Nil {
			fmt.Println("There is already a user listed under the email", line)
			//logic for redirecting
			continue
		}
		//if not, make new user - request password
		fmt.Println("Choose a password:")
		pass := password_create()

		//add password to params for user creation
		newUser, err := cfg.db.CreateUser(info, database.CreateUserParams{
			Email: email,
			Password: pass})
		if err != nil {
			fmt.Println("Error creating user:", err)
		}
		cfg.UID = newUser.ID
		response := fmt.Sprintf("New user created with email %s!", newUser.Email)
		fmt.Println(response)
		//cfg.pass = newUser.pass
		return 
	}
}

func password_create() string {
	//readline for repl commands
	rl, err := readline.New("* ")
	if err != nil {
		//change error handling
		panic(err)
	}
	defer rl.Close()

	line, err := rl.Readline()
	//password requirements/validation here
	//not in for loop currently, implement loop when making validation logic
	if line == "" {
		fmt.Println("Please input a valid password")
		//return
	}
	if err != nil {
		fmt.Println("Error creating password:", err)
		return ""
	}

	return line
}