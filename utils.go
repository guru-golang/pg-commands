package pgcommands

import (
	"fmt"
	"io"
	"os/exec"
	"sync"
)

type ExecOptions struct {
	StreamPrint       bool
	StreamDestination io.Writer
}

func streamOutput(stdErr io.ReadCloser, stdIn io.WriteCloser, opts ExecOptions, result *Result) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		buf := make([]byte, 1024)

		for {
			n, err := stdErr.Read(buf)
			if err != nil {
				if err == io.EOF {
					return // End of output
				}
				fmt.Println("Error reading from stdout:", err)
				return
			}
			fmt.Print(string(buf[:n]))
		}
	}()

	wg.Wait()
}

func CommandExist(command string) bool {
	_, err := exec.LookPath(command)

	return err == nil
}
