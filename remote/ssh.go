package remote

import (
	"bytes"
	"github.com/shiena/ansicolor"
	"github.com/zacscoding/myutils/types"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
)

type HostCmdResult struct {
	Host    *types.Host
	Command string
	Result  *types.ExecuteResult
	Err     error
}

type CommandGenerator func(h *types.Host) string
type CommandHandler func(result HostCmdResult)

// CreateSSHClient create ssh client given a host
func CreateSSHClient(h *types.Host) (*ssh.Client, error) {
	var auth ssh.AuthMethod
	if h.Password != "" {
		auth = ssh.Password(h.Password)
	} else {
		pemBytes, err := ioutil.ReadFile(h.KeyPath)
		if err != nil {
			return nil, err
		}
		key, err := ssh.ParsePrivateKey(pemBytes)
		auth = ssh.PublicKeys(key)
	}

	config := &ssh.ClientConfig{
		User: h.User,
		Auth: []ssh.AuthMethod{
			auth,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := h.Address + ":" + strconv.Itoa(h.Port)
	return ssh.Dial("tcp", addr, config)
}

// OpenRemoteShell start to open remote shell
func OpenRemoteShell(conn *ssh.Client) error {
	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdin = os.Stdin
	session.Stdout = ansicolor.NewAnsiColorWriter(os.Stdout)
	session.Stderr = ansicolor.NewAnsiColorWriter(os.Stderr)

	// copy from http://talks.rodaine.com/gosf-ssh/present.slide#9
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,      // please print what I type
		ssh.ECHOCTL:       0,      // please don't print control chars
		ssh.TTY_OP_ISPEED: 115200, // baud in
		ssh.TTY_OP_OSPEED: 115200, // baud out
	}

	termFD := int(os.Stdin.Fd())

	width, height, err := terminal.GetSize(termFD)
	if err != nil {
		return err
	}

	termState, _ := terminal.MakeRaw(termFD)
	defer terminal.Restore(termFD, termState)

	err = session.RequestPty("xterm-256color", height, width, modes)
	if err != nil {
		return err
	}
	err = session.Shell()
	if err != nil {
		return err
	}
	err = session.Wait()
	if err != nil {
		return err
	}
	return nil
}

// executesCommand execute command to given hosts with go routines
func ExecutesCommand(hosts []*types.Host, commandGen CommandGenerator, handler CommandHandler) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(hosts))
	cmdResults := make(chan HostCmdResult)

	for _, h := range hosts {
		go func(h *types.Host, w *sync.WaitGroup, ch chan HostCmdResult) {
			command := commandGen(h)
			conn, err := CreateSSHClient(h)
			if err != nil {
				ch <- HostCmdResult{h, command, nil, err}
				w.Done()
				return
			}
			defer conn.Close()

			session, err := conn.NewSession()
			if err != nil {
				ch <- HostCmdResult{h, command, nil, err}
				w.Done()
				return
			}
			defer session.Close()

			var stdOut bytes.Buffer
			var stdErr bytes.Buffer
			session.Stdout = &stdOut
			session.Stderr = &stdErr

			err = session.Run(command)
			ch <- HostCmdResult{
				Host: h,
				Result: &types.ExecuteResult{
					Error:  err,
					StdOut: stdOut.String(),
					StdErr: stdErr.String(),
				},
				Err: nil,
			}
			w.Done()
		}(h, &waitGroup, cmdResults)
	}
	go func() {
		for result := range cmdResults {
			handler(result)
		}
	}()
	waitGroup.Wait()
}
