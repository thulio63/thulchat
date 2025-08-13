package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/chzyer/readline"
	"github.com/google/uuid"
	"github.com/thulio63/thulchat/internal/database"
)

//specify server to connect to
func (cfg *config)Connect() {
	rl, err := readline.NewEx(cfg.rlc.Clone())
	if err != nil {
		//change error handling
		panic(err)
	}
	defer rl.Close()

	//swap icon for this command
	pStr := rl.Config.Prompt[:len(rl.Config.Prompt)-2] + "- "
	rl.Config.Prompt = pStr
	
	validSelection := false
	var servID uuid.UUID
	hostStr := ""
	portStr := ""
	hp := ""
	active := len(cfg.servers_active)
	if active > 0 {
		prompt := fmt.Sprintf("Would you like to choose from %d known server(s)? (Y/n)", active)
		response,_ := cfg.CleanPrompt(rl, prompt, 1)

		if response == "y" {
			for i := range len(cfg.servers_active) + 1 {
				mess := ""
				if i == len(cfg.servers_active) {
					mess = fmt.Sprintf("[%d]\tConnect to an unknown server", i+1)
				} else {
					addy := fmt.Sprintf("%s:%s", cfg.servers_active[i].host, cfg.servers_active[i].port)
					mess = fmt.Sprintf("[%d]\t%s", i+1, addy)
				}
				cfg.colorCon.info.Println(mess)
			}
			
			cfg.colorCon.info.Println("")
			line, err := rl.Readline()
			if err != nil {
				cfg.colorCon.err.Println("error reading input:", err)
			} //need to loop if error happens ****************************
			nl := Clean(line)
			num, err := strconv.Atoi(nl)
			if err != nil {
				cfg.colorCon.err.Println("error converting to int:", err)
			} //need to loop if error happens ****************************

			// logic for ensuring selection is viable/in given range
			if num > 0 && num <= len(cfg.servers_active) {
				validSelection = true
				selected := cfg.servers_active[num-1].serv
				hostStr = selected.Hostname
				portStr = selected.Port
				hp = fmt.Sprintf("%s:%s", hostStr, portStr)

				//query for requested server
				servID, err = cfg.db.FindServer(cfg.ctx, database.FindServerParams{
					Hostname: hostStr, 
					Port: portStr,
				})
				if err != nil {
					cfg.colorCon.err.Println("error finding server:", err)
					return
				}

			} else if num == len(cfg.servers_active) + 1 {
				cfg.colorCon.prompt.Println("")
			} else {
				//invalid selection, loop back -- maybe unnecessary
				cfg.colorCon.err.Println("error: Not a valid selection")
			}

		} 
	}

	if !validSelection {
		//add logic to select from available options
		prompt := "Please enter the desired server host:"
		hostStr, _ = cfg.CleanPrompt(rl, prompt, 1)
		
		//autocomplete for known ports, handle invalid port numbers
		prompt = "Select which port to connect to:"
		portStr, _ = cfg.CleanPrompt(rl, prompt, 1)
	
		//query for requested server
		servID, err = cfg.db.FindServer(cfg.ctx, database.FindServerParams{
			Hostname: hostStr, 
			Port: portStr,
		})
		if err != nil {
			cfg.colorCon.err.Println("error finding server:", err)
			return
		}
	
		hp = hostStr + ":" + portStr	
	}
	//form connection
	conn, err := net.Dial("tcp", hp)
	if err != nil {
		cfg.colorCon.err.Println("error connecting to server:", err)
		return
	}
	defer conn.Close()

	
	//clear screen here ***************************************
	// this is close but not quite
	idk, err := readline.ClearScreen(rl.Stdout())
	cfg.colorCon.info.Println(idk)
	if err != nil {
		cfg.colorCon.err.Println("error clearing screen:", err)
	}
	cfg.colorCon.success.Println("Connected to ThulChat Server")
	cfg.colorCon.success.Println("")


	if rl.Config.UniqueEditLine {
		rl.Config.UniqueEditLine = false
	}

	cfg.Chat(conn, servID)
	
	//reverts readline to not disappearing edit line
	rl.Config.UniqueEditLine = true
	// cfg.colorCon.prompt.Println("Enter your message:")
	
}

