package main

import (
	"fmt"
	"os"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/jedib0t/go-pretty/v6/table"
)

func listAll(db *badger.DB) error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Output"})
	t.SetStyle(table.StyleLight)

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				if string(k) == "id" {
					return nil
				}
				t.AppendRows([]table.Row{
					{string(k), string(v)},
				})
				//fmt.Printf("key=%s, value=\n%s\n", k, v)
				return nil
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("unable to list the commands: %s", err)
	}

	t.Render()
	return nil
}
