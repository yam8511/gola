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
	"gola/internal/logger"

	"github.com/spf13/cobra"
)

// greetClientCmd represents the greetClient command
var greetClientCmd = &cobra.Command{
	Use:     "client",
	Short:   "Greet Client",
	Long:    `Greet Client 用來呼叫RPC`,
	Example: program("grpc", "greet", "client"),
	Run: func(cmd *cobra.Command, args []string) {
		name := "World"
		if len(args) > 0 {
			name = args[0]
		}
		_, err := greet.Client(name)
		if err != nil {
			logger.Warn(err.Error())
		}
	},
}

// gorpcGreetClientCmd represents the greetClient command
var gorpcGreetClientCmd = &cobra.Command{
	Use:     "client",
	Short:   "Greet Client",
	Long:    `Greet Client 用來呼叫RPC`,
	Example: program("gorpc", "greet", "client"),
	Run: func(cmd *cobra.Command, args []string) {
		name := "World"
		if len(args) > 0 {
			name = args[0]
		}
		_, err := gogreet.Client(name)
		if err != nil {
			logger.Warn(err.Error())
		}
	},
}

func init() {
	greetCmd.AddCommand(greetClientCmd)
	goprcGreetCmd.AddCommand(gorpcGreetClientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// greetClientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// greetClientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
