package packet

import (
	"encoding/json"

	api "github.com/pharmer/pharmer/apis/v1beta1"
	packetconfig "github.com/pharmer/pharmer/apis/v1beta1/packet"
	"github.com/pharmer/pharmer/cloud"
	"github.com/pharmer/pharmer/cloud/utils/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

func (cm *ClusterManager) GetDefaultMachineProviderSpec(sku string, role api.MachineRole) (clusterapi.ProviderSpec, error) {
	if sku == "" {
		sku = "baremetal_0"
	}
	//config := cluster.Spec.Config
	spec := &packetconfig.PacketMachineProviderSpec{
		TypeMeta: metav1.TypeMeta{
			APIVersion: packetconfig.PacketProviderGroupName + "/" + packetconfig.PacketProviderAPIVersion,
			Kind:       packetconfig.PacketProviderKind,
		},
		Plan:         sku,
		SpotInstance: "Regular",
	}

	providerSpecValue, err := json.Marshal(spec)
	if err != nil {
		return clusterapi.ProviderSpec{}, err
	}

	return clusterapi.ProviderSpec{
		Value: &runtime.RawExtension{
			Raw: providerSpecValue,
		},
	}, nil
}

func (cm *ClusterManager) SetDefaultCluster() error {
	cluster := cm.Cluster
	config := &cluster.Spec.Config

	config.Cloud.InstanceImage = "ubuntu_16_04" // 1b9b78e3-de68-466e-ba00-f2123e89c112
	config.SSHUserName = "root"

	return packetconfig.SetPacketClusterProviderConfig(&cluster.Spec.ClusterAPI)
}

func (cm *ClusterManager) IsValid(cluster *api.Cluster) (bool, error) {
	return false, cloud.ErrNotImplemented
}

func (cm *ClusterManager) GetKubeConfig() (*api.KubeConfig, error) {
	return kube.GetAdminConfig(cm.Cluster, cm.GetCaCertPair())
}
