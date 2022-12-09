// Package cmd /*
package cmd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"os"

	"github.com/spf13/cobra"
	"go.etcd.io/etcd/server/v3/etcdmain"
)

// ServerCmd represents the server command
var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(os.Args)
		go etcdmain.Main(checkArgs(os.Args))

		cli, err := clientv3.New(clientv3.Config{
			Endpoints: []string{"127.0.0.1:2379"},
		})
		if err != nil {
			log.Fatal(err)
		}
		defer cli.Close()
		kvc := clientv3.NewKV(cli)
		response, err := kvc.Put(context.Background(), "key", "value")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(response)

		value, err := kvc.Get(context.Background(), "key")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(value)
	},
}

func checkArgs(args []string) (targetArgs []string) {
	for _, arg := range args {
		if arg != "server" {
			targetArgs = append(targetArgs, arg)
		} else {
			continue
		}
	}
	return targetArgs
}
