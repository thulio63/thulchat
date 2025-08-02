package main

import "github.com/fatih/color"

type colorConfig struct {
	//background color.Color
	prompt color.Color
	err color.Color
	success color.Color
	info color.Color
}

//CreateColorConfig(prompt color.Color, err color.Color, success color.Color, info color.Color)
func CreateColorConfig() colorConfig {
	promptColor := color.BgRGB(12, 12, 12).Add(color.FgHiCyan)
	errColor := color.BgRGB(12, 12, 12).Add(color.FgHiRed)
	successColor := color.BgRGB(12, 12, 12).Add(color.FgHiYellow)
	infoColor := color.BgRGB(12, 12, 12).Add(color.FgHiWhite)

	newConfig := colorConfig{
		prompt: *promptColor,
		err: *errColor,
		success: *successColor,
		info: *infoColor,
	}
	return newConfig
}