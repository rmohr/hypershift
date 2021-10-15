package create

import (
	"github.com/openshift/hypershift/cmd/cluster/aws"
	"github.com/spf13/cobra"

	"github.com/openshift/hypershift/cmd/bastion"
	"github.com/openshift/hypershift/cmd/infra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "destroy",
		Short:        "Commands for destroying HyperShift resources",
		SilenceUsage: true,
	}

	cmd.AddCommand(infra.NewDestroyCommand())
	cmd.AddCommand(infra.NewDestroyIAMCommand())
	cmd.AddCommand(aws.NewDestroyCommand())
	cmd.AddCommand(bastion.NewDestroyCommand())

	return cmd
}
