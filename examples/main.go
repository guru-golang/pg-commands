package main

import (
	"fmt"
	pg "github.com/guru-golang/pg-commands"
)

func main() {
	dump := pg.NewDump(&pg.Postgres{
		Host:     "localhost",
		Port:     5432,
		DB:       "dev_example",
		Username: "example",
		Password: "example",
	})
	dumpExec := dump.Exec(pg.ExecOptions{StreamPrint: false})
	if dumpExec.Error != nil {
		fmt.Println(dumpExec.Error.Err)
		fmt.Println(dumpExec.Output)

	} else {
		fmt.Println("Dump success")
		fmt.Println(dumpExec.Output)
	}

	restore := pg.NewRestore(&pg.Postgres{
		Host:     "localhost",
		Port:     5432,
		DB:       "dev_example",
		Username: "example",
		Password: "example",
	})
	restoreExec := restore.Exec(dumpExec.File, pg.ExecOptions{StreamPrint: false})
	if restoreExec.Error != nil {
		fmt.Println(restoreExec.Error.Err)
		fmt.Println(restoreExec.Output)

	} else {
		fmt.Println("Restore success")
		fmt.Println(restoreExec.Output)

	}
}
