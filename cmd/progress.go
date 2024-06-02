package main

import (
	"fmt"
	"strconv"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/urfave/cli/v2"
)

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
