package fixtures

import (
	"os"
	"strconv"

	pg "github.com/guru-golang/pg-commands"
)

func Setup() *pg.Postgres {
	config := &pg.Postgres{
		Host:     "localhost",
		Port:     5432,
		DB:       "dev_example",
		Username: "example",
		Password: "example",
	}
	if os.Getenv("TEST_DB_HOST") != "" {
		config.Host = os.Getenv("TEST_DB_HOST")
	}
	if os.Getenv("TEST_DB_PORT") != "" {
		port, err := strconv.Atoi(os.Getenv("TEST_DB_PORT"))
		if err != nil {
			panic(err)
		}
		config.Port = port
	}
	if os.Getenv("TEST_DB_NAME") != "" {
		config.DB = os.Getenv("TEST_DB_NAME")
	}
	if os.Getenv("TEST_DB_USER") != "" {
		config.Username = os.Getenv("TEST_DB_USER")
	}
	if os.Getenv("TEST_DB_PASS") != "" {
		config.Password = os.Getenv("TEST_DB_PASS")
	}
	return config
}
