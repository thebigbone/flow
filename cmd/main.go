package main

import (
	"log"
	"os"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/urfave/cli/v2"
)

func main() {
	db, err := badger.Open(badger.DefaultOptions("./badger").WithLogger(nil))
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	app := &cli.App{
		Name:  "flow",
		Usage: "easily manage your shell commands",
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add shell commands",
				Action: func(con *cli.Context) error {
					return addShell(con, db)
				},
			},
			{
				Name:    "log",
				Aliases: []string{"l"},
				Usage:   "log the output of a shell command",
				Action: func(con *cli.Context) error {
					return checkProgress(con, db)
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"d"},
				Usage:   "delete the command data by specifying it's id",
				Action: func(con *cli.Context) error {
					return deleteID(con, db)
				},
			},
			{
				Name:    "list",
				Aliases: []string{"o"},
				Usage:   "list all the commands and its output",
				Action: func(con *cli.Context) error {
					return listAll(db)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
