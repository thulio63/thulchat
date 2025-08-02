package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"strings"

	"github.com/chzyer/readline"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/thulio63/thulchat/internal/database"
)

type config struct {
	db *database.Queries
	User *User
	command_list map[string]*cli_command
	servers_active []*Server //CONNECT TO DB ON ENTRY TO POPULATE
	//servers_count *int
	MyIP net.IP
	colorCon colorConfig
	ctx context.Context
}

type cli_command struct {
	name string
	description string
	callback func()
	visible bool
	goro bool
}

func Clean(text string) string {
	trimmed := strings.TrimSpace(text)
	lower := strings.ToLower(trimmed)
	return lower
}

func (cfg *config)CleanPrompt(rl *readline.Instance, text string, args int) (string, []string) {
	cfg.colorCon.prompt.Println(text)

	line, err := rl.Readline()
	if err != nil {
		cfg.colorCon.err.Println("error reading input:", err)
	}
	words := strings.Split(line, " ")
	for num := range args {
		words[num] = Clean(words[num])
	}
	command := words[0]
	if args == 1 {
		return command, nil
	}
	return command, words
	
}

func main() {
	//create color config for user
	clrCfg := CreateColorConfig()

	//connect to database
	dbURL := "postgres://andrewthul:@localhost:5432/thulchat?sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		clrCfg.err.Println("Error opening database:", err)
	}
	defer db.Close()
	
	dbQueries := database.New(db)
	//store data on db, user id, available commands

	var servers []*Server
	opened := false

	//retrieve list of currently operational servers
	contex := context.Background()
	sqlServers, err := dbQueries.RetrieveServers(contex)
	if err != nil {
		clrCfg.err.Println("error retrieving servers from database:", err)
	}
	
	num := 0
	for i, s := range sqlServers {
		newServ := Server{
			host: s.Hostname,
			port: s.Port,
			context: contex,
			serv: &database.Server{},
		}
		servers = append(servers, &newServ)
		num = i
	}
	clrCfg.info.Println("Number of servers found:", num)

	//get local outbound ip
	myIP := GetOutboundIP()

	//create config for current user
	//ADD WAY TO REMEMBER USER ************************
	config := config{db: dbQueries, User: &User{}, servers_active: servers, MyIP: myIP, colorCon: clrCfg, ctx: contex}

	comm_list := config.PopulateCommands()
	config.command_list = comm_list

	//greeting
	config.colorCon.success.Println("\nHello! Welcome to ThulChat")
	//fmt.Println("\nHello! Welcome to ThulChat")
	config.colorCon.info.Println("For a list of available commands, type 'help'")
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
			if !config.command_list["create"].visible {
				UpdateVisibile(true ,config.command_list, "create")
				continue
			} else if len(config.servers_active) != 0 && !opened{
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
			config.colorCon.info.Println("deleting servers...")
			_, err := config.db.DeleteServer(config.ctx, config.User.UserID)
			if err != nil {
				config.colorCon.err.Println("error deleting servers:", err)
			}
			// for s := len(config.servers_active) - 1; s >= 0; s--{
			// 	config.db.DeleteServer(config.ctx, config.servers_active[s].serv.ServerID)
			// }

			farewell := fmt.Sprintf("\nClosing ThulChat. Goodbye%s!", personal)
			config.colorCon.success.Println(farewell)
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
			response := fmt.Sprintf("'%s' is not a valid command", trimmed)
			config.colorCon.err.Println(response)
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
	//color.Unset()
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