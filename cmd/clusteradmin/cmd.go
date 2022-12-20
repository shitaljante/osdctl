package clusteradmin

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	basePath         = "/api/accounts_mgmt/v1"
	labelsPath       = "/labels"
	subscriptionPath = basePath + "/subscriptions"
)

func NewCmdClusterAdmin() *cobra.Command {
	var clusteradminCmd = &cobra.Command{
		Use:   "cluster-admin",
		Short: "manage cluster admin permissions",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				fmt.Println("Error calling cmd.Help(): ", err.Error())
				return
			}
		},
	}

	// Add subcommands
	clusteradminCmd.AddCommand(newEnableCmd())
	clusteradminCmd.AddCommand(newDisableCmd())
	clusteradminCmd.AddCommand(newCheckCmd())

	return clusteradminCmd
}


