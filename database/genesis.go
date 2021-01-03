package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var genesisJSON = `
{
  "genesis_time": "2019-03-18T00:00:00.000000000Z",
  "chain_id": "the-blockchain-bar-ledger",
  "balances": {
    "andrej": 1000000
  }
}
`

type genesis struct {
	Balances map[Account]uint `json:"balances"`
}

func loadGenesis(path string) (genesis, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return genesis{}, fmt.Errorf("error while reading genesis file: %w", err)
	}

	var loadedGenesis genesis
	if err := json.Unmarshal(content, &loadedGenesis); err != nil {
		return genesis{}, fmt.Errorf("error while unmarshalling file content: %w", err)
	}

	return loadedGenesis, nil
}

func writeGenesisToDisk(path string) error {
	return ioutil.WriteFile(path, []byte(genesisJSON), 0644)
}
