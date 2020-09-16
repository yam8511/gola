/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"gola/app/console"
	"gola/internal/bootstrap"
	"gola/internal/logger"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gola",
	Short: "GoLa指令",
	Long: `GoLa指令
	啟動伺服器，提供「狼人殺」、「犯人在跳舞」...等遊戲服務
	亦可執行單一指令，或者執行背景排程功能
	`,
	Example: usage(),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	logger.Info(fmt.Sprintf("⚙  APP_ROOT: %s", bootstrap.GetAppRoot()))
	logger.Info(fmt.Sprintf("⚙  APP_ENV: %s", bootstrap.GetAppEnv()))
	logger.Info(fmt.Sprintf("⚙  APP_SITE: %s", bootstrap.GetAppSite()))

	bootstrap.LoadConfig()
}

func program(args ...string) string {
	program := os.Args[0]
	if strings.HasPrefix(program, "/var/") {
		program = "go run ."
	}

	if len(args) > 0 {
		program += " " + strings.Join(args, " ")
	}

	return program
}

func usage(extraMessage ...interface{}) string {
	commands := console.GetCommands()

	builder := new(strings.Builder)
	builder.WriteString("⚙  可執行的指令 (command name)")
	builder.WriteRune('\n')
	commandName := "<none>"
	for cmd := range commands {
		command := commands[cmd]
		commandName = cmd
		builder.WriteString(fmt.Sprintf("		✏ %s %s\n", cmd, command.Short))
	}

	program := program()

	return fmt.Sprintf(`
	📖 環境說明 📖

	需傳入以下環境變數：

	⚙  APP_ENV : 專案環境
		✏ default	預設值
		✏ docker	容器開發
		✏ local		本機開發
		✏ prod		正式

	⚙  APP_SITE : 專案端口
		✏ default	預設值

	--------------

	📖 指令說明 📖

	⚙  主要指令
		✏ server   運行伺服器
		✏ schedule 運行背景排程
		✏ run [command name] 執行指定命令

	%s

	📌  舉例： APP_ENV=local APP_SITE=default %s server
	📌  舉例： APP_ENV=local APP_SITE=default %s schedule
	📌  舉例： APP_ENV=local APP_SITE=default %s run %s
`, builder.String(),
		program,
		program,
		program,
		commandName,
	) + fmt.Sprintln(extraMessage...)
}
