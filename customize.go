package main

import (
	"time"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/thulio63/thulchat/internal/database"
)

type colorConfig struct {
	//background color.Color
	prompt color.Color
	err color.Color
	success color.Color
	info color.Color
	chat color.Color
}

//CreateColorConfig(prompt color.Color, err color.Color, success color.Color, info color.Color)
func CreateColorConfig() colorConfig {
	promptColor := color.BgRGB(12, 12, 12).Add(color.FgHiCyan)
	errColor := color.BgRGB(12, 12, 12).Add(color.FgHiRed)
	successColor := color.BgRGB(12, 12, 12).Add(color.FgHiYellow)
	infoColor := color.BgRGB(12, 12, 12).Add(color.FgHiWhite)
	chatColor := color.BgRGB(12, 12, 12).AddRGB(122, 231, 95)

	newConfig := colorConfig{
		prompt: *promptColor,
		err: *errColor,
		success: *successColor,
		info: *infoColor,
		chat: *chatColor,
	}
	return newConfig
}

func filterInput(r rune) (rune, bool) {
	switch r {
	case readline.CharCtrlL:
		return '!', true
	}		
	return r, true
}

func CreateReadlineConfig(prompt string) readline.Config {
	newCfg := readline.Config{
		FuncFilterInputRune: filterInput,
		Prompt: prompt,
	}
	return newCfg
}

func (cfg *config)SetNickname() {
	//ask for new nickname
	nname := ""
	update := database.SetNicknameParams{Nickname: nname, ID: cfg.User.UserID}
	upUser, err := cfg.db.SetNickname(cfg.ctx, update)
	if err != nil {
		cfg.colorCon.err.Println("error updating nickname:", err)
	}
	cfg.User.UpdatedAt = time.Now()
	cfg.User.Nickname = upUser.Nickname
}