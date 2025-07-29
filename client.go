package main

import (
	"fmt"
	"io"
	"net"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

//specify server to connect to
func (cfg *config)Connect() {
	rl, err := readline.New("- ")
	if err != nil {
		//change error handling
		panic(err)
	}
	defer rl.Close()
	
	active := len(cfg.servers_active)
	if active > 0 {
		prompt := fmt.Sprintf("Would you like to choose from %d known server(s)? (Y/n)", active)
		response,_ := Clean(rl, prompt, 1)

		if response == "y" {
			for _, server := range cfg.servers_active {
				fmt.Println(server.host + ":" + server.port)
			}
		}
	}

	//add logic to select from available options
	fmt.Println("Please enter the desired server host:")
	//autocomplete for known servers
	
	hostStr, err := rl.Readline()
	if err!= nil {
		fmt.Println("error reading input:", err)
	}
	
	fmt.Println("Select which port to connect to:")
	//autocomplete for known ports
	
	portStr, err := rl.Readline()
	if err!= nil {
		fmt.Println("error reading input:", err)
	}


	hp := hostStr + ":" + portStr
	conn, err := net.Dial("tcp", hp)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to ThulChat Server")

	fmt.Println("Enter your message:")
	line, err := rl.Readline()
	if err != nil {
		fmt.Println("Error reading your message:", err)
		return
	}
	
	_, err = conn.Write([]byte(line))
	if err != nil {
		fmt.Println("Error sending data:", err)
	}
	fmt.Println("Message sent")

	//read server response
	buf := make([]byte, len(line))
	if _, err := io.ReadFull(conn, buf); err != nil {
		fmt.Println("Error reading from server:", err)
		return
	}

	fmt.Println("Server says:", string(buf))
}

func (cfg *config)Find() {
	rl, err := readline.New("- ")
	if err != nil {
		//change error handling
		panic(err)
	}
	defer rl.Close()
	sel := false
	for !sel{
		response, _ := Clean(rl, "Would you like to search for a user or a server?", 1)
		switch response {
		case "user":
			sel = true
			cfg.FindUser()
		case "server":
			sel = true
			cfg.FindServer()
		default:
			fmt.Println("Not a valid selection. Please choose from the following: (user/server)")
		}
	}
}

func (cfg *config)FindUser() {
	//friends list, manual search

}

func (cfg *config)FindServer() {
	//list available servers
	num := len(cfg.servers_active)
	if num == 0 {
		color.Red("No known available servers")
		return
	} else if num == 1 {
		color.Cyan("1 known available server:")
	} else {
		mess := fmt.Sprintf("%d known available servers:", num)
		color.Cyan(mess)
	}
	fmt.Println("")
	for _, serv := range cfg.servers_active {
		color.HiWhite(serv.host + ":" + serv.port)
	}

	//logic for joining
}

func (cfg *config)myIP() {
	fmt.Println(cfg.MyIP.String())
	fmt.Println("")
}
