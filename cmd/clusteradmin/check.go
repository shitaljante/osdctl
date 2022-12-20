package clusteradmin

import (
	"encoding/json"
	"fmt"

	sdk "github.com/openshift-online/ocm-sdk-go"
	"github.com/openshift/osdctl/pkg/utils"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
)

type labelList struct {
	Items []labels `json:"items"`
}

type labels struct {
	HREF           string `json:"href"`
	ID             string `json:"id"`
	Key            string `json:"key"`
	Internal       bool   `json:"internal"`
	Kind           string `json:"kind"`
	SubscriptionID string `json:"subscription_id"`
	Value          string `json:"value"`
}

func newCheckCmd() *cobra.Command {
	checkCmd := &cobra.Command{
		Use:   "check [cluster-id]",
		Short: "check cluster-admin permissions for given cluster",
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(complete(cmd, args))
			util.CheckErr(runCheck(cmd, args[0]))
		},
	}
	return checkCmd
}

func complete(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Cluster ID is empty, please specifiy which cluster to add cluster-admin permission")
	}
	return nil
}

// runCheck to check if cluster-admin is enabled or disabled
func runCheck(cmd *cobra.Command, clusterid string) error {
	ocmClient := utils.CreateConnection()
	defer ocmClient.Close()

	//Fetch subscriptions from cluster_mgmt
	subID, res := GetSubscription(ocmClient, clusterid)
	if !res {
		return fmt.Errorf("Cluster admin disabled")
	}
	//check label value for "manage_cluster_admin" label
	clusterAdminEnabled := GetLabelValue(ocmClient, subID)
	if clusterAdminEnabled == "true" {
		fmt.Println("Cluster admin is enabled")
	} else {
		fmt.Println("Cluster admin disabled")
	}
	return nil
}

func GetSubscription(ocmClient *sdk.Connection, id string) (string, bool) {
	clusterResp, err := ocmClient.ClustersMgmt().V1().Clusters().Cluster(id).Get().Send()
	if err != nil {
		return err.Error(), false
	}
	clustersubs, _ := clusterResp.Body().GetSubscription()
	subID, res := clustersubs.GetID()
	return subID, res
}

func GetLabelValue(ocmClient *sdk.Connection, subID string) string {
	// Build href path for subscription label
	href := fmt.Sprintf(subscriptionPath+"/%s"+labelsPath, subID)

	// Create a Get request and add the href as the path
	request := ocmClient.Get()
	request.Path(href)

	// Send the request
	response, err := request.Send()
	var labelobj labelList
	err = json.Unmarshal(response.Bytes(), &labelobj)
	if err != nil {
		fmt.Printf("cannot send request: %q", err)
		return "false"
	}
	for _, item := range labelobj.Items{
		if item.Key == "capability.cluster.manage_cluster_admin"{
			return item.Value
		}
	}
	return "false"
}
