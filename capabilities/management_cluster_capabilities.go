package capabilities

import (
	routev1 "github.com/openshift/api/route/v1"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
)

// ManagementClusterCapabilities holds all information about optional capabilities of
// the management cluster.
type ManagementClusterCapabilities struct {
	hasRoutesCap bool
}

func (m *ManagementClusterCapabilities) HasRoutes() bool {
	return m.hasRoutesCap
}

// isGroupVersionRegistered determines if a specified groupVersion is registered on the cluster
func isGroupVersionRegistered(client discovery.ServerResourcesInterface, groupVersion schema.GroupVersion) (bool, error) {
	_, apis, err := client.ServerGroupsAndResources()
	if err != nil {
		if discovery.IsGroupDiscoveryFailedError(err) {
			// If the group we are looking for can't be fully discovered,
			// that does still mean that it exists.
			// Continue with the search in the discovered groups if not present here.
			e := err.(*discovery.ErrGroupDiscoveryFailed)
			if _, exists := e.Groups[groupVersion]; exists {
				return true, nil
			}
		} else {
			return false, err
		}
	}

	for _, api := range apis {
		if api.GroupVersion == groupVersion.String() {
			return true, nil
		}
	}

	return false, nil
}

func DetectManagementClusterCapabilities(client discovery.ServerResourcesInterface) (*ManagementClusterCapabilities, error) {
	hasRoutesCap, err := isGroupVersionRegistered(client, routev1.GroupVersion)
	if err != nil {
		return nil, err
	}
	return &ManagementClusterCapabilities{hasRoutesCap: hasRoutesCap}, nil
}
