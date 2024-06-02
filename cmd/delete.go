package main

import (
	"fmt"
	"log"
	"strconv"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/urfave/cli/v2"
)

func deleteID(con *cli.Context, db *badger.DB) error {
	id, err := strconv.Atoi(con.Args().First())

	if err != nil {
		log.Fatal(err)
	}

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
