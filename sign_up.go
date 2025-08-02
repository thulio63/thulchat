package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chzyer/readline"
	"github.com/google/uuid"
	"github.com/thulio63/thulchat/internal/auth"
	"github.com/thulio63/thulchat/internal/database"
)

func (cfg *config)sign_up() {
	if cfg.User.UserID != uuid.Nil {
		fmt.Println("You are already signed in as " + cfg.User.Username)
		return
	}

	fmt.Println("Create a username:")


	//readline for repl commands
	rl, err := readline.New("- ")
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
		username := line
		//if input is incorrect, continue loop
		if username == "" {
			fmt.Println("Please enter a username:")
			continue
		}
		//check if email already exists
		info := context.Background()
		uid, err := cfg.db.FindUser(info, username)
		//ensures an empty table doesn't break the request
		if err != nil && !errors.Is(err, sql.ErrNoRows){
			fmt.Println("Error connecting to the user database:", err)
			return 
		}
		if uid != uuid.Nil {
			fmt.Println("There is already a user named " + username)
			//redirect logic
			fmt.Println("Would you like to try again? (Y/n)")
			line, err = rl.Readline()
			if err != nil {
				fmt.Println("Error reading:", err)
				break
			}
			if line != "Y" && line != "y" {
				return
			}
			continue
		}
		//if not, make new user - request password
		fmt.Println("Choose a password:")
		pass := password_create()
		for pass == nil {
			pass = password_create()
		}

		//add password to params for user creation
		newUser, err := cfg.db.CreateUser(info, database.CreateUserParams{
			Username: username,
			Password: pass})
		if err != nil {
			fmt.Println("Error creating user:", err)
		}
		cfg.User.Username = newUser.Username
		cfg.User.Password = string(newUser.Password)
		cfg.User.UserID = newUser.ID
		cfg.User.CreatedAt = newUser.CreatedAt
		cfg.User.UpdatedAt = newUser.UpdatedAt
		response := fmt.Sprintf("New user created! \nYour username is %s", newUser.Username)
		fmt.Println(response)
		fmt.Println("Would you like to create a nickname? (Y/n)")
		nick := nickname_create(cfg.User.Username)
		cfg.User.Nickname = nick
		fmt.Println("")
		//cfg.pass = newUser.pass
		return 
	}
}

func password_create() []byte {
	//readline for repl commands
	rl, err := readline.New("* ")
	if err != nil {
		//change error handling
		panic(err)
	}
	defer rl.Close()

	line, err := rl.Readline()
	if err != nil {
		fmt.Println("Error creating password:", err)
		return nil
	}
	//password requirements/validation here
	//not in for loop currently, implement loop when making validation logic
	if line == "" {
		fmt.Println("Please input a valid password")
		return nil
	}
	pass := auth.EncodePassword(line)
	return pass
}

func nickname_create(username string) string {
	//readline for repl commands
	rl, err := readline.New("- ")
	if err != nil {
		//change error handling
		panic(err)
	}
	defer rl.Close()

	line, err := rl.Readline()
	if err != nil {
		//change error handling
		fmt.Println("Error creating nickname:", err)
	}
	if line != "Y" && line!="y" {
		return username
	}
	fmt.Println("Choose your nickname:")
	line, err = rl.Readline()
	if err != nil {
		//change error handling
		fmt.Println("Error creating nickname:", err)
	}
	return line
}