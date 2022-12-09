package connector

import (
	"bufio"
	"context"
	"fmt"
	"github.com/hchenc/luffy/pkg/host"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Config struct {
	Username       string
	Password       string
	Address        string
	Port           int
	PrivateKey     string
	PrivateKeyPath string
	Timeout        time.Duration
}

type connection struct {
	mu        sync.Mutex
	sshclient *ssh.Client
	ctx       context.Context
	cancel    context.CancelFunc
}

func (c *connection) session() (*ssh.Session, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.sshclient == nil {
		return nil, errors.New("connection closed")
	}

	sess, err := c.sshclient.NewSession()
	if err != nil {
		return nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	err = sess.RequestPty("xterm", 100, 50, modes)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (c *connection) Exec(host host.Host, cmd string) (string, int, error) {
	sess, err := c.session()
	if err != nil {
		return "", 1, errors.Wrap(err, "failed to get SSH session")
	}
	defer sess.Close()

	exitCode := 0

	in, _ := sess.StdinPipe()
	out, _ := sess.StdoutPipe()

	err = sess.Start(strings.TrimSpace(cmd))
	if err != nil {
		exitCode = -1
		if exitErr, ok := err.(*ssh.ExitError); ok {
			exitCode = exitErr.ExitStatus()
		}
		return "", exitCode, err
	}

	var (
		output []byte
		line   = ""
		r      = bufio.NewReader(out)
	)

	for {
		b, err := r.ReadByte()
		if err != nil {
			break
		}

		output = append(output, b)

		if b == byte('\n') {
			line = ""
			continue
		}

		line += string(b)

		if (strings.HasPrefix(line, "[sudo] password for ") || strings.HasPrefix(line, "Password")) && strings.HasSuffix(line, ": ") {
			_, err = in.Write([]byte(host.Password + "\n"))
			if err != nil {
				break
			}
		}
	}
	err = sess.Wait()
	if err != nil {
		exitCode = -1
		if exitErr, ok := err.(*ssh.ExitError); ok {
			exitCode = exitErr.ExitStatus()
		}
	}
	outStr := strings.TrimPrefix(string(output), fmt.Sprintf("[sudo] password for %s:", host.User))

	// preserve original error
	return strings.TrimSpace(outStr), exitCode, errors.Wrapf(err, "Failed to exec command: %s \n%s", cmd, strings.TrimSpace(outStr))
}

func (c *connection) PipeExec(host host.Host, cmd string, stdin io.Reader, stdout io.Writer, stderr io.Writer) (code int, err error) {
	//TODO implement me
	panic("implement me")
}

func (c *connection) Fetch(host host.Host, local, remote string) error {
	//TODO implement me
	panic("implement me")
}

func (c *connection) Scp(host host.Host, local, remote string) error {
	//TODO implement me
	panic("implement me")
}

func (c *connection) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.sshclient == nil {
		return
	}
	c.cancel()

	if c.sshclient != nil {
		c.sshclient.Close()
		c.sshclient = nil
	}
}

func NewConnection(config Config) (Connection, error) {
	config, err := validateOptions(config)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to validate ssh connection parameters")
	}

	authMethods := make([]ssh.AuthMethod, 0)

	if config.Password != "" {
		authMethods = append(authMethods, ssh.Password(config.Password))
	}

	if config.PrivateKey != "" {
		signer, parseErr := ssh.ParsePrivateKey([]byte(config.PrivateKey))
		if parseErr != nil {
			return nil, errors.Wrap(parseErr, "The given SSH key could not be parsed")
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	sshConfig := &ssh.ClientConfig{
		User:            config.Username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         config.Timeout,
	}

	targetHost := config.Address
	targetPort := strconv.Itoa(config.Port)

	endpoint := net.JoinHostPort(targetHost, targetPort)

	client, err := ssh.Dial("tcp", endpoint, sshConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "could not establish connection to %s", endpoint)
	}

	ctx, cancel := context.WithCancel(context.Background())

	sshConn := &connection{
		ctx:       ctx,
		sshclient: client,
		cancel:    cancel,
	}

	return sshConn, nil
}

func validateOptions(config Config) (Config, error) {
	if len(config.Username) == 0 {
		return config, errors.New("No username specified for SSH connection")
	}

	if len(config.Address) == 0 {
		return config, errors.New("No address specified for SSH connection")
	}

	if len(config.Password) == 0 && len(config.PrivateKey) == 0 {
		return config, errors.New("Must specify at least one of password, private key")
	}

	if len(config.PrivateKey) == 0 && len(config.PrivateKeyPath) > 0 {
		content, err := os.ReadFile(config.PrivateKeyPath)
		if err != nil {
			return config, errors.Wrapf(err, "Failed to read keyfile %q", config.PrivateKeyPath)
		}

		config.PrivateKey = string(content)
		config.PrivateKeyPath = ""
	}

	if config.Port <= 0 {
		config.Port = 22
	}

	if config.Timeout == 0 {
		config.Timeout = 15 * time.Second
	}

	return config, nil
}
