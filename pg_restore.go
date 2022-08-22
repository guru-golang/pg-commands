package pg_commands

import (
	"fmt"
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
	// verbose mode
	verbose bool
	// role: do SET ROLE before restore
	role string
	// path: setup path for source restore
	path string
	// format: input file format (custom, directory, tar, plain text (default))
	format string
	// Extra pg_dump options
	// e.g []string{"--inserts"}
	options []string
	// schemas: list of database schema
	schemas []string
}

func NewRestore(pg *Postgres) *Restore {
	return &Restore{options: RestoreStdOpts, Postgres: pg, schemas: []string{"public"}}
}

// Exec `pg_restore` of the specified database, and restore from a gzip compressed tarball archive.
func (x *Restore) Exec(filename string, opts ExecOptions) Result {
	result := Result{}
	options := append(x.restoreOptions(), fmt.Sprintf("%s%s", x.path, filename))
	result.FullCommand = strings.Join(options, " ")
	cmd := exec.Command(RestoreCmd, options...)
	stderrIn, _ := cmd.StderrPipe()
	go func() {
		result.Output = streamExecOutput(stderrIn, opts)
	}()

	if err := cmd.Start(); err != nil {
		result.Error = &ResultError{Err: err, CmdOutput: result.Output}
	}
	if err := cmd.Wait(); err != nil {
		result.Error = &ResultError{Err: err, CmdOutput: result.Output}
	}

	return result
}

func (x *Restore) SetVerbose(verbose bool) {
	x.verbose = verbose
}
func (x *Restore) GetVerbose() bool {
	return x.verbose
}

func (x *Restore) SetPath(path string) {
	x.path = path
}
func (x *Restore) GetPath() string {
	return x.path
}

func (x *Restore) SetFormat(f string) {
	x.format = f
}
func (x *Restore) GetFormat() string {
	return x.format
}

func (x *Restore) SetOptions(o []string) {
	x.options = o
}
func (x *Restore) GetOptions() []string {
	return x.options
}

func (x *Restore) SetSchemas(schemas []string) {
	x.schemas = schemas
}
func (x *Restore) GetSchemas() []string {
	return x.schemas
}

func (x *Restore) restoreOptions() []string {
	options := x.options
	options = append(options, x.Postgres.Parse()...)

	if x.format != "" {
		options = append(options, fmt.Sprintf(`--format=%v`, x.format))
	} else {
		options = append(options, fmt.Sprintf(`--format=%v`, RestoreDefaultFormat))
	}
	if x.role != "" {
		options = append(options, fmt.Sprintf(`--role=%v`, x.role))
	} else if x.DB != "" {
		x.role = x.DB
		options = append(options, fmt.Sprintf(`--role=%v`, x.DB))
	}

	if x.verbose {
		options = append(options, "-v")
	}
	for _, schema := range x.schemas {
		options = append(options, "--schema="+schema)
	}

	return options
}
