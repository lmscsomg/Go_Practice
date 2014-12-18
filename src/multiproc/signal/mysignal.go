package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	go func() {
		time.Sleep(5 * time.Second)
		sendSigDemo()
	}()

	mySigalDemo()
}
func sendSigDemo() {
	cmds := []*exec.Cmd{
		exec.Command("ps", "aux"),
		exec.Command("grep", "mysignal"),
		exec.Command("grep", "-v", "grep"),
		exec.Command("grep", "-v", "go run"),
		exec.Command("awk", "{print $2}"),
	}
	output, err := runCmd(cmds)

	if err != nil {
		fmt.Printf("Command Execution Error: %s\n", err)
		return
	}

	pids, err := getPids(output)
	if err != nil {
		fmt.Printf("PID Parsing Error: %s\n", err)
		return
	}
	fmt.Printf("Target PID(s):\n%v\n", pids)
	for _, pid := range pids {
		proc, err := os.FindProcess(pid)
		if err != nil {
			fmt.Printf("Process Finding Error: %s\n", err)
			return
		}
		sig := syscall.SIGQUIT
		fmt.Printf("Send signal '%s' to the process (pid=%d)...\n", sig, pid)
		err = proc.Signal(sig)
		if err != nil {
			fmt.Printf("Signal Sending Error: %s\n", err)
			return
		}
	}
}

func getPids(strs []string) ([]int, error) {
	pids := make([]int, 0)
	for _, str := range strs {
		pid, err := strconv.Atoi(strings.TrimSpace(str))
		if err != nil {
			return nil, err
		}
		pids = append(pids, pid)
	}
	return pids, nil
}

func runCmd(cmds []*exec.Cmd) ([]string, error) {
	if cmds == nil || len(cmds) == 0 {
		return nil, errors.New("The cmd slice is invalid!")
	}
	first := true

	var output []byte
	var err error

	for _, cmd := range cmds {
		fmt.Printf("Run command: %v ... \n", getCmdPlaintext(cmd))
		if !first {
			var stdinBuf bytes.Buffer
			stdinBuf.Write(output)
			cmd.Stdin = &stdinBuf
		}

		var stdoutBuf bytes.Buffer
		cmd.Stdout = &stdoutBuf

		if err = cmd.Start(); err != nil {
			return nil, getError(err, cmd)
		}

		if err = cmd.Wait(); err != nil {
			return nil, getError(err, cmd)
		}
		output = stdoutBuf.Bytes()

		if first {
			first = false
		}
	}

	lines := make([]string, 0)
	var outputBuf bytes.Buffer
	outputBuf.Write(output)

	for {
		line, err := outputBuf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, getError(err, nil)
			}
		}
		lines = append(lines, string(line))
	}
	return lines, nil
}

func getCmdPlaintext(cmd *exec.Cmd) string {
	var buf bytes.Buffer
	buf.WriteString(cmd.Path)

	for _, arg := range cmd.Args[1:] {
		buf.WriteRune(' ')
		buf.WriteString(arg)
	}
	return buf.String()
}
func mySigalDemo() {
	sigRev1 := make(chan os.Signal, 1)
	sig1 := []os.Signal{syscall.SIGINT, syscall.SIGQUIT}
	// display the signals we set
	fmt.Printf("Set notification for %s ...[sigRev1]\n", sig1)
	signal.Notify(sigRev1, sig1...)

	sigRev2 := make(chan os.Signal, 1)
	sig2 := []os.Signal{syscall.SIGQUIT}
	fmt.Printf("Set notification for %s ...[sigRev1]\n", sig2)
	signal.Notify(sigRev2, sig2...)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for sig := range sigRev1 {
			fmt.Printf("Received a signal [ %s ] from sigRev1\n", sig)
		}
		fmt.Printf("End! [sigRev1]")
		wg.Done()
	}()

	go func() {
		for sig := range sigRev2 {
			fmt.Printf("Received a signal [ %s ] from sigRev2\n", sig)
		}
		fmt.Printf("End! [sigRev2]\n")
		wg.Done()
	}()

	fmt.Printf("\nWait for 5 seconds\n")
	time.Sleep(5 * time.Second)

	fmt.Printf("Stop notification!\n")
	signal.Stop(sigRev1)
	close(sigRev1)
	fmt.Printf("done. [sigRev1]\n")

	wg.Wait()
}

func getError(err error, cmd *exec.Cmd, extraInfo ...string) error {
	var errMsg string
	if cmd != nil {
		errMsg = fmt.Sprintf("%s  [%s %v]", err, (*cmd).Path, (*cmd).Args)
	} else {
		errMsg = fmt.Sprintf("%s", err)
	}
	if len(extraInfo) > 0 {
		errMsg = fmt.Sprintf("%s (%v)", errMsg, extraInfo)
	}
	return errors.New(errMsg)
}
