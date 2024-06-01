package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/urfave/cli/v2"
)

var nextID int

func addShell(con *cli.Context, db *badger.DB) error {
	cmd := exec.Command(con.Args().First())
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	seq, err := db.GetSequence([]byte("nextID"), 1)
	if err != nil {
		return err
	}

	id, err := seq.Next()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("command executed with id: %d\n", id)

	err = db.Update(func(txn *badger.Txn) error {
		key := []byte(strconv.FormatUint(id, 10))
		value := output
		return txn.Set(key, value)
	})
	if err != nil {
		return err
	}

	return nil
}

func checkProgress(con *cli.Context, db *badger.DB) error {
	id, err := strconv.Atoi(con.Args().First())

	if err != nil {
		return fmt.Errorf("invalid ID: %s", con.Args().First())
	}

	var output []byte
	err = db.View(func(txn *badger.Txn) error {
		key := []byte(strconv.Itoa(id))
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		output, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("no command found with id: %d", id)
	}

	fmt.Println(string(output))
	return nil
}

func deleteID(con *cli.Context, db *badger.DB) error {
	id, err := strconv.Atoi(con.Args().First())

	err = db.Update(func(txn *badger.Txn) error {
		key := []byte(strconv.Itoa(id))
		err := txn.Delete([]byte(key))
		return err
	})

	if err != nil {
		return fmt.Errorf("unable to delete: %d", err)
	}

	fmt.Printf("deleted the data with %d id\n", id)
	return nil
}

func listAll(db *badger.DB) error {

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				fmt.Printf("key=%s, value=\n%s\n", k, v)
				return nil
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("unable to list the commands.")
	}

	return nil
}

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
