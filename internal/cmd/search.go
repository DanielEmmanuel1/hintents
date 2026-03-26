// Copyright 2026 Erst Users
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/dotandev/hintents/internal/config"
	"github.com/dotandev/hintents/internal/errors"
	"github.com/dotandev/hintents/internal/rpc"
	"github.com/spf13/cobra"
)

var (
	searchLimitFlag      int
	searchNetworkFlag    string
	searchHorizonURLFlag string
)

var searchCmd = &cobra.Command{
	Use:     "search <query>",
	GroupID: "management",
	Short:   "Find contracts by symbol, creator, or contract ID",
	Long: `Search contracts on Horizon/Soroban-backed networks using one query string.

The query is matched against:
  - contract symbol (when available from Horizon metadata)
  - creator/sponsor account address
  - partial contract ID`,
	Example: `  # Find token contracts by symbol
  erst search usdc --network testnet

  # Find contracts by creator account
  erst search GABCDEF... --network testnet

  # Find contracts by partial contract ID
  erst search CAXYZ --network testnet`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		network := strings.TrimSpace(searchNetworkFlag)
		switch rpc.Network(network) {
		case rpc.Testnet, rpc.Mainnet, rpc.Futurenet:
		default:
			return errors.WrapInvalidNetwork(network)
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		horizonURL := strings.TrimSpace(searchHorizonURLFlag)
		if horizonURL == "" {
			switch rpc.Network(network) {
			case rpc.Mainnet:
				horizonURL = rpc.MainnetHorizonURL
			case rpc.Futurenet:
				horizonURL = rpc.FuturenetHorizonURL
			default:
				horizonURL = rpc.TestnetHorizonURL
			}
		}

		results, err := rpc.SearchContracts(cmd.Context(), rpc.SearchContractsOptions{
			Query:      args[0],
			HorizonURL: horizonURL,
			Limit:      searchLimitFlag,
			Timeout:    time.Duration(cfg.RequestTimeout) * time.Second,
		})
		if err != nil {
			return fmt.Errorf("search failed: %w", err)
		}

		if len(results) == 0 {
			fmt.Println("No matching contracts found.")
			return nil
		}

		fmt.Printf("Found %d matching contracts on %s:\n", len(results), network)
		for _, contract := range results {
			fmt.Println("--------------------------------------------------")
			fmt.Printf("Contract ID: %s\n", contract.ID)
			if contract.Symbol != "" {
				fmt.Printf("Symbol: %s\n", contract.Symbol)
			}
			if contract.Creator != "" {
				fmt.Printf("Creator: %s\n", contract.Creator)
			}
			if contract.LastModifiedLedger > 0 {
				fmt.Printf("Latest Activity Ledger: %d\n", contract.LastModifiedLedger)
			}
			if contract.LastModifiedTime != "" {
				fmt.Printf("Latest Activity Time: %s\n", contract.LastModifiedTime)
			}
		}
		fmt.Println("--------------------------------------------------")
		return nil
	},
}

func init() {
	searchCmd.Flags().IntVar(&searchLimitFlag, "limit", 10, "Maximum number of results to return")
	searchCmd.Flags().StringVarP(&searchNetworkFlag, "network", "n", string(rpc.Testnet), "Stellar network to search (testnet, mainnet, futurenet)")
	searchCmd.Flags().StringVar(&searchHorizonURLFlag, "horizon-url", "", "Override Horizon URL (advanced)")
	_ = searchCmd.Flags().MarkHidden("horizon-url")
	_ = searchCmd.RegisterFlagCompletionFunc("network", completeNetworkFlag)

	rootCmd.AddCommand(searchCmd)
}
