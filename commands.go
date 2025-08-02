package main

func (config *config)PopulateCommands() map[string]*cli_command{
	command_list := map[string]*cli_command{
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
			"create": {
				name: "create",
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
		return command_list
}