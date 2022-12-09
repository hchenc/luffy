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
var server *host.Config

// ShellCmd represents the shell command
var ShellCmd = &cobra.Command{
	Use:   "shell",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
	ServerCmd.Flags().StringVar(&cfgFile, "config", "/Users/joey/.luffy.yaml", "config file (default is $HOME/.luffy.yaml)")
}
