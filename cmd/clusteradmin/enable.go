package clusteradmin

import (
	"encoding/json"
	"fmt"

	"github.com/openshift/osdctl/pkg/utils"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
)

type subscriptionLabelBody struct {
	Internal     bool   `json:"internal"`
	Key          string `json:"key"`
	ResourceType string `json:"type"`
	Value        string `json:"value"`
}

func newEnableCmd() *cobra.Command {
	//ops := CmdOptions{}
	enableCmd := &cobra.Command{
		Use:   "enable [cluster-id]",
		Short: "enable cluster-admin permissions for given cluster",
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(complete(cmd, args))
			util.CheckErr(runEnable(cmd, args[0]))
		},
	}
	return enableCmd
}

func runEnable(cmd *cobra.Command, clusterid string) error {

	//call func to check if enabled, if true skip else post subscription
	ocmClient := utils.CreateConnection()
	defer ocmClient.Close()

	subID, _ := GetSubscription(ocmClient, clusterid)
	clusterAdminEnabled := GetLabelValue(ocmClient, subID)
	if clusterAdminEnabled == "true" {
		return fmt.Errorf("Cluster admin already enabled")
	}

	body := subscriptionLabelBody{Internal: true, Value: "true"}
	request := ocmClient.Post()
	body.ResourceType = "subscription"
	body.Key = "capability.cluster.manage_cluster_admin"
	request.Path(fmt.Sprintf(subscriptionPath+"/%s"+labelsPath, subID))

	messageBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("cannot marshal template to json: %v", err)
	}
	request.Bytes(messageBytes)

	// Post request
	response, err := request.Send()
	if err != nil {
		return fmt.Errorf("cannot send request: %q", err)
	}
	fmt.Println(response.Status())

	return nil
}
