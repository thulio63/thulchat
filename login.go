package main

import (
	"context"
	"fmt"

	"github.com/chzyer/readline"
	"github.com/google/uuid"
)

func (cfg *config)login() {
	
	fmt.Println("\nPlease enter your email address:")


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
		
		// query db with email, return uid if present, return uuid nil if not
		info := context.Background()
		uid, err := cfg.db.FindUser(info, email)
		if err != nil {
			fmt.Println("Error connecting to the user database:", err)
			return 
		}

		if uid == uuid.Nil {
			//logic for redirecting
			fmt.Println("No account found with the email address of", email)
			fmt.Println("\nPlease enter your email address:")
			continue
		}

		pass := enterPassword()
		//test password
		pID, err := cfg.db.CheckPassword(info, pass)
		if err != nil {
			fmt.Println("Error checking password:", err)
		}
		if pID == uid {
			//success
			fmt.Println("Login successful")
			cfg.UID = uid
			return
		}
		//implement failure into the enterpassword function (maybe) for retries
		fmt.Println("Incorrect password")
		return 
	}
	//return 
}			

func enterPassword() string {
	rl, err := readline.New("* ")
	if err != nil {
		//change error handling
		panic(err)
	}
	defer rl.Close()

	line, err := rl.Readline()
	if err != nil {
		fmt.Println("Error reading line:", err)
		return ""
	}

	return line
}