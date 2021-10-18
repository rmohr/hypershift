package cloud

import (
	"github.com/openshift/hypershift/control-plane-operator/controllers/hostedcontrolplane/cloud/aws"
	"github.com/openshift/hypershift/control-plane-operator/controllers/hostedcontrolplane/cloud/kubevirt"
)

func ProviderConfigKey(provider string) string {
	switch provider {
	case aws.Provider:
		return aws.ProviderConfigKey
	case kubevirt.Provider:
		return kubevirt.ProviderConfigKey
	default:
		return ""
	}
}
