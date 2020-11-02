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
	"gola/internal/bootstrap"
	"gola/internal/server"
	defaultR "gola/router/default"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "啟動伺服器",
	Long:  `啟動伺服器`,
	Run: func(cmd *cobra.Command, args []string) {
		bootstrap.SetRunMode(bootstrap.ServerMode)
		bootstrap.SetupGracefulSignal() // 設定優雅結束程序

		conf := bootstrap.GetAppConf()
		switch conf.App.Site {
		case "admin":
			// 專屬admin的route
		case "member":
			// 專屬member的route
		default:
			// 檢查資料庫資料表
			// modelList := []model.IModel{
			// 	new(model.User),
			// }
			// model.SetupTable(modelList...)

			server.Run(defaultR.LoadRoutes)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	serverCmd.Flags().BoolP("port", "p", false, "自定義Port")
}
