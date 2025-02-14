package pgcommands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	// PGRestoreCmd is the path to the `pg_restore` executable
	psqlCmd     = "psql"
	psqlStdOpts = []string{"--echo-all"}
)

type Psql struct {
	*Postgres
	// Verbose mode
	Verbose bool
	// Path: setup path for source restore
	Path string
	// Extra pg_dump options
	// e.g []string{"--inserts"}
	Options []string
}

func NewPsql(pg *Postgres) (*Psql, error) {
	if !CommandExist(psqlCmd) {
		return nil, &ErrCommandNotFound{Command: psqlCmd}
	}

	return &Psql{Options: psqlStdOpts, Postgres: pg}, nil
}

// Exec `pg_restore` of the specified database, and restore from a gzip compressed tarball archive.
func (x *Psql) Exec(filename string, opts ExecOptions) Result {
	result := Result{}
	options := append(x.psqlOptions(), fmt.Sprintf("> %s%s", x.Path, filename))
	result.FullCommand = strings.Join(options, " ")
	cmd := exec.Command(psqlCmd, options...)
	cmd.Env = append(os.Environ(), x.EnvPassword)
	stdErr, _ := cmd.StderrPipe()
	stdIn, _ := cmd.StdinPipe()
	go streamOutput(stdErr, stdIn, opts, &result)
	err := cmd.Start()
	if err != nil {
		result.Error = &ResultError{Err: err, CmdOutput: result.Output}
	}
	err = cmd.Wait()
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		result.Error = &ResultError{Err: err, ExitCode: exitError.ExitCode(), CmdOutput: result.Output}
	}

	return result
}

func (x *Psql) ResetOptions() {
	x.Options = []string{}
}

func (x *Psql) EnableVerbose() {
	x.Verbose = true
}

func (x *Psql) SetPath(path string) {
	x.Path = path
}

func (x *Psql) psqlOptions() []string {
	options := x.Options
	options = append(options, x.Postgres.Parse()...)

	if x.Verbose {
		options = append(options, "-v")
	}

	return options
}
