package main

import (
	"fmt"

	pg "github.com/habx/pg-commands"
	"os"
)

func main() {
	dump, err := pg.NewDump(&pg.Postgres{
		Host:     "localhost",
		Port:     5432,
		DB:       "dev_example",
		Username: "example",
		Password: "example",
	})
	if err != nil {
		panic(err)
	}

	dump.EnableVerbose()

	// Default options
	//dumpExec := dump.Exec(pg.ExecOptions{})

	// Log to stdout
	// Note that any io.Writer could be assigned to StreamDestination. I'm just using stdout here for brevity.
	// But we could write to a file, a database, a RabbitMQ queue or whatever
	// We could even write to all of the above using io.MultiWriter(...)
	dumpExec := dump.Exec(pg.ExecOptions{StreamPrint: true, StreamDestination: os.Stdout})
	if dumpExec.Error != nil {
		fmt.Println(dumpExec.Error.Err)
		fmt.Println(dumpExec.Output)

	} else {
		fmt.Println("Dump success")
		fmt.Println(dumpExec.Output)
	}

	restore, err := pg.NewRestore(&pg.Postgres{
		Host:     "localhost",
		Port:     5432,
		DB:       "dev_example",
		Username: "example",
		Password: "example",
	})
	if err != nil {
		panic(err)
	}
	restoreExec := restore.Exec(dumpExec.File, pg.ExecOptions{StreamPrint: false})
	if restoreExec.Error != nil {
		fmt.Println(restoreExec.Error.Err)
		fmt.Println(restoreExec.Output)

	} else {
		fmt.Println("Restore success")
		fmt.Println(restoreExec.Output)

	}
}
