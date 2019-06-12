package aws

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// AddToManagerFuncs is a list of functions to add all Controllers to the Manager
var AddToManagerFuncs []func(context.Context, manager.Manager) error

// AddToManager adds all Controllers to the Manager
func (cm *ClusterManager) AddToManager(ctx context.Context, m manager.Manager) error {
	return nil
}
