package main

import (
	"fmt"
	"github.com/francisco-serrano/sample-blockchain/database"
	"os"
	"time"
)

func main() {
	state, err := database.NewStateFromDisk()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("error while reading state from disk: %w", err))
		os.Exit(1)
	}

	defer state.Close()

	block0 := database.NewBlock(database.Hash{}, uint64(time.Now().Unix()), []database.Tx{
		database.NewTx("andrej", "andrej", 3, ""),
		database.NewTx("andrej", "andrej", 700, "reward"),
	})

	if err := state.AddBlock(block0); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("error while adding block 0: %w", err))
		os.Exit(1)
	}

	block0hash, err := state.Persist()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("error while persisting state for block 0: %w", err))
		os.Exit(1)
	}

	block1 := database.NewBlock(block0hash, uint64(time.Now().Unix()), []database.Tx{
		database.NewTx("andrej", "babayaga", 2000, ""),
		database.NewTx("andrej", "andrej", 100, "reward"),
		database.NewTx("babayaga", "andrej", 1, ""),
		database.NewTx("babayaga", "caesar", 1000, ""),
		database.NewTx("babayaga", "babayaga", 50, ""),
		database.NewTx("andrej", "andrej", 600, "reward"),
	})

	if err := state.AddBlock(block1); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("error while adding block 1: %w", err))
		os.Exit(1)
	}

	if _, err := state.Persist(); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("error while persisting state for block 1: %w", err))
		os.Exit(1)
	}
}
