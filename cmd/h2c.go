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
	"crypto/tls"
	"crypto/x509"
	"gola/internal/logger"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

// h2cCmd represents the h2c command
var h2cCmd = &cobra.Command{
	Use:     "h2c",
	Short:   "http2 client",
	Long:    `測試呼叫http2(https)的server`,
	Example: program("h2c"),
	Run: func(cmd *cobra.Command, args []string) {
		urlStr, err := cmd.Flags().GetString("url")
		if err != nil {
			logger.Error(err)
			return
		}

		certFile, err := cmd.Flags().GetString("cert")
		if err != nil {
			logger.Error(err)
			return
		}

		link, err := url.Parse(urlStr)
		if err != nil {
			logger.Danger("請提供有效的網址 --url, -u | %s", err.Error())
			return
		}

		cert, err := ioutil.ReadFile(certFile)
		if err != nil {
			logger.Error(err)
			return
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(cert)

		// Setup HTTPS client
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:            caCertPool,
					InsecureSkipVerify: true,
				},
			},
		}

		res, err := client.Get(link.String())
		if err != nil {
			logger.Error(err)
			return
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			logger.Warn(res.Status)
			return
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logger.Danger("Read Body Failed: %s", err.Error())
			return
		}

		logger.Success(string(body))
	},
}

func init() {
	rootCmd.AddCommand(h2cCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// h2cCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	h2cCmd.Flags().StringP("url", "u", "https://gola.local:30004", "要呼叫的網址")
	h2cCmd.Flags().StringP("cert", "c", "cert.pem", "帶入的憑證")
}
