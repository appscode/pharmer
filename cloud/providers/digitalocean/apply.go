package digitalocean

import (
	"context"

	api "github.com/pharmer/pharmer/apis/v1beta1"
	"github.com/pharmer/pharmer/cloud"
	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

func (cm *ClusterManager) EnsureMaster(leaderMachine *v1alpha1.Machine) error {
	log := cm.Logger.WithName("ensure-master").WithValues("machine-name", leaderMachine.Name)
	log.Info("ensuring master machine")

	if d, _ := cm.conn.instanceIfExists(leaderMachine); d == nil {
		log.Info("Creating master instance")
		nodeAddresses := make([]core.NodeAddress, 0)

		if cm.Cluster.Status.Cloud.LoadBalancer.DNS != "" {
			nodeAddresses = append(nodeAddresses, core.NodeAddress{
				Type:    core.NodeExternalDNS,
				Address: cm.Cluster.Status.Cloud.LoadBalancer.DNS,
			})
		} else if cm.Cluster.Status.Cloud.LoadBalancer.IP != "" {
			nodeAddresses = append(nodeAddresses, core.NodeAddress{
				Type:    core.NodeExternalIP,
				Address: cm.Cluster.Status.Cloud.LoadBalancer.IP,
			})
		}

		script, err := cloud.RenderStartupScript(cm, leaderMachine, "", customTemplate)
		if err != nil {
			return err
		}

		err = cm.conn.CreateInstance(cm.Cluster, leaderMachine, script)
		if err != nil {
			return err
		}

		if err = cm.Cluster.SetClusterAPIEndpoints(nodeAddresses); err != nil {
			return err
		}
	}
	log.Info("success")

	var err error
	cm.Cluster, err = cm.StoreProvider.Clusters().Update(cm.Cluster)
	if err != nil {
		return err
	}

	return nil
}

func (cm *ClusterManager) PrepareCloud() error {
	log := cm.Logger.WithName("[prepare-cloud]")
	log.Info("preparing cloud infra")

	var found bool
	var err error

	found, _, err = cm.conn.getPublicKey()
	if err != nil {
		return err
	}
	if !found {
		_, err = cm.conn.importPublicKey()
		if err != nil {
			return err
		}
	}

	// ignore errors, since tags are simply informational.
	found, err = cm.conn.getTags()
	if err != nil {
		return err
	}
	if !found {
		if err = cm.conn.createTags(); err != nil {
			return err
		}
	}

	lb, err := cm.conn.lbByName(context.Background(), cm.namer.LoadBalancerName())
	if err == errLBNotFound {
		lb, err = cm.conn.createLoadBalancer(context.Background(), cm.namer.LoadBalancerName())
		if err != nil {
			return err
		}
	}

	cm.Cluster.Status.Cloud.LoadBalancer = api.LoadBalancer{
		IP:   lb.IP,
		Port: lb.ForwardingRules[0].EntryPort,
	}

	nodeAddresses := []corev1.NodeAddress{
		{
			Type:    corev1.NodeExternalIP,
			Address: cm.Cluster.Status.Cloud.LoadBalancer.IP,
		},
	}

	if err = cm.Cluster.SetClusterAPIEndpoints(nodeAddresses); err != nil {
		return errors.Wrap(err, "Error setting controlplane endpoints")
	}

	log.Info("successfully created cloud infra")
	return nil
}

// Deletes master(s) and releases other cloud resources
func (cm *ClusterManager) ApplyDelete() error {
	log := cm.Logger.WithName("[apply-delete]")

	kc, err := cm.GetAdminClient()
	if err != nil {
		return err
	}
	var masterInstances *core.NodeList
	masterInstances, err = kc.CoreV1().Nodes().List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(map[string]string{
			api.RoleMasterKey: "",
		}).String(),
	})
	if err != nil && !kerr.IsNotFound(err) {
		log.Error(err, "master instance not found")
	} else if err == nil {
		for _, mi := range masterInstances.Items {
			err = cm.conn.DeleteInstanceByProviderID(mi.Spec.ProviderID)
			if err != nil {
				log.Error(err, "Failed to delete instance", "instanceID", mi.Spec.ProviderID)
			}
		}
	}

	// delete by tag
	tag := "KubernetesCluster:" + cm.Cluster.Name
	_, err = cm.conn.client.Droplets.DeleteByTag(context.Background(), tag)
	if err != nil {
		log.Error(err, "Failed to delete resources", "tag", tag)
	}
	log.Info("Deleted droplet", "tag", tag)

	// Delete SSH key
	found, _, err := cm.conn.getPublicKey()
	if err != nil {
		return err
	}
	if found {
		err = cm.conn.deleteSSHKey()
		if err != nil {
			return err
		}
	}

	_, err = cm.conn.lbByName(context.Background(), cm.namer.LoadBalancerName())
	if err != errLBNotFound {
		if err = cm.conn.deleteLoadBalancer(context.Background(), cm.namer.LoadBalancerName()); err != nil {
			return err
		}

	}

	cm.Cluster.Status.Phase = api.ClusterDeleted
	_, err = cm.StoreProvider.Clusters().UpdateStatus(cm.Cluster)
	if err != nil {
		return err
	}

	log.Info("successfully deleted cluster")
	return err
}

func (cm *ClusterManager) GetMasterSKU(totalNodes int32) string {
	cm.Logger.Info("setting master sku", "sku", "2gb")
	return "2gb"
}
