package dst_cli_windows

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type PythonProcess struct {
	cmd    *exec.Cmd
	stdin  *os.File
	stdout *os.File
}

func NewPythonProcess(command string, args ...string) (*PythonProcess, error) {
	cmd := exec.Command("cmd", "/C", command)

	// 创建用于输入的匿名管道
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	// 创建用于输出的匿名管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	p := &PythonProcess{
		cmd:    cmd,
		stdin:  stdin.(*os.File),
		stdout: stdout.(*os.File),
	}

	return p, nil
}

func (p *PythonProcess) Start() error {
	err := p.cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func (p *PythonProcess) Stop() error {
	err := p.stdin.Close()
	if err != nil {
		return err
	}

	err = p.stdout.Close()
	if err != nil {
		return err
	}

	err = p.cmd.Process.Signal(os.Interrupt)
	if err != nil {
		return err
	}

	err = p.cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (p *PythonProcess) SendInput(input string) string {
	p.stdin.WriteString(input + "\n")
	p.stdin.Sync()

	reader := bufio.NewReader(p.stdout)
	output, _ := reader.ReadString('\n')
	return strings.TrimSuffix(output, "\n")
}

func main() {
	pythonProcess, err := NewPythonProcess("python", "your_script.py")
	if err != nil {
		log.Fatal(err)
	}

	err = pythonProcess.Start()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := pythonProcess.Stop()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// 向子进程发送输入
	pythonProcess.SendInput("Hello")
	pythonProcess.SendInput("World")

	// 在此处可以执行其他操作

	// 等待主进程退出
	fmt.Println("Press Enter to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
