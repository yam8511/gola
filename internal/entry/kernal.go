package entry

import (
	"fmt"
	"gola/app/console"
	"gola/internal/bootstrap"
	"gola/internal/schedule"
	"gola/internal/server"
	"os"
	"strings"
)

func usage(exitCode int, extraMessage ...interface{}) {
	commands := console.GetCommands()

	builder := new(strings.Builder)
	builder.WriteString("⚙  可執行的指令 (command name)")
	builder.WriteRune('\n')
	commandName := "<none>"
	for cmd := range commands {
		command := commands[cmd]
		commandName = cmd
		builder.WriteString(fmt.Sprintf("		✏ %s %s\n", cmd, command.Description))
	}

	fmt.Printf(`
	📖 環境說明 📖

	需傳入以下環境變數：

	⚙  APP_ENV : 專案環境
		✏ docker 容器開發
		✏ local 本機開發
		✏ qatest 測試
		✏ prod 正式

	📌  舉例： APP_ENV=local ./app

	--------------

	📖 指令說明 📖

	⚙  主要指令
		✏ server   運行伺服器
		✏ schedule 運行背景排程
		✏ run [command name] 執行指定命令

	%s

	📌  舉例： APP_ENV=local ./app server
	📌  舉例： APP_ENV=local ./app schedule
	📌  舉例： APP_ENV=local ./app run %s

`, builder.String(), commandName)

	if len(extraMessage) > 0 {
		fmt.Println(extraMessage...)
	}

	os.Exit(exitCode)
}

// Run 執行CronJob的 Command Line
func Run() {
	if bootstrap.GetAppEnv() == "" {
		usage(0)
	} else {
		bootstrap.WriteLog("INFO", fmt.Sprintf("⚙  APP_ENV: %s", bootstrap.GetAppEnv()))
	}
	args := os.Args
	if len(args) < 2 {
		usage(0)
		return
	}

	// 設定優雅結束程序
	bootstrap.SetupGracefulSignal()

	mainCmd := args[1]
	switch mainCmd {
	case "server":
		server.Run()
	case "schedule":
		schedule.Run()
	case "run":
		if len(args) < 3 {
			usage(1, "請輸入欲執行命令")
			return
		}

		commandName := args[2]

		cmd := console.GetCommand(commandName)

		if cmd == nil {
			usage(1, fmt.Sprintf("命令尚未註冊: %s", commandName))
		}

		go func() {
			<-bootstrap.WaitOnceSignal()
			bootstrap.WriteLog("WARNING", `🚦  收到訊號囉，再收到一次，強制結束~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ 🚦`)
			<-bootstrap.WaitOnceSignal()
			bootstrap.WriteLog("WARNING", `🚦  收到第二次訊號，強制結束~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ 🚦`)
			os.Exit(2)
		}()

		err := cmd.Run()
		if err != nil {
			bootstrap.WriteLog("ERROR", fmt.Sprintf("背景[%s] (%s) 運行時，發生錯誤！ ---> %s\n", commandName, cmd.Description, err.Error()))
			os.Exit(1)
		}
		bootstrap.WriteLog("INFO", fmt.Sprintf("背景[%s] (%s) 運行結束\n", commandName, cmd.Description))

	default:
		bootstrap.WriteLog("WARNING", fmt.Sprintf("Unknown Command : %s", strings.Join(args, " ")))
		usage(1)
	}
}
