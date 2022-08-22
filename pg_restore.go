package pg_commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	// RestoreCmd is the path to the `pg_restore` executable
	RestoreCmd           = "pg_restore"
	RestoreStdOpts       = []string{"--exit-on-error"}
	RestoreDefaultFormat = "p" // p  c  d  t
)

type Restore struct {
	*Postgres
	// Verbose mode
	Verbose bool
	// Role: do SET ROLE before restore
	Role string
	// Path: setup path for source restore
	Path string
	// Format: input file format (custom, directory, tar, plain text (default))
	Format string
	// Extra pg_dump options
	// e.g []string{"--inserts"}
	Options []string
	// Schemas: list of database schema
	Schemas []string
}

func NewRestore(pg *Postgres) *Restore {
	return &Restore{Options: RestoreStdOpts, Postgres: pg, Schemas: []string{"public"}}
}

// Exec `pg_restore` of the specified database, and restore from a gzip compressed tarball archive.
func (x *Restore) Exec(filename string, opts ExecOptions) Result {
	result := Result{}
	options := append(x.restoreOptions(), fmt.Sprintf("%s%s", x.Path, filename))
	result.FullCommand = strings.Join(options, " ")
	cmd := exec.Command(RestoreCmd, options...)

	cmd.Env = append(os.Environ(), x.EnvPassword)
	stderrIn, _ := cmd.StderrPipe()
	go func() {
		result.Output = streamExecOutput(stderrIn, opts)
	}()
	err := cmd.Start()
	if err != nil {
		result.Error = &ResultError{Err: err, CmdOutput: result.Output}
	}
	err = cmd.Wait()
	if exitError, ok := err.(*exec.ExitError); ok {
		result.Error = &ResultError{Err: err, ExitCode: exitError.ExitCode(), CmdOutput: result.Output}
	}

	return result
}

func (x *Restore) ResetOptions() {
	x.Options = []string{}
}

func (x *Restore) EnableVerbose() {
	x.Verbose = true
}

func (x *Restore) SetPath(path string) {
	x.Path = path
}

func (x *Restore) SetupFormat(f string) {
	x.Format = f
}

func (x *Restore) SetSchemas(schemas []string) {
	x.Schemas = schemas
}

func (x *Restore) restoreOptions() []string {
	options := x.Options
	options = append(options, x.Postgres.Parse()...)

	if x.Format != "" {
		options = append(options, fmt.Sprintf(`-F %v`, x.Format))
	} else {
		options = append(options, fmt.Sprintf(`-F %v`, RestoreDefaultFormat))
	}
	if x.Role != "" {
		options = append(options, fmt.Sprintf(`--role=%v`, x.Role))
	} else if x.DB != "" {
		x.Role = x.DB
		options = append(options, fmt.Sprintf(`--role=%v`, x.DB))
	}

	if x.Verbose {
		options = append(options, "-v")
	}
	for _, schema := range x.Schemas {
		options = append(options, "--schema="+schema)
	}

	return options
}
