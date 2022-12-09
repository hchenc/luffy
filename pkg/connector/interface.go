package connector

import (
	"github.com/hchenc/luffy/pkg/host"
	"io"
)

type Connection interface {
	Exec(host host.Host, cmd string) (string, int, error)
	PipeExec(host host.Host, cmd string, stdin io.Reader, stdout io.Writer, stderr io.Writer) (code int, err error)
	Fetch(host host.Host, local, remote string) error
	Scp(host host.Host, local, remote string) error
	Close()
}
