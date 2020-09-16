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
	"gola/internal/schedule"

	"github.com/spf13/cobra"
)

// scheduleCmd represents the schedule command
var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "啟動背景排程",
	Long:  `啟動背景排程管理，自動記錄排程運行時間與狀態`,
	Run: func(cmd *cobra.Command, args []string) {
		bootstrap.SetRunMode(bootstrap.CommandMode)
		bootstrap.SetupGracefulSignal() // 設定優雅結束程序
		schedule.Run()
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scheduleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scheduleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
