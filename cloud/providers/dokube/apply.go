package dokube

import (
	"context"

	"github.com/pharmer/pharmer/cloud"

	"github.com/appscode/go/log"
	api "github.com/pharmer/pharmer/apis/v1beta1"
	"github.com/pharmer/pharmer/store"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

func (cm *ClusterManager) PrepareCloud() error {

	if cm.Cluster.Spec.Config.Cloud.Dokube.ClusterID == "" {
		cluster, err := cm.conn.createCluster(cm.Cluster)
		if err != nil {
			return err
		}
		cm.Cluster.Spec.Config.Cloud.Dokube.ClusterID = cluster.ID
		if _, err = store.StoreProvider.Clusters().Update(cm.Cluster); err != nil {
			return err
		}
		if err := cm.retrieveClusterStatus(cluster); err != nil {
			return err
		}

		err = cm.StoreCertificate(cm.conn.client)
		if err != nil {
			log.Infof(err.Error())
			return err
		}
		certs, err := cloud.GetPharmerCerts(cm.Cluster.Name)
		if err != nil {
			return err
		}

		cm.Certs = certs
	}

	return nil
}

func (cm *ClusterManager) ApplyScale() error {
	var nodeGroups []*clusterapi.MachineSet
	nodeGroups, err := store.StoreProvider.MachineSet(cm.Cluster.Name).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, ng := range nodeGroups {
		igm := NewDokubeNodeGroupManager(cm.conn, ng)

		err = igm.Apply()
		if err != nil {
			return err
		}
	}
	_, err = store.StoreProvider.Clusters().UpdateStatus(cm.Cluster)
	if err != nil {
		return err
	}
	_, err = store.StoreProvider.Clusters().Update(cm.Cluster)
	if err != nil {
		return err
	}
	return nil
}

func (cm *ClusterManager) ApplyDelete() error {
	if cm.Cluster.Status.Phase == api.ClusterReady {
		cm.Cluster.Status.Phase = api.ClusterDeleting
	}
	_, err := store.StoreProvider.Clusters().UpdateStatus(cm.Cluster)
	if err != nil {
		return err
	}
	_, err = cm.conn.client.Kubernetes.Delete(context.Background(), cm.conn.Cluster.Spec.Config.Cloud.Dokube.ClusterID)
	if err != nil {
		return err
	}
	cm.Cluster.Status.Phase = api.ClusterDeleted
	_, err = store.StoreProvider.Clusters().Update(cm.Cluster)
	if err != nil {
		return err
	}

	return nil
}
