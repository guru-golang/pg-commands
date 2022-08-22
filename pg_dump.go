package pg_commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	// DumpCmd is the path to the `pg_dump` executable
	DumpCmd           = "pg_dump"
	DumpStdOpts       = []string{}
	DumpDefaultFormat = "p" // p  c  d  t
)

// Dump is an `Exporter` interface that backs up a Postgres database via the `pg_dump` command.
type Dump struct {
	*Postgres
	// verbose mode
	verbose bool
	// path: setup path dump out
	path string
	// file: output file
	file string
	// format: output file format (custom, directory, tar, plain text (default))
	format string
	// Extra pg_dump x.FullOptions
	// e.g []string{"--inserts"}
	options []string

	IgnoreTableData []string
}

func NewDump(pg *Postgres) *Dump {
	return &Dump{options: DumpStdOpts, Postgres: pg}
}

// Exec `pg_dump` of the specified database, and creates a gzip compressed tarball archive.
func (x *Dump) Exec(opts ExecOptions) Result {
	result := Result{Mine: "application/sql"}
	result.File = x.GetFile()
	options := append(x.dumpOptions(), fmt.Sprintf(`--file=%s%v`, x.path, result.File))
	result.FullCommand = strings.Join(options, " ")
	cmd := exec.Command(DumpCmd, options...)
	cmd.Env = append(os.Environ(), x.EnvPassword)
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

func (x *Dump) SetVerbose(verbose bool) {
	x.verbose = verbose
}
func (x *Dump) GetVerbose() bool {
	return x.verbose
}

func (x *Dump) SetPath(path string) {
	x.path = path
}
func (x *Dump) GetPath() string {
	return x.path
}

func (x *Dump) SetFile(filename string) {
	x.file = filename
}
func (x *Dump) GetFile() string {
	if x.file == "" {
		// Use default file name
		x.file = x.newFileName()
	}
	return x.file
}

func (x *Dump) SetFormat(f string) {
	x.format = f
}
func (x *Dump) GetFormat() string {
	return x.format
}

func (x *Dump) SetOptions(o []string) {
	x.options = o
}
func (x *Dump) GetOptions() []string {
	return x.options
}

func (x *Dump) newFileName() string {
	return fmt.Sprintf(`%v_%v.sql`, x.DB, time.Now().Unix())
}

func (x *Dump) dumpOptions() []string {
	options := x.options
	options = append(options, x.Postgres.Parse()...)

	if x.format != "" {
		options = append(options, fmt.Sprintf(`--format=%v`, x.format))
	} else {
		options = append(options, fmt.Sprintf(`--format=%v`, DumpDefaultFormat))
	}
	if x.verbose {
		options = append(options, "-v")
	}
	if len(x.IgnoreTableData) > 0 {
		options = append(options, x.IgnoreTableDataToString()...)
	}

	return options
}

func (x *Dump) IgnoreTableDataToString() []string {
	var t []string
	for _, tables := range x.IgnoreTableData {
		t = append(t, "--exclude-table-data="+tables)
	}
	return t
}