func (cfg *config)Chat(conn net.Conn, servID uuid.UUID) {
	rl, err := readline.NewEx(cfg.rlc.Clone())
	if err != nil {
		//change error handling
		panic(err)
	}
	for {

		//when signal received from server, trigger printing of messages and restart loop
		//use goroutines - start one to listen for incoming data from server, and when it receives data it saves what you're typing, clears screen, prints data, and pastes what you were typing back into the line so you don't get interrupted
		//on death, kill the goroutine that handles messages from server

		//message loop
		line, err := rl.Readline()
		if err != nil {
			cfg.colorCon.err.Println("Error reading your message:", err)
			return
		}

		mess := Message{
			Body: line,
			TimeSent: time.Now(),
			Sender: cfg.User.UserID,
		}
		jDat, err := json.Marshal(mess)
		if err != nil {
			cfg.colorCon.err.Println("error marshalling data:", err)
		}
		
		_, err = conn.Write(jDat)
		if err != nil {
			cfg.colorCon.err.Println("Error sending data:", err)
		}
		//cfg.colorCon.success.Println("Message sent")

		// cfg.db.SendMessage(cfg.ctx, database.SendMessageParams{
		// 	SenderID: cfg.User.UserID,
		// 	Body: line,
		// 	Hostname: hostStr,
		// 	Port: portStr,
		// })

		// read server response
		// buf := make([]byte, len(line))
		// if _, err := io.ReadFull(conn, buf); err != nil {
		// 	cfg.colorCon.err.Println("Error reading from server:", err)
		// 	return
		// }

		// cfg.colorCon.info.Println("Server says:", string(buf))

		//utilize <- channel to know when to receive/display data and query db
		//messes, err := cfg.db.RetrieveMessages(cfg.ctx, servID)
		if err != nil {
			cfg.colorCon.err.Println("error retrieving messages:", err)
			return
		}


		//on command from user - maybe \q or something similar - send the string 'death' to the server that triggers it to shutdown
		//should this be a seperate command?
		//leaving a chatroom and deleting a server should be different commands
		//shit this is a lot man



		
		//break
	}
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
	servers, err := cfg.db.RetrieveServers(cfg.ctx)
	if err != nil {
		cfg.colorCon.err.Println("error retrieving servers:", err)
	}
	//update list of active servers
	cfg.servers_active = nil
	for _, s := range servers {
		cfg.servers_active = append(cfg.servers_active, &Server{
			host: s.Hostname,
			port: s.Port,
			context: cfg.ctx,
			serv: &database.Server{
				Hostname: s.Hostname,
				Port: s.Port,
				CreatorID: s.CreatorID,
				ServerID: s.ServerID,
				CreatedAt: s.CreatedAt,
			},
		})
	}

	num := len(servers)
	if num == 0 {
		cfg.colorCon.err.Println("No known available servers")
		cfg.colorCon.err.Println("")
		return
	} else if num == 1 {
		cfg.colorCon.success.Println("1 known available server:")
	} else {
		mess := fmt.Sprintf("%d known available servers:", num)
		cfg.colorCon.success.Println(mess)
	}
	fmt.Println("")
	for _, serv := range servers {
		cfg.colorCon.info.Println(serv.Hostname + ":" + serv.Port)
	}
	fmt.Println("")

	//logic for joining
}

func (cfg *config)myIP() {
	cfg.colorCon.info.Println(cfg.MyIP.String())
	fmt.Println("")
}

func (cfg *config)handleIncomingMessages(rl readline.Instance, c net.Conn) {
//parse data, copy current line, lock out text line, clear screen, paste incoming data, paste current line back in, unlock text line
	for {
		var buf bytes.Buffer
			_, err := buf.ReadFrom(c)
			if err != nil {
				fmt.Println("error reading from connection:", err)
			}
		var allMessages []Message //adjust server to get sender username/nickname********
		err = json.Unmarshal(buf.Bytes(), &allMessages)
		if err != nil {
			fmt.Println("error reading sent data:", err)
			return
		}

		//grab current text from message line with rl

		//lock out message line with rl

		//clear screen

		for _, mess := range allMessages {
			nextLine := fmt.Sprintf("%v:\t%s", mess.Sender, mess.Body)
			fmt.Println(nextLine)
		}

		//paste text back to message line

		//unlock message line
	}
}