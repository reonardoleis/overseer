package discord

import (
	"fmt"
	"strings"
)

type CommandInfo struct {
	argc      int
	isMaximum bool
	argNames  []string
}

var (
	commands = map[string]*CommandInfo{
		"join":           {0, false, []string{}},
		"audio":          {1, false, []string{"ID or alias"}},
		"favoritecreate": {2, false, []string{"ID", "alias"}},
		"favoritelist":   {0, false, []string{}},
		"randomaudios":   {1, false, []string{"number of audios"}},
		"chatgpt":        {800, true, []string{"prompt"}},
		"skip":           {0, false, []string{}},
		"help":           {0, false, []string{}},
		"leave":          {0, false, []string{}},
		"loop":           {0, false, []string{}},
		"chatgpttts":     {200, true, []string{"prompt"}},
		"image":          {50, true, []string{"prompt"}},
		"fncreate":       {5000, true, []string{"name", "code..."}},
		"fnrun":          {100, true, []string{"name, args..."}},
	}
)

func commandHelp() string {
	s := "```"
	for cmd, info := range commands {
		s += cmd
		for _, arg := range info.argNames {
			s += " <" + arg + ">"
		}

		if info.isMaximum {
			s += fmt.Sprintf(" (maximum %d)", info.argc)
		}

		s += "\n"
	}
	s += "```"

	return s
}

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
