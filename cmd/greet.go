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
	gogreet "gola/gorpc/greet"
	"gola/grpc/greet"

	"github.com/spf13/cobra"
)

// greetCmd represents the greet command
var greetCmd = &cobra.Command{
	Use:     "greet",
	Short:   "Greet服務",
	Long:    `啟動Greet服務的RPC Server`,
	Example: program("grpc", "greet"),
	Run: func(cmd *cobra.Command, args []string) {
		greet.Server()
	},
}

// goprcGreetCmd represents the greet command
var goprcGreetCmd = &cobra.Command{
	Use:     "greet",
	Short:   "Greet服務",
	Long:    `啟動Greet服務的RPC Server`,
	Example: program("gorpc", "greet"),
	Run: func(cmd *cobra.Command, args []string) {
		gogreet.Server()
	},
}

func init() {
	grpcCmd.AddCommand(greetCmd)
	gorpcCmd.AddCommand(goprcGreetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// greetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// greetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
