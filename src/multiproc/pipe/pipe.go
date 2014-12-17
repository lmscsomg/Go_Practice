package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
)

/*
	1. cmd1输出管道
	2. cmd1开启
	3. 输出管道存入缓冲区

	4. cmd2输入管道
	5. 将缓冲区中得输出管道写入输入管道
	6. cmd2开启
	7. 关闭输入管道
	8. 等待并输出
*/

func main() {

	fmt.Println("Run command `ps aux | grep pipe`: ")
	cmd1 := exec.Command("ps", "aux")
	cmd2 := exec.Command("grep", "pipe")

	stdout1, err := cmd1.StdoutPipe()

	if err != nil {
		fmt.Printf("Error, the command cannot be started up %s\n", err)
	}
	if err := cmd1.Start(); err != nil {
		fmt.Printf("Error, the command cannot be started up %s\n", err)
	}
	output1 := bufio.NewReader(stdout1)

	stdin2, err := cmd2.StdinPipe()

	if err != nil {
		fmt.Printf("Error, the command cannot be started up %s\n", err)
	}
	//connect
	output1.WriteTo(stdin2)

	//output to the buffer
	var outputbuf2 bytes.Buffer
	cmd2.Stdout = &outputbuf2

	if err := cmd2.Start(); err != nil {
		fmt.Printf("Error, the command cannot be started up %s\n", err)
	}

	err = stdin2.Close()
	if err != nil {
		fmt.Printf("Error: Cannot close the stdio pipe %s\n", err)
	}

	if err := cmd2.Wait(); err != nil {
		fmt.Printf("Error: Can not wait for the command: %s\n", err)
	}
	fmt.Printf("%s\n", outputbuf2.Bytes())
}
