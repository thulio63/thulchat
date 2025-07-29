package main

import (
	"database/sql"
	"fmt"
	"net"
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/thulio63/thulchat/internal/database"
)

type config struct {
	db *database.Queries
	User *User
	command_list map[string]*cli_command
	servers_active []*Server
	servers_count *int
	MyIP net.IP
}

type cli_command struct {
	name string
	description string
	callback func()
	visible bool
	goro bool
}

func Clean(rl *readline.Instance, text string, args int) (string, []string) {
	fmt.Println(text)

	line, err := rl.Readline()
	if err != nil {
		fmt.Println("error reading input:", err)
	}
	trimmed := strings.TrimSpace(line)
	words := strings.Split(trimmed, " ")
	for num := range args {
		words[num] = strings.ToLower(words[num])
	}
	command := words[0]
	if args == 1 {
		return command, nil
	}
	return command, words
	
}

func Space() {
	
}

func main() {
	//connect to database
	dbURL := "postgres://andrewthul:@localhost:5432/thulchat?sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Error opening database:", err)
	}
	defer db.Close()
	dbQueries := database.New(db)
	//store data on db, user id, available commands
	empty := 0
	var servers []*Server
	opened := false
	//get local outbound ip
	myIP := GetOutboundIP()
	config := config{db: dbQueries, User: &User{}, servers_count: &empty, servers_active: servers, MyIP: myIP}
	config.command_list = map[string]*cli_command{
		"login": {
			name: "login",
			description: "Enter a username and password to log in to your account",
			callback: config.login,
			visible: true,
			goro: false,
		},
		"signup": {
			name: "signup",
			description: "Create an account with a username and a password",
			callback: config.sign_up,
			visible: true,
			goro: false,
		},
		"help": {
			name: "help",
			description: "Displays available commands and their descriptions",
			callback: config.help,
			visible: true,
			goro: false,
		},
		"exit": {
			name: "exit",
			description: "Exits the application",
			callback: exit,
			visible: true,
			goro: false,
		},
		"s_c": {
			name: "s_c",
			description: "Creates a server for communication",
			callback: config.New, 
			visible: false,
			goro: false,
		},
		"connect": { // change to enter, make funtion for loop to create "chatroom"
			name: "connect",
			description: "Connect to a server",
			callback: config.Connect,
			visible: false,
			goro: false,
		},
		"find": {
			name: "find",
			description: "Search for other users or available servers",
			callback: config.Find,
			visible: false,
			goro: false,
		},
		"myip": {
			name: "myip",
			description: "Prints outbound IP address for this device",
			callback: config.myIP,
			visible: true,
			goro: false,
		},
	}
	
	//myColor := color.BgRGB(12, 12, 12)
	color.Set().AddBgRGB(12, 12,12)
	color.Set().AddRGB(122, 231, 95)

	//greeting
	color.Cyan("\nHello! Welcome to ThulChat")
	//fmt.Println("\nHello! Welcome to ThulChat")
	color.HiWhite("For a list of available commands, type 'help'")
	fmt.Println("")
	
	//myColor.Add(color.FgGreen)
	prompt := "> "
	if config.User.Nickname != "" {
		prompt = config.User.Nickname + " > "
	} else if config.User.UserID != uuid.Nil {
		prompt = config.User.Username + " > "
	}
	rl, err := readline.New(prompt)
	if err != nil {
		//change error handling
		panic(err)
	}
	//defer fmt.Println("Closing database...")
	defer rl.Close()

	for {
		//logic for revealing commands to user
		if config.User.UserID != uuid.Nil {
			if !config.command_list["s_c"].visible {
				UpdateVisibile(true ,config.command_list, "s_c")
				continue
			} else if *config.servers_count != 0 && !opened{
				UpdateVisibile(true, config.command_list, "connect")
				//prevents infinite loop
				opened = true
				continue
			} else if !config.command_list["find"].visible {
				UpdateVisibile(true, config.command_list, "find")
			}
		}

		//flag for confirming command is valid
		found := false
		line, err := rl.Readline()
		if err != nil {
			break
		}
		trimmed := strings.TrimSpace(line)
		if trimmed == "exit" {
			//ui.Close()
			personal := ""
			if config.User.Username != "" {
				personal = ", " + config.User.Username
			}
			farewell := fmt.Sprintf("\nClosing ThulChat. Goodbye%s!", personal)
			color.Cyan(farewell)
			fmt.Println("")
			return
		}
		for _, command := range config.command_list {
			if trimmed == command.name && command.visible{
				if command.goro {
					go command.callback()
				} else {
					command.callback()
				}
				//fmt.Println("command executed")
				found = true
			}
		}
		//fmt.Println("loop escaped")
		
		if !found {
			response := fmt.Sprintf("%s is not a valid command", trimmed)
			color.Red(response)
			fmt.Println("")
		} else { //sets prompt name to nickname or username
			if config.User.Nickname != "" {
				rl.SetPrompt(config.User.Nickname + " > ")
				continue
			} else if config.User.UserID != uuid.Nil {
				rl.SetPrompt(config.User.Username + " > ")
				continue
			} else {
				continue
			}
		}
	}	
	color.Unset()
}

func UpdateVisibile(change bool,list map[string]*cli_command, comms ...string) {

	for _, comm := range comms {
		for k, v := range list {
			if k == comm {
				v.visible = change
			}
		}
	}

}