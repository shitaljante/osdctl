package clusteradmin

import (
	"fmt"

	"github.com/openshift/osdctl/pkg/utils"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
)

func newDisableCmd() *cobra.Command {
	//ops := CmdOptions{}
	disableCmd := &cobra.Command{
		Use:   "disable [cluster-id]",
		Short: "disable cluster-admin permissions for given cluster",
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(complete(cmd, args))
			util.CheckErr(runDisable(cmd, args[0]))
		},
	}

	return disableCmd
}

func runDisable(cmd *cobra.Command, clusterid string) error {
	//before delete request check if cluster-admin already disabled, if true skip or delete

	ocmClient := utils.CreateConnection()
	defer ocmClient.Close()

	subID, _ := GetSubscription(ocmClient, clusterid)
	clusterAdminEnabled := GetLabelValue(ocmClient, subID)
	if clusterAdminEnabled == "false" {
		return fmt.Errorf("Cluster admin already disabled")
	}

	//fetch users
	usersResponse, _ := ocmClient.ClustersMgmt().V1().Clusters().Cluster(clusterid).Groups().Group("cluster-admins").Users().List().Send()
	userList := usersResponse.Items().Slice()
	for _, user := range userList {
		//delete users
		resp, _ := ocmClient.ClustersMgmt().V1().Clusters().Cluster(clusterid).Groups().Group("cluster-admins").Users().User(user.ID()).Delete().Send()
		fmt.Println(resp.Status())
	}

	href := fmt.Sprintf(subscriptionPath +"/%s"+ labelsPath + "/capability.cluster.manage_cluster_admin", subID)
	request := ocmClient.Delete()
	request.Path(href)

	// Send the request
	response, err := request.Send()
	fmt.Println(response.Status())
	if err != nil {
		return fmt.Errorf("cannot send request: %q", err)
	}
	return nil
}
