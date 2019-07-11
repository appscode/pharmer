package eks

import (
	"fmt"

	api "github.com/pharmer/pharmer/apis/v1beta1"
)

type namer struct {
	cluster *api.Cluster
}

func (n namer) GetStackServiceRole() string {
	return fmt.Sprintf("EKS-%v-ServiceRole", n.cluster.Name)
}

func (n namer) GetClusterVPC() string {
	return fmt.Sprintf("EKS-%v-VPC", n.cluster.Name)
}
