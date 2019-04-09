package azure

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeadmv1beta1 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AzureClusterProviderSpec is the providerConfig for Azure in the cluster.
// +k8s:openapi-gen=true
type AzureClusterProviderSpec struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// NetworkSpec encapsulates all things related to Azure network.
	NetworkSpec NetworkSpec `json:"networkSpec,omitempty"`

	ResourceGroup string `json:"resourceGroup"`
	Location      string `json:"location"`

	// CAKeyPair is the key pair for CA certs.
	CAKeyPair KeyPair `json:"caKeyPair,omitempty"`

	// EtcdCAKeyPair is the key pair for etcd.
	EtcdCAKeyPair KeyPair `json:"etcdCAKeyPair,omitempty"`

	// FrontProxyCAKeyPair is the key pair for the front proxy.
	FrontProxyCAKeyPair KeyPair `json:"frontProxyCAKeyPair,omitempty"`

	// SAKeyPair is the service account key pair.
	SAKeyPair KeyPair `json:"saKeyPair,omitempty"`

	// AdminKubeconfig generated using the certificates part of the spec
	// do not move to status, since it uses on disk ca certs, which causes issues during regeneration
	AdminKubeconfig string `json:"adminKubeconfig,omitempty"`

	// DiscoveryHashes generated using the certificates part of the spec, used by master and nodes bootstrapping
	// this never changes until ca is rotated
	// do not move to status, since it uses on disk ca certs, which causes issues during regeneration
	DiscoveryHashes []string `json:"discoveryHashes,omitempty"`

	// ClusterConfiguration holds the cluster-wide information used during a
	// kubeadm init call.
	ClusterConfiguration kubeadmv1beta1.ClusterConfiguration `json:"clusterConfiguration,omitempty"`
}

// KeyPair is how operators can supply custom keypairs for kubeadm to use.
type KeyPair struct {
	// base64 encoded cert and key
	Cert []byte `json:"cert"`
	Key  []byte `json:"key"`
}

// HasCertAndKey returns whether a keypair contains cert and key of non-zero length.
func (kp *KeyPair) HasCertAndKey() bool {
	return len(kp.Cert) != 0 && len(kp.Key) != 0
}

// NetworkSpec encapsulates all things related to Azure network.
type NetworkSpec struct {
	// Vnet configuration.
	// +optional
	Vnet VnetSpec `json:"vnet,omitempty"`

	// Subnets configuration.
	// +optional
	Subnets Subnets `json:"subnets,omitempty"`
}

// VnetSpec configures an Azure virtual network.
type VnetSpec struct {
	// ID is the identifier of the virtual network this provider should use to create resources.
	ID string `json:"id,omitempty"`

	// Name defines a name for the virtual network resource.
	Name string `json:"name"`

	// CidrBlock is the CIDR block to be used when the provider creates a managed virtual network.
	CidrBlock string `json:"cidrBlock,omitempty"`

	// Tags is a collection of tags describing the resource.
	// TODO: Uncomment once tagging is implemented.
	//Tags tags.Map `json:"tags,omitempty"`
}

// IsProvided returns true if the virtual network is not managed by Cluster API.
// TODO: Uncomment once tagging is implemented.
/*
func (v *VnetSpec) IsProvided() bool {
	return v.ID != "" && !v.Tags.HasManaged()
}
*/

// SubnetSpec configures an Azure subnet.
type SubnetSpec struct {
	// ID defines a unique identifier to reference this resource.
	ID string `json:"id,omitempty"`

	// Name defines a name for the subnet resource.
	Name string `json:"name"`

	// VnetID defines the ID of the virtual network this subnet should be built in.
	VnetID string `json:"vnetId"`

	// CidrBlock is the CIDR block to be used when the provider creates a managed Vnet.
	CidrBlock string `json:"cidrBlock,omitempty"`

	// SecurityGroup defines the NSG (network security group) that should be attached to this subnet.
	SecurityGroup SecurityGroup `json:"securityGroup"`

	// Tags is a collection of tags describing the resource.
	// TODO: Uncomment once tagging is implemented.
	//Tags tags.Map `json:"tags,omitempty"`
}
