package kubevirt

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	hyperapi "github.com/openshift/hypershift/api"
	apifixtures "github.com/openshift/hypershift/api/fixtures"
	hyperv1 "github.com/openshift/hypershift/api/v1alpha1"
	"github.com/openshift/hypershift/cmd/util"
	"github.com/openshift/hypershift/cmd/version"
	"github.com/spf13/cobra"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Options struct {
	Namespace                      string
	Name                           string
	ReleaseImage                   string
	PullSecretFile                 string
	KubeConfig                     string
	ControlPlaneOperatorImage      string
	NodePoolReplicas               int32
	ControlPlaneAvailabilityPolicy string
	Render                         bool
	RootVolumeType                 string
	RootVolumeIOPS                 int64
	RootVolumeSize                 int64
	NetworkType                    string
}

func NewCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "kubevirt",
		Short:        "Creates basic functional HostedCluster resources on KubeVirt",
		SilenceUsage: true,
	}

	opts := Options{
		Namespace:                      "clusters",
		Name:                           "example",
		ReleaseImage:                   "",
		PullSecretFile:                 "",
		NodePoolReplicas:               2,
		ControlPlaneAvailabilityPolicy: "SingleReplica",
		Render:                         false,
		NetworkType:                    string(hyperv1.OpenShiftSDN),
		RootVolumeType:                 "gp2",
		RootVolumeSize:                 16,
		RootVolumeIOPS:                 0,
	}

	cmd.Flags().StringVar(&opts.Namespace, "namespace", opts.Namespace, "A namespace to contain the generated resources")
	cmd.Flags().StringVar(&opts.Name, "name", opts.Name, "A name for the cluster")
	cmd.Flags().StringVar(&opts.ReleaseImage, "release-image", opts.ReleaseImage, "The OCP release image for the cluster")
	cmd.Flags().StringVar(&opts.PullSecretFile, "pull-secret", opts.PullSecretFile, "Path to a pull secret (required)")
	cmd.Flags().StringVar(&opts.KubeConfig, "kubevirt-kubeconfig", opts.KubeConfig, "Path to the kubeconfig where kubevirt is installed")
	cmd.Flags().Int32Var(&opts.NodePoolReplicas, "node-pool-replicas", opts.NodePoolReplicas, "If >-1, create a default NodePool with this many replicas")
	cmd.Flags().StringVar(&opts.ControlPlaneAvailabilityPolicy, "control-plane-availability-policy", opts.ControlPlaneAvailabilityPolicy, "Availability policy for hosted cluster components. Supported options: SingleReplica, HighlyAvailable")
	cmd.Flags().BoolVar(&opts.Render, "render", opts.Render, "Render output as YAML to stdout instead of applying")

	cmd.MarkFlagRequired("pull-secret")
	cmd.MarkFlagRequired("kubevirt-kubeconfig")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT)
		go func() {
			<-sigs
			cancel()
		}()

		if err := CreateCluster(ctx, opts); err != nil {
			log.Error(err, "Failed to create cluster")
			os.Exit(1)
		}
	}

	return cmd
}

func CreateCluster(ctx context.Context, opts Options) error {
	if len(opts.ReleaseImage) == 0 {
		defaultVersion, err := version.LookupDefaultOCPVersion()
		if err != nil {
			return fmt.Errorf("release image is required when unable to lookup default OCP version: %w", err)
		}
		opts.ReleaseImage = defaultVersion.PullSpec
	}

	client := util.GetClientOrDie()

	pullSecret, err := ioutil.ReadFile(opts.PullSecretFile)
	if err != nil {
		return fmt.Errorf("failed to read pull secret file: %w", err)
	}

	exampleObjects := apifixtures.ExampleOptions{
		Namespace:                      opts.Namespace,
		Name:                           opts.Name,
		ReleaseImage:                   opts.ReleaseImage,
		PullSecret:                     pullSecret,
		NodePoolReplicas:               opts.NodePoolReplicas,
		NetworkType:                    hyperv1.NetworkType(opts.NetworkType),
		ControlPlaneAvailabilityPolicy: hyperv1.AvailabilityPolicy(opts.ControlPlaneAvailabilityPolicy),
		InfraID:                        opts.Name,
		BaseDomain:                     "example.com",
	}.Resources().AsObjects()

	switch {
	case opts.Render:
		for _, object := range exampleObjects {
			err := hyperapi.YamlSerializer.Encode(object, os.Stdout)
			if err != nil {
				return fmt.Errorf("failed to encode objects: %w", err)
			}
			fmt.Println("---")
		}
	default:
		for _, object := range exampleObjects {
			key := crclient.ObjectKeyFromObject(object)
			if err := client.Patch(ctx, object, crclient.Apply, crclient.ForceOwnership, crclient.FieldOwner("hypershift-cli")); err != nil {
				return fmt.Errorf("failed to apply object %q: %w", key, err)
			}
			log.Info("Applied Kube resource", "kind", object.GetObjectKind().GroupVersionKind().Kind, "namespace", key.Namespace, "name", key.Name)
		}
		return nil
	}
	return nil
}
