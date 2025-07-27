package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chzyer/readline"
	"github.com/google/uuid"
	"github.com/thulio63/thulchat/internal/auth"
)

func (cfg *config)login() {
	
	fmt.Println("\nPlease enter your username:")


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
		
		// query db with username, return uid if present, return uuid nil if not
		info := context.Background()
		uid, err := cfg.db.FindUser(info, username)
		//ensures an empty table doesn't break the request
		if err != nil && !errors.Is(err, sql.ErrNoRows){
			fmt.Println("Error connecting to the user database:", err)
			return 
		}

		if uid == uuid.Nil {
			fmt.Println("No account found with the username", username)
			//logic for redirecting
			fmt.Println("Would you like to try again? (Y/n)")
			line, err = rl.Readline()
			if err != nil {
				fmt.Println("Error reading:", err)
				break
			}
			if line != "Y" && line != "y" {
				return
			}
			fmt.Println("\nPlease enter your username:")
			continue
		}
		fmt.Println("Enter your password:")
		pass := enterPassword()
		//test password
		pID, err := cfg.db.CheckPassword(info, pass)
		if err != nil && !errors.Is(err, sql.ErrNoRows){
			fmt.Println("Error checking password:", err)
		}
		if pID.ID == uid {
			//success
			fmt.Println("Login successful!")
			if pID.Nickname.Valid {
				cfg.User.Nickname = pID.Nickname.String
			}
			cfg.User.UserID = uid
			cfg.User.CreatedAt = pID.CreatedAt
			cfg.User.UpdatedAt = pID.UpdatedAt
			cfg.User.Username = pID.Username
			return
		}
		//implement failure into the enterpassword function (maybe) for retries
		fmt.Println("Incorrect password\nReturning to menu")
		return 
	}
	//return 
}			

func enterPassword() []byte {
	rl, err := readline.New("* ")
	if err != nil {
		//change error handling
		panic(err)
	}
	defer rl.Close()

	line, err := rl.Readline()
	if err != nil {
		fmt.Println("Error reading line:", err)
		return nil
	}
	revealed := auth.EncodePassword(line)
	return revealed
}