package linode

import (
	api "github.com/pharmer/pharmer/apis/v1beta1"
	"github.com/pharmer/pharmer/cloud"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

func (cm *ClusterManager) EnsureMaster(leaderMachine *v1alpha1.Machine) error {
	log := cm.Logger.WithName("ensure-master").WithValues("machine-name", leaderMachine.Name)
	log.Info("ensuring master machine")

	script, err := cloud.RenderStartupScript(cm, leaderMachine, "", customTemplate)
	if err != nil {
		return err
	}

	if _, err = cm.conn.createOrUpdateStackScript(leaderMachine, script); err != nil {
		return err
	}

	if d, _ := cm.conn.instanceIfExists(leaderMachine); d == nil {
		log.Info("Creating master instance")
		var masterServer *api.NodeInfo
		nodeAddresses := make([]corev1.NodeAddress, 0)
		if cm.Cluster.Status.Cloud.LoadBalancer.DNS != "" {
			nodeAddresses = append(nodeAddresses, corev1.NodeAddress{
				Type:    corev1.NodeExternalDNS,
				Address: cm.Cluster.Status.Cloud.LoadBalancer.DNS,
			})
		} else if cm.Cluster.Status.Cloud.LoadBalancer.IP != "" {
			nodeAddresses = append(nodeAddresses, corev1.NodeAddress{
				Type:    corev1.NodeExternalIP,
				Address: cm.Cluster.Status.Cloud.LoadBalancer.IP,
			})
		}

		script, err := cloud.RenderStartupScript(cm, leaderMachine, "", customTemplate)
		if err != nil {
			return err
		}

		masterServer, err = cm.conn.CreateInstance(leaderMachine, script)
		if err != nil {
			return err
		}

		if err = cm.conn.addNodeToBalancer(cm.namer.LoadBalancerName(), leaderMachine.Name, masterServer.PrivateIP); err != nil {
			return err
		}

		if err = cm.Cluster.SetClusterAPIEndpoints(nodeAddresses); err != nil {
			return err
		}
		if _, err = cm.StoreProvider.Clusters().Update(cm.Cluster); err != nil {
			return err
		}
	}

	return nil
}

func (cm *ClusterManager) PrepareCloud() error {
	log := cm.Logger.WithName("[prepare-cloud]")
	log.Info("preparing cloud infra")

	lb, err := cm.conn.lbByName(cm.namer.LoadBalancerName())

	if err != nil {

		lb, err = cm.conn.createLoadBalancer(cm.namer.LoadBalancerName())
		if err != nil {
			return err
		}
	}

	cm.Cluster.Status.Cloud.LoadBalancer = api.LoadBalancer{
		IP:   *lb.IPv4,
		Port: api.DefaultKubernetesBindPort,
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

	return err
}

//Deletes master(s) and releases other cloud resources
func (cm *ClusterManager) ApplyDelete() error {
	log := cm.Logger.WithName("[apply-delete]")

	var kc kubernetes.Interface
	kc, err := cm.GetAdminClient()
	if err != nil {
		return err
	}
	var nodeInstances *corev1.NodeList
	nodeInstances, err = kc.CoreV1().Nodes().List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(map[string]string{
			api.RoleNodeKey: "",
		}).String(),
	})
	if err != nil && !kerr.IsNotFound(err) {
		log.Error(err, "node instance not found.=")
	} else if err == nil {
		for _, mi := range nodeInstances.Items {

			if err = cm.conn.DeleteStackScript(mi.Name, api.RoleNode); err != nil {
				log.Error(err, "Unable to delete stack script")
			}
			err = kc.CoreV1().Nodes().Delete(mi.Name, nil)
			if err != nil {
				log.Error(err, "Failed to delete node.", "node-name", mi.Name)
			}
		}
	}

	var masterInstances *corev1.NodeList
	masterInstances, err = kc.CoreV1().Nodes().List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(map[string]string{
			api.RoleMasterKey: "",
		}).String(),
	})
	if err != nil && !kerr.IsNotFound(err) {
		log.Error(err, "master instance not found.")
	} else if err == nil {
		for _, mi := range masterInstances.Items {
			if err = cm.conn.DeleteStackScript(mi.Name, api.RoleMaster); err != nil {
				log.Error(err, "failed to delete stack script")
			}

			err = cm.conn.DeleteInstanceByProviderID(mi.Spec.ProviderID)
			if err != nil {
				log.Error(err, "Failed to delete instance.", "Instance ID", mi.Spec.ProviderID)
			}
		}
	}

	_, err = cm.conn.lbByName(cm.namer.LoadBalancerName())
	if err != errLBNotFound {
		if err = cm.conn.deleteLoadBalancer(cm.namer.LoadBalancerName()); err != nil {
			return err
		}
	}

	cm.Cluster.Status.Phase = api.ClusterDeleted
	_, err = cm.StoreProvider.Clusters().UpdateStatus(cm.Cluster)
	if err != nil {
		return err
	}

	log.Info("successfully deleted cluster")
	return nil
}

func (cm *ClusterManager) GetMasterSKU(totalNodes int32) string {
	cm.Logger.Info("setting master sku", "sku", "g6-standard-2")
	return "g6-standard-2"
}
