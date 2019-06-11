package aws

import (
	api "github.com/pharmer/pharmer/apis/v1beta1"
	"github.com/pharmer/pharmer/cloud"
	. "github.com/pharmer/pharmer/cloud"
)

type ClusterManager struct {
	*cloud.CloudManager

	conn  *cloudConnector
	namer namer
}

func (cm *ClusterManager) GetConnector() ClusterApiProviderComponent {
	panic(1)
	return nil
}

func (cm *ClusterManager) CreateCCMCredential() error {
	return nil
}

var _ Interface = &ClusterManager{}

const (
	UID = "aws"
)

func init() {
	RegisterCloudManager(UID, func(cluster *api.Cluster, certs *PharmerCertificates) Interface {
		return New(cluster, certs)
	})
}

func New(cluster *api.Cluster, certs *PharmerCertificates) cloud.Interface {
	return &ClusterManager{
		CloudManager: &cloud.CloudManager{
			Cluster: cluster,
			Certs:   certs,
		},
		namer: namer{
			cluster: cluster,
		},
	}
}
