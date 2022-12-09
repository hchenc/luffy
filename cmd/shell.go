// Package cmd /*
package cmd

import (
	"fmt"
	"github.com/hchenc/luffy/pkg/connector"
	"github.com/hchenc/luffy/pkg/host"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"os"
)

var cfgFile string
var command string
var server *host.Config

// ShellCmd represents the shell command
var ShellCmd = &cobra.Command{
	Use:   "shell",
	Short: "run shell command on remote hosts",
	Run: func(cmd *cobra.Command, args []string) {
		content, err := os.ReadFile(cfgFile)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(content, &server)
		if err != nil {
			panic(err)
		}
		var errs []error

		for _, host := range server.Hosts {
			c, err := connector.NewConnection(connector.Config{
				Username: host.User,
				Password: host.Password,
				Address:  host.Address,
			})
			if err != nil {
				errs = append(errs, err)
			}
			if command != "" {
				server.Do = command
			}
			output, code, err := c.Exec(host, server.Do)
			if err != nil {
				errs = append(errs, err)
			}
			fmt.Println(output)
			fmt.Println(code)
		}
	},
}

func init() {
	ShellCmd.Flags().StringVar(&cfgFile, "config", "/Users/joey/.luffy.yaml", "config file (default is $HOME/.luffy.yaml)")
	ShellCmd.Flags().StringVar(&command, "cmd", "", "shell command")
}
