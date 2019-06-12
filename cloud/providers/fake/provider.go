package fake

import (
	"context"

	v1 "github.com/pharmer/cloud/pkg/apis/cloud/v1"
	api "github.com/pharmer/pharmer/apis/v1beta1"
	. "github.com/pharmer/pharmer/cloud"
	core "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type ClusterManager struct {
	cluster *api.Cluster
	certs   *PharmerCertificates

	cfg   *api.PharmerConfig
	owner string
}

func (cm *ClusterManager) GetCluster() *api.Cluster {
	panic("implement me")
}

func (cm *ClusterManager) GetCaCertPair() *CertKeyPair {
	panic("implement me")
}

func (cm *ClusterManager) GetPharmerCertificates() *PharmerCertificates {
	panic("implement me")
}

func (cm *ClusterManager) GetCredential() (*v1.Credential, error) {
	panic("implement me")
}

func (cm *ClusterManager) GetAdminClient() (kubernetes.Interface, error) {
	panic("implement me")
}

func (cm *ClusterManager) CreateCCMCredential() error {
	panic("implement me")
}

func (cm *ClusterManager) GetConnector() ClusterApiProviderComponent {
	panic("implement me")
}

func (cm *ClusterManager) GetCloudConnector() error {
	panic("implement me")
}

func (cm *ClusterManager) ApplyCreate(dryRun bool) (acts []api.Action, leaderMachine *v1alpha1.Machine, machines []*v1alpha1.Machine, err error) {
	panic("implement me")
}

func (cm *ClusterManager) ApplyDelete(dryRun bool) ([]api.Action, error) {
	panic("implement me")
}

func (cm *ClusterManager) SetDefaultCluster() error {
	panic("implement me")
}

func (cm *ClusterManager) GetDefaultMachineProviderSpec(sku string, role api.MachineRole) (v1alpha1.ProviderSpec, error) {
	panic("implement me")
}

func (cm *ClusterManager) NewMasterTemplateData(machine *v1alpha1.Machine, token string, td TemplateData) TemplateData {
	panic("implement me")
}

func (cm *ClusterManager) NewNodeTemplateData(machine *v1alpha1.Machine, token string, td TemplateData) TemplateData {
	panic("implement me")
}

// AddToManager adds all Controllers to the Manager
func (cm *ClusterManager) AddToManager(ctx context.Context, m manager.Manager) error {
	return ErrNotImplemented
}

//func (cm *ClusterManager) GetDefaultMachineProviderSpec(cluster *api.Cluster, sku string, role api.MachineRole) (v1alpha1.ProviderSpec, error) {
//	return v1alpha1.ProviderSpec{}, ErrNotImplemented
//}

func (cm *ClusterManager) InitializeMachineActuator(mgr manager.Manager) error {
	return ErrNotImplemented
}

//func (cm *ClusterManager) SetDefaultCluster(in *api.Cluster) error {
//	return ErrNotImplemented
//}

var _ Interface = &ClusterManager{}

const (
	UID = "fake"
)

func init() {
	RegisterCloudManager(UID, func(cluster *api.Cluster, certs *PharmerCertificates) Interface {
		return New(cluster, certs)
	})
}

func New(cluster *api.Cluster, certs *PharmerCertificates) Interface {
	return &ClusterManager{
		cluster: cluster,
		certs:   certs,
	}
}

func (cm *ClusterManager) SetDefaults(in *api.Cluster) error {
	return nil
}

func (cm *ClusterManager) SetOwner(owner string) {
	cm.owner = owner
}

func (cm *ClusterManager) GetDefaultNodeSpec(cluster *api.Cluster, sku string) (api.NodeSpec, error) {
	return api.NodeSpec{}, nil
}

func (cm *ClusterManager) Apply(in *api.Cluster, dryRun bool) ([]api.Action, error) {
	return nil, ErrNotImplemented
}

func (cm *ClusterManager) IsValid(cluster *api.Cluster) (bool, error) {
	return false, ErrNotImplemented
}

func (cm *ClusterManager) UploadStartupConfig() error {
	return nil
}

func (cm *ClusterManager) runFakeJob(requestType string) {
	//c.Logger().Infof("starting %v job", requestType)
	//for i := 1; i <= 10; i++ {
	//	c.Logger().Info(fmt.Sprint("Job completed: ", i*10, "%"))
	//	time.Sleep(time.Second * 3)
	//}
}

func (cm *ClusterManager) GetSSHConfig(cluster *api.Cluster, node *core.Node) (*api.SSHConfig, error) {
	return nil, ErrNotImplemented
}

func (cm *ClusterManager) GetKubeConfig(cluster *api.Cluster) (*api.KubeConfig, error) {
	return nil, ErrNotImplemented
}
