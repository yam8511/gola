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
	"gola/app/common/def"
	"gola/internal/logger"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// traceCmd represents the trace command
var traceCmd = &cobra.Command{
	Use:   "trace",
	Short: "追蹤請求",
	Long:  `提供一組URL，會追蹤請求連線的DNS、TCP、Server Processing、Content Transfer的資訊`,
	Run: func(cmd *cobra.Command, args []string) {
		m, _ := cmd.Flags().GetString("method")
		logger.Info(fmt.Sprintf("curl -X %s %s", m, args))

		req, err := http.NewRequest(m, args[0], nil)
		if err != nil {
			logger.Error(err)
			return
		}

		req, trace := def.WithHttpTrace(req)
		trace.Begin = time.Now()
		defer func() {
			trace.End = time.Now()
			logger.Success(fmt.Sprintf("%+v", trace.Calculate()))
		}()

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			logger.Error(err)
			return
		}
		defer func() {
			trace.BodyCloseStart = time.Now()
			res.Body.Close()
			trace.BodyCloseDone = time.Now()
		}()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logger.Error(err)
			return
		}
		trace.TransferDone = time.Now()

		trace.ResultParseStart = time.Now()
		logger.Info(string(body))
		trace.ResultParseDone = time.Now()
	},
}

func init() {
	rootCmd.AddCommand(traceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// traceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// traceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	traceCmd.Flags().StringP("method", "X", "GET", "請求方式")
}
