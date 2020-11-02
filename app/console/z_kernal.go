package console

import (
	consoleDemo "gola/app/console/demo"

	"github.com/spf13/cobra"
)

// 註冊指令
var registeredCommand = map[string]*cobra.Command{
	"demo": {Use: "demo", Short: "範例指令", RunE: consoleDemo.Run},
}

// GetCommands 取所有註冊指令
func GetCommands() map[string]*cobra.Command {
	return registeredCommand
}

// GetCommand 取指定的註冊指令
func GetCommand(name string) *cobra.Command {
	cmd, ok := registeredCommand[name]
	if !ok {
		return nil
	}
	return cmd
}
