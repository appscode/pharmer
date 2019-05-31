package cloud

import (
	"context"

	api "github.com/pharmer/pharmer/apis/v1beta1"
	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var (
	ErrNotFound       = errors.New("node not found")
	ErrNotImplemented = errors.New("not implemented")
	ErrNoMasterNG     = errors.New("Cluster has no master NodeGroup")
)

type Interface interface {
	CloudManagerInterface

	CreateCCMCredential() error

	GetConnector() ClusterApiProviderComponent

	InitializeMachineActuator(mgr manager.Manager) error
	AddToManager(ctx context.Context, m manager.Manager) error

	SSHGetter
	ProviderKubeConfig

	GetCloudConnector() error
	ApplyCreate(dryRun bool) (acts []api.Action, leaderMachine *clusterv1.Machine, machines []*clusterv1.Machine, err error)
	ApplyDelete(dryRun bool) ([]api.Action, error)
	IsValid(cluster *api.Cluster) (bool, error)
	SetDefaultCluster(in *api.Cluster) error
	GetDefaultMachineProviderSpec(cluster *api.Cluster, sku string, role api.MachineRole) (clusterv1.ProviderSpec, error)

	NewMasterTemplateData(machine *clusterapi.Machine, token string, td TemplateData) TemplateData
	NewNodeTemplateData(machine *clusterapi.Machine, token string, td TemplateData) TemplateData
}

type SSHGetter interface {
	GetSSHConfig(cluster *api.Cluster, node *core.Node) (*api.SSHConfig, error)
}

type NodeGroupManager interface {
	//	Apply(dryRun bool) (acts []api.Action, err error)
	//	AddNodes(count int64) error
	//	DeleteNodes(nodes []core.Node) error
}

type InstanceManager interface {
	CreateInstance(cluster *api.Cluster, machine *clusterv1.Machine, token string) (*api.NodeInfo, error)
	DeleteInstanceByProviderID(providerID string) error
}

// TODO: change name
type ClusterApiProviderComponent interface {
	CreateCredentialSecret(kc kubernetes.Interface, data map[string]string) error
	GetControllerManager() (string, error)
}

type UpgradeManager interface {
	GetAvailableUpgrades() ([]*api.Upgrade, error)
	PrintAvailableUpgrades([]*api.Upgrade)
	Apply(dryRun bool) ([]api.Action, error)
	MasterUpgrade(oldMachine *clusterv1.Machine, newMachine *clusterv1.Machine) error
	NodeUpgrade(oldMachine *clusterv1.Machine, newMachine *clusterv1.Machine) error
}

type ProviderKubeConfig interface {
	GetKubeConfig(cluster *api.Cluster) (*api.KubeConfig, error)
}

type HookFunc func() error
