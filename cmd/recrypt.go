package cmd

import (
	"vervet/vervet"

	"github.com/spf13/cobra"
)

func init() {
	recryptCmd.AddCommand(recryptKeySubCmd)
	recryptCmd.AddCommand(recryptClusterSubCmd)

	rootCmd.AddCommand(recryptCmd)
}

var recryptCmd = &cobra.Command{
	Use:   "recrypt",
	Short: "Re-encrypt the unseal key with a new pubkey",
	Long:  `Decrypt PGP-encrypted unseal key and unseal Vault.`,
}

var recryptKeySubCmd = &cobra.Command{
	Use:   "key <unseal key path> <new public key>",
	Short: "Re-encrypt a specific key",
	Long:  `Decrypt unseal key and re-encrypt it using the new public key.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		keyPath := args[0]
		pubKeyPath := args[1]

		keys, err := vervet.ReadKeyFile(keyPath)
		if err != nil {
			vervet.PrintFatal(err.Error(), 1)
		}

		if err := vervet.Recrypt(pubKeyPath, keys); err != nil {
			vervet.PrintFatal(err.Error(), 1)
		}
	},
}

var recryptClusterSubCmd = &cobra.Command{
	Use:   "cluster <cluster name> <new public key>",
	Short: "Re-encrypt the unseal key)s_ for a cluster",
	Long:  `Decrypt unseal key and re-encrypt it using the provided pubkey.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		clusterName := args[0]
		pubKeyPath := args[1]

		cluster, err := getVaultClusterConfig(clusterName)
		if err != nil {
			vervet.PrintFatal(err.Error(), 1)
		}

		keys, err := cluster.keyring()
		if err != nil {
			vervet.PrintFatal(err.Error(), 1)
		}

		if err := vervet.Recrypt(pubKeyPath, keys); err != nil {
			vervet.PrintFatal(err.Error(), 1)
		}
	},
}
