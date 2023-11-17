package discord

import "strings"

type CommandInfo struct {
	argc      int
	isMaximum bool
}

var (
	commands = map[string]*CommandInfo{
		"join":           {0, false},
		"audio":          {1, false},
		"favoritecreate": {2, false},
		"favoritelist":   {0, false},
		"randomaudios":   {1, false},
		"chatgpt":        {200, true},
		"skip":           {0, false},
	}
)

func getCommandInfo(command string) *CommandInfo {
	return commands[command]
}

func parseArguments(message string) (command string, args []string) {
	s := strings.Split(message, " ")

	s[0] = strings.Replace(s[0], "!", "", 1)
	return s[0], s[1:]
}

func (arg *CommandInfo) validateArguments(args []string) (expected int, ok bool) {
	return arg.argc, (len(args) == arg.argc || (arg.isMaximum && len(args) <= arg.argc))
}
