package discord

import "strings"

type ArgumentInfo struct {
	argc int
}

var (
	commandsArgc = map[string]*ArgumentInfo{
		"join":           {0},
		"audio":          {1},
		"createfavorite": {1},
		"favoritelist":   {0},
		"randomaudios":   {1},
	}
)

func parseArguments(message string) (command string, args []string) {
	s := strings.Split(message, " ")

	s[0] = strings.Replace(s[0], "!", "", 1)
	return s[0], s[1:]
}

func (arg *ArgumentInfo) validateArguments(args []string) (expected int, ok bool) {
	return arg.argc, len(args) == arg.argc
}
