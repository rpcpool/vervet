package cmd

import (
	"vervet/vervet"

	"github.com/spf13/cobra"
)

func init() {
	rekeyServerSubCmd.Flags().IntVarP(&vaultPort, "port", "p", 8200, "Vault API port")
	rekeyServerSubCmd.Flags().BoolVarP(&vaultTLSDisable, "insecure", "i", false, "disable TLS")
	rekeyServerSubCmd.Flags().StringVarP(&vaultRekeyNonce, "nonce", "n", "", "nonce for root token generation")

	rekeyClusterSubCmd.Flags().StringVarP(&vaultRekeyNonce, "nonce", "n", "", "nonce for root token generation")

	rekeyCmd.AddCommand(rekeyServerSubCmd)
	rekeyCmd.AddCommand(rekeyClusterSubCmd)

	rootCmd.AddCommand(rekeyCmd)

}

var rekeyCmd = &cobra.Command{
	Use:   "rekey",
	Short: "Rekey Vault",
	Long:  `Decrypt the unseal key and rekey the Vault cluster.`,
}

var rekeyServerSubCmd = &cobra.Command{
	Use:   "server <vault address> <unseal key path> -n <nonce>",
	Short: "Rekey Vault cluster",
	Long:  `Decrypt the unseal key and rekey Vault root token.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		vaultAddr := getVaultAddress(args[0])
		keyPath := args[1]

		keys, err := vervet.ReadKeyFile(keyPath)
		if err != nil {
			vervet.PrintFatal(err.Error(), 1)
		}

		if err := vervet.Rekey(vaultAddr, keys); err != nil {
			vervet.PrintFatal(err.Error(), 1)
		}
	},
}

var rekeyClusterSubCmd = &cobra.Command{
	Use:   "cluster <cluster name> -n <nonce>",
	Short: "Rekey the Vault cluster",
	Long:  `Decrypt the unseal key and rekey Vault.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		clusterName := args[0]

		cluster, err := getVaultClusterConfig(clusterName)
		if err != nil {
			vervet.PrintFatal(err.Error(), 1)
		}

		keys, err := cluster.keyring()
		if err != nil {
			vervet.PrintFatal(err.Error(), 1)
		}

		if len(cluster.Servers) == 0 {
			vervet.PrintFatal("no Vault servers in configuration", 1)
		}

		if err := vervet.Rekey(cluster.Servers[0], vervet.Unique(keys)); err != nil {
			vervet.PrintFatal(err.Error(), 1)
		}
	},
}
