package create

import (
	"github.com/openshift/hypershift/cmd/cluster/aws"
	"github.com/openshift/hypershift/cmd/cluster/kubevirt"
	"github.com/spf13/cobra"

	"github.com/openshift/hypershift/cmd/bastion"
	"github.com/openshift/hypershift/cmd/infra"
	"github.com/openshift/hypershift/cmd/kubeconfig"
	"github.com/openshift/hypershift/cmd/nodepool"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "create",
		Short:        "Commands for creating HyperShift resources",
		SilenceUsage: true,
	}

	clusterCmd := &cobra.Command{
		Use:          "cluster",
		Short:        "Creates basic functional HostedCluster resources",
		SilenceUsage: true,
	}
	clusterCmd.AddCommand(aws.NewCreateCommand())
	clusterCmd.AddCommand(kubevirt.NewCreateCommand())
	cmd.AddCommand(clusterCmd)
	cmd.AddCommand(infra.NewCreateCommand())
	cmd.AddCommand(infra.NewCreateIAMCommand())
	cmd.AddCommand(kubeconfig.NewCreateCommand())
	cmd.AddCommand(nodepool.NewCreateCommand())
	cmd.AddCommand(bastion.NewCreateCommand())

	return cmd
}
