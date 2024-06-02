package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/urfave/cli/v2"
)

func addShell(con *cli.Context, db *badger.DB) error {
	cmd := exec.Command(con.Args().First())
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	seq, err := db.GetSequence([]byte("id"), 1)
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
