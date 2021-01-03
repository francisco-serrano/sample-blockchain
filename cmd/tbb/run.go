package main

import (
	"github.com/spf13/cobra"
)

const flagFrom = "from"
const flagTo = "to"
const flagValue = "value"
const flagData = "data"

func runCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use: "tx",
		Short: "Interact with txs (add...).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	addDefaultRequiredFlags(runCmd)

	return runCmd
}
