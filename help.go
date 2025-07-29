package main

import (
	"fmt"

	"github.com/fatih/color"
)

func (cfg *config)help() {
	color.Cyan("\nCommand options:")
	fmt.Println("")
	myColor := color.BgRGB(12, 12, 12)
	myColor.Add(color.FgHiWhite)
	for key, val := range cfg.command_list {
		if val.visible {
			mess := fmt.Sprintf("%s: %s", key, val.description)
			myColor.Println(mess)
		}
	}
	fmt.Println("")
	//return 
}