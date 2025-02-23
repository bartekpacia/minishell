package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

// Processes being currently executed.
var cmds []*exec.Cmd = make([]*exec.Cmd, 0)

func main() {
	log.SetFlags(0)
	log.SetPrefix("msh: ")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		for sig := range sigChan {
			log.Println("Interrupt signal received")
			for _, cmd := range cmds {
				cmd.Process.Signal(sig)
			}
		}
	}()

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
		fmt.Print("> ")
		cmdline, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalln("scan", err)
		}
		execute(cmdline)
	}
}

func execute(cmdline string) {
	parts := strings.FieldsFunc(cmdline, splitter)

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
	defer func() {
		cmds = cmds[:0]
	}()

	if len(cmds) == 1 {
		// There is only a single program running. Connect its stdin, stdout, stderr to our own.
		cmd := cmds[0]

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		err := cmd.Start()
		if err != nil {
			log.Printf("failed to start single command %s: %v\n", cmd.Path, err)
			return
		}
		err = cmd.Wait() // Correctly assign the error from cmd.Wait()
		if err != nil {
			log.Printf("failed to wait for single command %s: %v\n", cmd.Path, err)
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
			log.Printf("failed to start command %s: %v\n", cmd.Path, err)
		}
	}

	// Wait for all commands to finish
	for _, cmd := range cmds {
		err := cmd.Wait()
		if err != nil {
			log.Printf("failed to wait for command %s: %v\n", cmd.Path, err)
		}
	}
}

func splitter(r rune) bool {
	return r == '|' || r == '>'
}
