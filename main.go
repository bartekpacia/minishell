package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) == 1 {
		runRepl()
	} else if len(os.Args) == 2 {
		cmdline := os.Args[1]
		execute(cmdline)
	} else {
		log.Fatalln("invalid number of args")
	}
}

func runRepl() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		cmdline, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln("scan", err)
		}
		execute(cmdline)
	}
}

func execute(cmdline string) {
	parts := strings.FieldsFunc(cmdline, splitter)

	cmds := make([]*exec.Cmd, 0)
	for _, part := range parts {
		fields := strings.Fields(part)
		if len(fields) == 0 {
			return
		}
		cmdname := strings.Fields(part)[0]
		args := strings.Fields(part)[1:]

		newcmd := exec.Command(cmdname, args...)

		cmds = append(cmds, newcmd)
	}

	if len(cmds) == 1 {
		// There is only a single program running. Connect its stdin, stdout, stderr to our own.
		cmd := cmds[0]

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		err := cmd.Start()
		if err != nil {
			log.Fatalln(err)
		}
		err = cmd.Wait() // Correctly assign the error from cmd.Wait()
		if err != nil {
			log.Fatalln("failed to wait for command:", err)
		}

		return
	}

	// Connect stdout and stdins
	for i := 1; i < len(cmds); i++ {
		prevcmd := cmds[i-1]
		cmd := cmds[i]

		if i == 1 {
			prevcmd.Stdin = os.Stdin
		}

		cmd.Stdin, _ = prevcmd.StdoutPipe()

		if i == len(cmds)-1 {
			cmd.Stdout = os.Stdout
		}
	}

	// Start all commands
	for _, cmd := range cmds {
		err := cmd.Start()
		if err != nil {
			log.Fatalln(err)
		}
	}

	// Wait for all commands to finish
	for _, cmd := range cmds {
		err := cmd.Wait()
		if err != nil {
			log.Fatalln("failed to wait for command:", err)
		}
	}
}

func splitter(r rune) bool {
	return r == '|' || r == '>'
}
