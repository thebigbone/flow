package main

import (
	"encoding/json"
	"fmt"
	"os"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/jedib0t/go-pretty/v6/table"
)

func listAll(db *badger.DB) error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Command", "Output", "Start Time", "Total Time"})
	t.SetStyle(table.StyleLight)

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			if string(k) == "id" {
				continue
			}

			v, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			var sc shellCommand
			err = json.Unmarshal(v, &sc)
			if err != nil {
				return err
			}

			t.AppendRows([]table.Row{
				{string(k), sc.Command, sc.Output, sc.StartTime, sc.EndTime},
			})
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("unable to list the commands: %s", err)
	}

	t.Render()
	return nil
}
