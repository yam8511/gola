package console

import (
	consoleDemo "gola/app/console/demo"
)

// 註冊指令
var registeredCommand = map[string]Command{
	"demo": {"範例指令", consoleDemo.Run},
}

// GetCommands 取所有註冊指令
func GetCommands() map[string]Command {
	return registeredCommand
}

// GetCommand 取指定的註冊指令
func GetCommand(name string) *Command {
	cmd, ok := registeredCommand[name]
	if !ok {
		return nil
	}
	return &cmd
}

// Command 執行命令
type Command struct {
	// 指令描述
	Description string
	// 指令執行func
	Run func() error
}
