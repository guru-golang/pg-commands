package pgcommands

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

type ExecOptions struct {
	StreamPrint       bool
	StreamDestination io.Writer
}

func streamExecOutput(out io.ReadCloser, options ExecOptions) (string, error) {
	output := ""
	reader := bufio.NewReader(out)
	for {
		line, err := reader.ReadString('\n')
		fmt.Println("streamExecOutput:", line)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return output, nil
			}

			return output, fmt.Errorf("error reading output: %w", err)
		}

		if options.StreamPrint {
			_, err = fmt.Fprint(options.StreamDestination, line)
			if err != nil {
				return output, fmt.Errorf("error writing output: %w", err)
			}
		}

		output += line
	}
}

func streamExecInput(out io.ReadCloser, options ExecOptions) (string, error) {
	output := ""
	reader := bufio.NewReader(out)
	for {
		line, err := reader.ReadString('\n')
		fmt.Println("streamExecOutput:", line)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return output, nil
			}

			return output, fmt.Errorf("error reading output: %w", err)
		}

		if options.StreamPrint {
			_, err = fmt.Fprint(options.StreamDestination, line)
			if err != nil {
				return output, fmt.Errorf("error writing output: %w", err)
			}
		}

		output += line
	}
}

func streamOutput(stderrIn io.ReadCloser, stderrOut io.ReadCloser, opts ExecOptions, result *Result) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		output, err := streamExecInput(stderrIn, opts)
		if err != nil {
			result.Error = &ResultError{Err: err, CmdOutput: output}
		}
		fmt.Println("streamExecInput", output)
		result.Output = output
	}()

	go func() {
		defer wg.Done()
		output, err := streamExecOutput(stderrOut, opts)
		if err != nil {
			result.Error = &ResultError{Err: err, CmdOutput: output}
		}
		fmt.Println("streamExecOutput", output)
		//result.Output = output
	}()

	wg.Wait()
}

func CommandExist(command string) bool {
	_, err := exec.LookPath(command)

	return err == nil
}
