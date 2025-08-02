package main

import (
	"fmt"
	"io"
	"net"

	"github.com/chzyer/readline"
	"github.com/thulio63/thulchat/internal/database"
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
		response,_ := cfg.CleanPrompt(rl, prompt, 1)

		if response == "y" {
			for _, server := range cfg.servers_active {
				fmt.Println(server.host + ":" + server.port)
			}
		}
	}

	//add logic to select from available options
	prompt := "Please enter the desired server host:"
	hostStr, _ := cfg.CleanPrompt(rl, prompt, 1)
	//autocomplete for known servers
	
	// cfg.colorCon.prompt.Println("Please enter the desired server host:")
	// hostStr, err := rl.Readline()
	// if err!= nil {
	// 	cfg.colorCon.err.Println("error reading input:", err)
	// 	return
	// }
	
	//autocomplete for known ports
	prompt = "Select which port to connect to:"
	portStr, _ := cfg.CleanPrompt(rl, prompt, 1)

	// cfg.colorCon.prompt.Println("Select which port to connect to:")
	// portStr, err := rl.Readline()
	// if err!= nil {
	// 	cfg.colorCon.err.Println("error reading input:", err)
	// }


	hp := hostStr + ":" + portStr
	conn, err := net.Dial("tcp", hp)
	if err != nil {
		cfg.colorCon.err.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()
	
	//clear screen here ***************************************
	cfg.colorCon.success.Println("Connected to ThulChat Server")

	for {
		//message loop





		//clear line after input












		
		break
	}

	cfg.colorCon.prompt.Println("Enter your message:")
	line, err := rl.Readline()
	if err != nil {
		cfg.colorCon.err.Println("Error reading your message:", err)
		return
	}
	
	_, err = conn.Write([]byte(line))
	if err != nil {
		cfg.colorCon.err.Println("Error sending data:", err)
	}
	cfg.colorCon.success.Println("Message sent")

	cfg.db.SendMessage(cfg.ctx, database.SendMessageParams{
		SenderID: cfg.User.UserID,
		Body: line,
		Hostname: hostStr,
		Port: portStr,
})

	//read server response
	buf := make([]byte, len(line))
	if _, err := io.ReadFull(conn, buf); err != nil {
		cfg.colorCon.err.Println("Error reading from server:", err)
		return
	}

	cfg.colorCon.info.Println("Server says:", string(buf))
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
		response, _ := cfg.CleanPrompt(rl, "Would you like to search for a user or a server?", 1)
		switch response {
		case "user":
			sel = true
			cfg.FindUser()
		case "server":
			sel = true
			cfg.FindServer()
		default:
			cfg.colorCon.err.Println("Not a valid selection. Please choose from the following: (user/server)")
		}
	}
}

func (cfg *config)FindUser() {
	//friends list, manual search

}

func (cfg *config)FindServer() {
	//CONNECT TO SQL QUERY ********************************************
	//list available servers
	num := len(cfg.servers_active)
	if num == 0 {
		cfg.colorCon.err.Println("No known available servers")
		return
	} else if num == 1 {
		cfg.colorCon.success.Println("1 known available server:")
	} else {
		mess := fmt.Sprintf("%d known available servers:", num)
		cfg.colorCon.success.Println(mess)
	}
	fmt.Println("")
	for _, serv := range cfg.servers_active {
		cfg.colorCon.info.Println(serv.host + ":" + serv.port)
	}
	fmt.Println("")

	//logic for joining
}

func (cfg *config)myIP() {
	cfg.colorCon.info.Println(cfg.MyIP.String())
	fmt.Println("")
}
