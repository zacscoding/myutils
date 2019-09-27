package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/shiena/ansicolor"
	"github.com/urfave/cli"
	"github.com/zacscoding/myutils/host"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
)

var (
	sshCommand = cli.Command{
		Action:   ShowSubCommand,
		Name:     "ssh",
		Usage:    "command for ssh [shell]",
		Category: "SSH COMMANDS",
		Subcommands: []cli.Command{
			{
				Name:      "shell",
				Usage:     "open remote shell",
				Action:    openRemoteShell,
				ArgsUsage: "[host name]",
			},
		},
	}
)

func openRemoteShell(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return errors.New("invalid arguments")
	}
	conn, err := createSSHClient(ctx.Args()[0])
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer conn.Close()

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

// createSSHClient create a ssh client given ssh server info.
func createSSHClient(hostName string) (*ssh.Client, error) {
	h, err := host.GetHost(app.db, hostName)
	if err != nil {
		return nil, err
	}

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

// connectRemote connect to remote server with console.
func connectRemote(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return errors.New("invalid arguments")
	}

	hostName := ctx.Args()[0]
	h, err := host.GetHost(app.db, hostName)
	if err != nil {
		return err
	}

	var auth ssh.AuthMethod
	if h.Password != "" {
		auth = ssh.Password(h.Password)
	} else {
		pemBytes, err := ioutil.ReadFile(h.KeyPath)
		if err != nil {
			return err
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
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdout = ansicolor.NewAnsiColorWriter(os.Stdout)
	session.Stderr = ansicolor.NewAnsiColorWriter(os.Stderr)
	in, _ := session.StdinPipe()

	modes := ssh.TerminalModes{
		ssh.ECHO:  0, // Disable echoing
		ssh.IGNCR: 1, // Ignore CR on input.
	}

	if err := session.RequestPty("vt100", 80, 40, modes); err != nil {
		log.Fatalf("request for pseudo terminal failed: %s", err)
	}

	// Start remote shell
	if err := session.Shell(); err != nil {
		log.Fatalf("failed to start shell: %s", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for {
			<-c
			os.Exit(0)
		}
	}()

	// accepting commands
	for {
		reader := bufio.NewReader(os.Stdin)
		str, _ := reader.ReadString('\n')
		_, _ = fmt.Fprint(in, str)
	}
	return nil
}
