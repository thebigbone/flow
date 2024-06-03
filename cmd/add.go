package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"time"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/urfave/cli/v2"
)

type shellCommand struct {
	Command   string        `json:"command"`
	Output    string        `json:"output"`
	StartTime time.Time     `json:"start"`
	EndTime   time.Duration `json:"end"`
}

func addShell(con *cli.Context, db *badger.DB) error {
	cmd_name := con.Args().First()

	start := time.Now()
	cmd := exec.Command(cmd_name)

	elapsed := time.Since(start)

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

	sc := shellCommand{
		Command:   cmd_name,
		Output:    string(output),
		StartTime: start,
		EndTime:   elapsed,
	}

	data, err := json.Marshal(sc)
	if err != nil {
		return err
	}

	err = db.Update(func(txn *badger.Txn) error {
		key := []byte(strconv.FormatUint(id, 10))
		// value := output
		return txn.Set(key, data)
	})
	if err != nil {
		return err
	}

	return nil
}
