package connector

import (
	"fmt"
	"github.com/hchenc/luffy/pkg/host"
	"testing"
)

func TestSSH(t *testing.T) {
	c, _ := NewConnection(Config{
		Username: "root",
		Password: "test",
		Address:  "127.0.0.1",
		Port:     22,
	})
	stdout, errCode, err := c.Exec(host.Host{
		Address:  "127.0.0.1",
		User:     "root",
		Password: "test",
	}, "ls")
	fmt.Println(stdout)
	fmt.Println(errCode)
	fmt.Println(err)
}
