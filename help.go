package main

import (
	"fmt"
)

func (cfg *config)help() {
	cfg.colorCon.prompt.Println("\nCommand options:")
	fmt.Println("")
	// myColor := color.BgRGB(12, 12, 12)
	// myColor.Add(color.FgHiWhite)
	for key, val := range cfg.command_list {
		if val.visible {
			mess := fmt.Sprintf("%s: %s", key, val.description)
			cfg.colorCon.info.Println(mess)
		}
	}
	fmt.Println("")
	//return 
}