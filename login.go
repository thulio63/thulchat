package main

import (
	"bufio"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (cfg *config)login(input bufio.Scanner) (uuid.UUID) {
	
	fmt.Println("\nPlease enter your email:")
	input.Scan()
	//implement verification of proper format... regex?
	email := input.Text()
	// query db with email, return uid if present, return uuid nil if not
	info := context.Background()
	uid, err := cfg.db.FindUser(info, email)
	if err != nil {
		fmt.Println("Error connecting to the user database:", err)
		return uuid.Nil
	}
	
	return uid
}			