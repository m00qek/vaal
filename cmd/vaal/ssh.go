package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func remoteAddr(host ConfigHost) string {
	return *host.Server.User + "@" + *host.Server.Addr
}

func execute(command string, args []string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// TODO: Maybe (maybe!) do this in go instead off calling ssh
func SecureShell(host ConfigHost, dryRun bool, arguments ...string) error {
	args := append([]string{remoteAddr(host)}, arguments...)

	if dryRun {
		fmt.Println("ssh " + strings.Join(args[:], " "))
		return nil
	}

	return execute("ssh", args)
}

func SecureCopyToRemote(host ConfigHost, dryRun bool, from string, to string) error {
	remoteTo := remoteAddr(host) + ":/" + to

	if dryRun {
		fmt.Println("scp -O " + from + " " + remoteTo)
		return nil
	}

	fmt.Println("Sending '" + to + "' to server ...")
	return execute("scp", []string{"-O", from, remoteTo})
}
