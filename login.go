package main

import (
	"database/sql"
	"errors"

	"github.com/chzyer/readline"
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
		username, _ := cfg.CleanPrompt(rl, "\nPlease enter your username:",1)
		//if input is incorrect, continue loop
		if username == "" {
			cfg.colorCon.err.Println("No text was entered.")
			continue
		}
		
		// query db with username, return uid if present, return uuid nil if not
		uid, err := cfg.db.FindUser(cfg.ctx, username)
		//ensures an empty table doesn't break the request
		if err != nil && !errors.Is(err, sql.ErrNoRows){
			cfg.colorCon.err.Println("Error connecting to the user database:", err)
			return 
		}

		if uid == uuid.Nil {
			cfg.colorCon.err.Println("No account found with the username", username)
			//logic for redirecting
			//fmt.Println("Would you like to try again? (Y/n)")

			response,_ := cfg.CleanPrompt(rl, "Would you like to try again? (Y/n)", 1)

			if response != "y" {
				return
			}
			//fmt.Println("\nPlease enter your username:")
			continue
		}

		cfg.colorCon.prompt.Println("")
		cfg.colorCon.prompt.Println("Enter your password:")
		pass := cfg.enterPassword()
		//test password
		pID, err := cfg.db.CheckPassword(cfg.ctx, pass)
		if err != nil && !errors.Is(err, sql.ErrNoRows){
			cfg.colorCon.err.Println("Error checking password:", err)
		}
		if pID.ID == uid {
			//success
			cfg.colorCon.success.Println("")
			cfg.colorCon.success.Println("Login successful!")
			cfg.colorCon.success.Println("")
			if pID.Nickname.Valid {
				cfg.User.Nickname = pID.Nickname.String
				cfg.colorCon.info.Println("Nickname found:", pID.Nickname.String)
				cfg.colorCon.info.Println()
			}
			cfg.User.UserID = uid
			cfg.User.CreatedAt = pID.CreatedAt
			cfg.User.UpdatedAt = pID.UpdatedAt
			cfg.User.Username = pID.Username
			return
		}
		//implement failure into the enterpassword function (maybe) for retries
		cfg.colorCon.err.Println("Incorrect password\nReturning to menu")
		return 
	}
	//return 
}			

func (cfg *config)enterPassword() []byte {
	rl, err := readline.New("* ")
	if err != nil {
		//change error handling
		panic(err)
	}
	defer rl.Close()

	line, err := rl.Readline()
	if err != nil {
		cfg.colorCon.err.Println("Error reading line:", err)
		return nil
	}
	revealed := auth.EncodePassword(line)
	return revealed
}