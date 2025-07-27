package main

import "fmt"

func (cfg *config)help() {
	fmt.Println("\nCommand options:")
	fmt.Println("")
	for key, val := range cfg.command_list {
		if val.visible {
			mess := fmt.Sprintf("%s: %s", key, val.description)
			fmt.Println(mess)
		}
	}
	fmt.Println("")
	//return 
}