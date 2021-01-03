package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type State struct {
	Balances  map[Account]uint
	txMempool []Tx

	dbFile          *os.File
	latestBlockHash Hash
}

func NewStateFromDisk(dataDir string) (*State, error) {
	if err := initDataDirIfNotExists(dataDir); err != nil {
		return nil, fmt.Errorf("error while initializing data: %w", err)
	}

	gen, err := loadGenesis(getGenesisJsonFilePath(dataDir))
	if err != nil {
		return nil, fmt.Errorf("error while loading genesis: %w", err)
	}

	balances := make(map[Account]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	f, err := os.OpenFile(getBlocksDbFilePath(dataDir), os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, fmt.Errorf("error while opening tx file: %w", err)
	}

	scanner := bufio.NewScanner(f)

	state := &State{
		Balances:        balances,
		txMempool:       make([]Tx, 0),
		dbFile:          f,
		latestBlockHash: Hash{},
	}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error while scanning: %w", err)
		}

		var blockFS BlockFS
		if err := json.Unmarshal(scanner.Bytes(), &blockFS); err != nil {
			return nil, fmt.Errorf("error while unmarshalling bytes: %w", err)
		}

		if err := state.applyBlock(blockFS.Value); err != nil {
			return nil, fmt.Errorf("error while applying block: %w", err)
		}

		state.latestBlockHash = blockFS.Key
	}

	return state, nil
}

func (s *State) LatestBlockHash() Hash {
	return s.latestBlockHash
}

func (s *State) AddBlock(b Block) error {
	for _, tx := range b.TXs {
		if err := s.AddTx(tx); err != nil {
			return fmt.Errorf("error while adding tx: %w", err)
		}
	}

	return nil
}

func (s *State) AddTx(tx Tx) error {
	if err := s.apply(tx); err != nil {
		return fmt.Errorf("error while applying tx: %w", err)
	}

	s.txMempool = append(s.txMempool, tx)

	return nil
}

func (s *State) Persist() (Hash, error) {
	block := NewBlock(s.latestBlockHash, uint64(time.Now().Unix()), s.txMempool)

	blockHash, err := block.Hash()
	if err != nil {
		return Hash{}, fmt.Errorf("error while hashing block: %w", err)
	}

	blockFS := BlockFS{blockHash, block}

	blockFsJson, err := json.Marshal(blockFS)
	if err != nil {
		return Hash{}, fmt.Errorf("error while marshalling blockFS: %w", err)
	}

	fmt.Printf("persisting new Block to disk:\n")
	fmt.Printf("\t%s\n", blockFsJson)

	if _, err := s.dbFile.Write(append(blockFsJson, '\n')); err != nil {
		return Hash{}, fmt.Errorf("error while writing on file: %w", err)
	}

	s.latestBlockHash = blockHash

	// reset mempool after flush to disk
	s.txMempool = []Tx{}

	return s.latestBlockHash, nil
}

func (s *State) Close() error {
	return s.dbFile.Close()
}

func (s *State) applyBlock(b Block) error {
	for _, tx := range b.TXs {
		if err := s.apply(tx); err != nil {
			return fmt.Errorf("error while applying tx: %w", err)
		}
	}

	return nil
}

func (s *State) apply(tx Tx) error {
	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}

	if s.Balances[tx.From] < tx.Value {
		return fmt.Errorf("insufficient balance")
	}

	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value

	return nil
}
