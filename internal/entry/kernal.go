package entry

import (
	"fmt"
	"gola/app/console"
	"gola/internal/bootstrap"
	"gola/internal/logger"
	"gola/internal/schedule"
	"gola/internal/server"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
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

	📌  舉例： APP_ENV=local APP_SITE=default ./gola server
	📌  舉例： APP_ENV=local APP_SITE=default ./gola schedule
	📌  舉例： APP_ENV=local APP_SITE=default ./gola run %s

`, builder.String(), commandName)

	if len(extraMessage) > 0 {
		fmt.Println(extraMessage...)
	}

	os.Exit(exitCode)
}

// Run 執行CronJob的 Command Line
func Run(payload ...func()) {
	log.Println(color.HiCyanString("⚙  APP_ROOT: %s", bootstrap.GetAppRoot()))
	log.Println(color.HiCyanString("⚙  APP_ENV: %s", bootstrap.GetAppEnv()))
	log.Println(color.HiCyanString("⚙  APP_SITE: %s", bootstrap.GetAppSite()))

	args := os.Args
	if len(args) < 2 {
		usage(0)
		return
	}

	// 載入設定檔
	bootstrap.LoadConfig()

	// 設定優雅結束程序
	bootstrap.SetupGracefulSignal()

	for _, fn := range payload {
		fn()
	}

	mainCmd := args[1]
	switch mainCmd {
	case "server":
		bootstrap.SetRunMode(bootstrap.ServerMode)
		server.Run()
	case "schedule":
		bootstrap.SetRunMode(bootstrap.CommandMode)
		schedule.Run()
	case "run":
		bootstrap.SetRunMode(bootstrap.CommandMode)
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
			<-bootstrap.GracefulDown()
			logger.Warn(`🚦  收到第一次訊號囉，若再收到一次，將會強制結束 🚦`)
			<-bootstrap.WaitOnceSignal()
			logger.Danger(`🚦  收到第二次訊號，強制結束 🚦`)
			os.Exit(2)
		}()

		err := cmd.Run()
		if err != nil {
			logger.Danger(fmt.Sprintf("指令[%s] (%s) 運行時，發生錯誤！ ---> %s\n", commandName, cmd.Description, err.Error()))
			os.Exit(1)
		}
		logger.Success(fmt.Sprintf("背景[%s] (%s) 運行結束\n", commandName, cmd.Description))

	default:
		logger.Warn(fmt.Sprintf("Unknown Command : %s", strings.Join(args, " ")))
		usage(1)
	}
}
