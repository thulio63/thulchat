package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/thulio63/thulchat/internal/auth"
)

func (cfg *config)login() {
	
	//fmt.Println("\nPlease enter your username:")


	//readline for repl commands
	rl, err := readline.New("- ")
	if err != nil {
		//change error handling
		panic(err)
	}
	defer rl.Close()

	for {
		// line, err := rl.Readline()
		// if err != nil {
		// 	break
		// }
		//implement verification of proper format... regex?
		username, _ := Clean(rl, "\nPlease enter your username:",1)
		//if input is incorrect, continue loop
		if username == "" {
			fmt.Println("No text was entered.")
			continue
		}
		
		// query db with username, return uid if present, return uuid nil if not
		info := context.Background()
		uid, err := cfg.db.FindUser(info, username)
		//ensures an empty table doesn't break the request
		if err != nil && !errors.Is(err, sql.ErrNoRows){
			color.Red("Error connecting to the user database:", err)
			return 
		}

		if uid == uuid.Nil {
			fmt.Println("No account found with the username", username)
			//logic for redirecting
			//fmt.Println("Would you like to try again? (Y/n)")

			response,_ := Clean(rl, "Would you like to try again? (Y/n)", 1)

			if response != "y" {
				return
			}
			//fmt.Println("\nPlease enter your username:")
			continue
		}
		fmt.Println("Enter your password:")
		pass := enterPassword()
		//test password
		pID, err := cfg.db.CheckPassword(info, pass)
		if err != nil && !errors.Is(err, sql.ErrNoRows){
			color.Red("Error checking password:", err)
		}
		if pID.ID == uid {
			//success
			color.Yellow("Login successful!")
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
		color.Red("Incorrect password\nReturning to menu")
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
		color.Red("Error reading line:", err)
		return nil
	}
	revealed := auth.EncodePassword(line)
	return revealed
}