package cloud

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pharmer/pharmer/store"

	"k8s.io/kubernetes/pkg/apis/core"

	"github.com/appscode/go/log"
	. "github.com/appscode/go/types"
	"github.com/pharmer/pharmer/apis/v1beta1"
	api "github.com/pharmer/pharmer/apis/v1beta1"
	"github.com/pharmer/pharmer/cloud/cmds/options"
	"github.com/pkg/errors"
	"gomodules.xyz/cert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/clusterdeployer/clusterclient"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"sigs.k8s.io/cluster-api/pkg/util"
)

var managedProviders = sets.NewString("aks", "gke", "eks", "dokube")

func List(ctx context.Context, opts metav1.ListOptions, owner string) ([]*api.Cluster, error) {
	return Store(ctx).Owner(owner).Clusters().List(opts)
}

func Get(name string, owner string) (*api.Cluster, error) {
	return store.StoreProvider.Owner(owner).Clusters().Get(name)
}

func Create(cluster *api.Cluster, owner string) (*api.Cluster, error) {
	config := cluster.Spec.Config
	if cluster == nil {
		return nil, errors.New("missing cluster")
	} else if cluster.Name == "" {
		return nil, errors.New("missing cluster name")
	} else if config.KubernetesVersion == "" {
		return nil, errors.New("missing cluster version")
	}

	exists := false
	_, err := store.StoreProvider.Owner(owner).Clusters().Get(cluster.Name)
	if err == nil {
		exists = true
	}

	cm, err := GetCloudManager(config.Cloud.CloudProvider)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get cloud manager")
	}
	cm.SetOwner(owner)

	// set common cluster configs
	err = SetDefaultCluster(cluster)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set default cluster")
	}

	// set cloud-specific configs
	if err = cm.SetDefaultCluster(cluster); err != nil {
		return nil, errors.Wrap(err, "failed to set default cluster")
	}
	if exists {
		if cluster, err = store.StoreProvider.Owner(owner).Clusters().Update(cluster); err != nil {
			return nil, errors.Wrap(err, "failed to update cluster object")
		}
	} else {
		if err = CreateCACertificates(store.StoreProvider, cluster, owner); err != nil {
			return nil, err
		}
		if err = CreateServiceAccountKey(store.StoreProvider, cluster, owner); err != nil {
			return nil, err
		}
		if err = CreateEtcdCertificates(store.StoreProvider, cluster, owner); err != nil {
			return nil, err
		}
		if err = CreateSSHKey(store.StoreProvider, cluster, owner); err != nil {
			return nil, err
		}

		if cluster, err = store.StoreProvider.Owner(owner).Clusters().Create(cluster); err != nil {
			return nil, errors.Wrap(err, "failed to create cluster object in store")
		}
	}

	if !managedProviders.Has(cluster.ClusterConfig().Cloud.CloudProvider) {
		for i := 0; i < cluster.Spec.Config.MasterCount; i++ {
			master, err := CreateMasterMachines(cluster, i)
			if err != nil {
				return nil, err
			}
			if _, err = store.StoreProvider.Owner(owner).Machine(cluster.Name).Create(master); err != nil {
				return nil, err
			}
		}
	}

	return store.StoreProvider.Owner(owner).Clusters().Update(cluster)
}

func SetDefaultCluster(cluster *api.Cluster) error {
	config := cluster.Spec.Config

	if err := api.AssignTypeKind(cluster); err != nil {
		return err
	}
	if err := api.AssignTypeKind(cluster.Spec.ClusterAPI); err != nil {
		return err
	}

	cluster.SetNetworkingDefaults(config.Cloud.NetworkProvider)

	config.APIServerExtraArgs = map[string]string{
		"kubelet-preferred-address-types": strings.Join([]string{
			string(core.NodeInternalIP),
			string(core.NodeInternalDNS),
			string(core.NodeExternalDNS),
			string(core.NodeExternalIP),
		}, ","),
		"cloud-provider": cluster.Spec.Config.Cloud.CloudProvider,
	}

	cluster.Spec.Config.Cloud.Region = cluster.Spec.Config.Cloud.Zone[0 : len(cluster.Spec.Config.Cloud.Zone)-1]
	cluster.Spec.Config.Cloud.SSHKeyName = cluster.GenSSHKeyExternalID()

	return nil
}

func CreateMasterMachines(cluster *api.Cluster, index int) (*clusterapi.Machine, error) {
	cm, err := GetCloudManager(cluster.ClusterConfig().Cloud.CloudProvider)
	if err != nil {
		return nil, err
	}
	providerSpec, err := cm.GetDefaultMachineProviderSpec(cluster, "", api.MasterRole)
	if err != nil {
		return nil, err
	}

	/*role := api.RoleMember
	if ind == 0 {
		role = api.RoleLeader
	}*/
	machine := &clusterapi.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%v-master-%v", cluster.Name, index),
			//	UID:               uuid.NewUUID(),
			CreationTimestamp: metav1.Time{Time: time.Now()},
			Labels: map[string]string{
				//ref: https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/94a3a3abc7b1ebdd88ea89889347f5e644e160cf/pkg/cloud/aws/actuators/machine_scope.go#L90-L93
				//ref: https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/94a3a3abc7b1ebdd88ea89889347f5e644e160cf/pkg/cloud/aws/actuators/machine/actuator.go#L89-L92
				"set":                              "controlplane",
				api.RoleMasterKey:                  "",
				clusterapi.MachineClusterLabelName: cluster.Name,
			},
		},
		Spec: clusterapi.MachineSpec{
			ProviderSpec: providerSpec,
			Versions: clusterapi.MachineVersionInfo{
				Kubelet:      cluster.ClusterConfig().KubernetesVersion,
				ControlPlane: cluster.ClusterConfig().KubernetesVersion,
			},
		},
	}
	if err := api.AssignTypeKind(machine); err != nil {
		return nil, err
	}

	return machine, nil
}

func CreateMachineSet(cluster *api.Cluster, owner, role, sku string, nodeType api.NodeType, count int32, spotPriceMax float64) error {
	var err error
	cm, err := GetCloudManager(cluster.ClusterConfig().Cloud.CloudProvider)
	if err != nil {
		return err
	}

	providerSpec, err := cm.GetDefaultMachineProviderSpec(cluster, sku, api.NodeRole)
	if err != nil {
		return err
	}

	machineSet := clusterapi.MachineSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:              strings.Replace(strings.ToLower(sku), "_", "-", -1) + "-pool",
			CreationTimestamp: metav1.Time{Time: time.Now()},
		},
		Spec: clusterapi.MachineSetSpec{
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					clusterapi.MachineClusterLabelName: cluster.Name,
					api.MachineSlecetor:                sku,
				},
			},
			Replicas: Int32P(count),
			Template: clusterapi.MachineTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						api.PharmerCluster:                 cluster.Name,
						api.RoleNodeKey:                    "",
						api.MachineSlecetor:                sku,
						"set":                              "node",
						clusterapi.MachineClusterLabelName: cluster.Name, //ref:https://github.com/kubernetes-sigs/cluster-api/blob/master/pkg/controller/machine/controller.go#L229-L232
					},
					CreationTimestamp: metav1.Time{Time: time.Now()},
				},
				Spec: clusterapi.MachineSpec{
					ProviderSpec: providerSpec,
					Versions: clusterapi.MachineVersionInfo{
						Kubelet: cluster.ClusterConfig().KubernetesVersion,
					},
				},
			},
		},
	}

	_, err = store.StoreProvider.Owner(owner).MachineSet(cluster.Name).Create(&machineSet)

	return err
}

func Delete(ctx context.Context, name string, owner string) (*api.Cluster, error) {
	if name == "" {
		return nil, errors.New("missing cluster name")
	}

	cluster, err := Store(ctx).Owner(owner).Clusters().Get(name)
	if err != nil {
		return nil, errors.Errorf("cluster `%s` does not exist. Reason: %v", name, err)
	}
	cluster.DeletionTimestamp = &metav1.Time{Time: time.Now()}
	cluster.Status.Phase = api.ClusterDeleting

	return Store(ctx).Owner(owner).Clusters().Update(cluster)
}

func DeleteMachineSet(ctx context.Context, clusterName, setName, owner string) error {
	if clusterName == "" {
		return errors.New("missing cluster name")
	}
	if setName == "" {
		return errors.New("missing machineset name")
	}

	mSet, err := Store(ctx).Owner(owner).MachineSet(clusterName).Get(setName)
	if err != nil {
		return errors.Errorf(`machinset not found in pharmer db, try using kubectl`)
	}
	tm := metav1.Now()
	mSet.DeletionTimestamp = &tm
	_, err = Store(ctx).Owner(owner).MachineSet(clusterName).Update(mSet)
	return err
}

func GetSSHConfig(ctx context.Context, owner, nodeName string, cluster *api.Cluster) (*api.SSHConfig, error) {
	var err error
	ctx, err = LoadCACertificates(ctx, cluster, owner)
	if err != nil {
		return nil, err
	}
	client, err := NewAdminClient(ctx, cluster)
	if err != nil {
		return nil, err
	}
	node, err := client.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	ctx, err = LoadSSHKey(ctx, cluster, owner)
	if err != nil {
		return nil, err
	}

	cm, err := GetCloudManager(cluster.ClusterConfig().Cloud.CloudProvider)
	if err != nil {
		return nil, err
	}
	return cm.GetSSHConfig(cluster, node)
}

func GetAdminConfig(ctx context.Context, cluster *api.Cluster, owner string) (*api.KubeConfig, error) {
	if managedProviders.Has(cluster.ClusterConfig().Cloud.CloudProvider) {
		cm, err := GetCloudManager(cluster.ClusterConfig().Cloud.CloudProvider)
		if err != nil {
			return nil, err
		}
		cm.SetOwner(owner)
		return cm.GetKubeConfig(cluster)
	}
	var err error
	ctx, err = LoadCACertificates(ctx, cluster, owner)
	if err != nil {
		return nil, err
	}
	adminCert, adminKey, err := CreateAdminCertificate(ctx)
	if err != nil {
		return nil, err
	}

	var (
		clusterName = fmt.Sprintf("%s.pharmer", cluster.Name)
		userName    = fmt.Sprintf("cluster-admin@%s.pharmer", cluster.Name)
		ctxName     = fmt.Sprintf("cluster-admin@%s.pharmer", cluster.Name)
	)
	cfg := api.KubeConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "KubeConfig",
		},
		Preferences: api.Preferences{
			Colors: true,
		},
		Cluster: api.NamedCluster{
			Name:                     clusterName,
			Server:                   cluster.APIServerURL(),
			CertificateAuthorityData: cert.EncodeCertPEM(CACert(ctx)),
		},
		AuthInfo: api.NamedAuthInfo{
			Name:                  userName,
			ClientCertificateData: cert.EncodeCertPEM(adminCert),
			ClientKeyData:         cert.EncodePrivateKeyPEM(adminKey),
		},
		Context: api.NamedContext{
			Name:     ctxName,
			Cluster:  clusterName,
			AuthInfo: userName,
		},
	}
	return &cfg, nil
}

func Apply(ctx context.Context, opts *options.ApplyConfig) ([]api.Action, error) {
	if opts.ClusterName == "" {
		return nil, errors.New("missing cluster name")
	}

	cluster, err := Store(ctx).Owner(opts.Owner).Clusters().Get(opts.ClusterName)
	if err != nil {
		return nil, errors.Errorf("cluster `%s` does not exist. Reason: %v", opts.ClusterName, err)
	}

	cm, err := GetCloudManager(cluster.ClusterConfig().Cloud.CloudProvider)
	if err != nil {
		return nil, err
	}
	cm.SetOwner(opts.Owner)

	return cm.Apply(cluster, opts.DryRun)
}

func CheckForUpdates(ctx context.Context, name, owner string) (string, error) {
	if name == "" {
		return "", errors.New("missing cluster name")
	}

	cluster, err := Store(ctx).Owner(owner).Clusters().Get(name)
	if err != nil {
		return "", errors.Errorf("cluster `%s` does not exist. Reason: %v", name, err)
	}
	if cluster.Status.Phase == "" {
		return "", errors.Errorf("cluster `%s` is in unknown phase", cluster.Name)
	}
	if cluster.Status.Phase != api.ClusterReady {
		return "", errors.Errorf("cluster `%s` is not ready", cluster.Name)
	}
	if cluster.Status.Phase == api.ClusterDeleted {
		return "", nil
	}
	if ctx, err = LoadCACertificates(ctx, cluster, owner); err != nil {
		return "", err
	}
	if ctx, err = LoadSSHKey(ctx, cluster, owner); err != nil {
		return "", err
	}
	kc, err := NewAdminClient(ctx, cluster)
	if err != nil {
		return "", err
	}

	upm := NewUpgradeManager(ctx, kc, cluster, owner)
	upgrades, err := upm.GetAvailableUpgrades()
	if err != nil {
		return "", err
	}
	upm.PrintAvailableUpgrades(upgrades)
	return "", nil
}

func UpdateSpec(ctx context.Context, cluster *api.Cluster, owner string) (*api.Cluster, error) {
	if cluster == nil {
		return nil, errors.New("missing cluster")
	} else if cluster.Name == "" {
		return nil, errors.New("missing cluster name")
	} else if cluster.ClusterConfig().KubernetesVersion == "" {
		return nil, errors.New("missing cluster version")
	}

	existing, err := Store(ctx).Owner(owner).Clusters().Get(cluster.Name)
	if err != nil {
		return nil, errors.Errorf("cluster `%s` does not exist. Reason: %v", cluster.Name, err)
	}
	cluster.Status = existing.Status
	cluster.Generation = time.Now().UnixNano()

	return Store(ctx).Owner(owner).Clusters().Update(cluster)
}

func GetBooststrapClient(ctx context.Context, cluster *api.Cluster, owner string) (clusterclient.Client, error) {
	clientFactory := clusterclient.NewFactory()
	kubeConifg, err := GetAdminConfig(ctx, cluster, owner)
	if err != nil {
		return nil, err
	}

	config := api.Convert_KubeConfig_To_Config(kubeConifg)
	data, err := clientcmd.Write(*config)
	bootstrapClient, err := clientFactory.NewClientFromKubeconfig(string(data))
	if err != nil {
		return nil, fmt.Errorf("unable to create bootstrap client: %v", err)
	}
	return bootstrapClient, nil
}

func GetKubernetesClient(ctx context.Context, cluster *api.Cluster, owner string) (kubernetes.Interface, error) {
	kubeConifg, err := GetAdminConfig(ctx, cluster, owner)
	if err != nil {
		return nil, err
	}

	config := api.NewRestConfig(kubeConifg)

	return kubernetes.NewForConfig(config)
}

func GetLeaderMachine(ctx context.Context, cluster *v1beta1.Cluster, owner string) (*clusterapi.Machine, error) {
	machine, err := Store(ctx).Owner(owner).Machine(cluster.Name).Get(cluster.Name + "-master-0")
	if err != nil {
		return nil, err
	}
	return machine, nil
}

/*func GetSSHConfig(ctx context.Context, hostip string) *api.SSHConfig {
	return &api.SSHConfig{
		PrivateKey: SSHKey(ctx).PrivateKey,
		User:       "root",
		HostPort:   int32(22),
		HostIP: hostip,
	}
}*/

// DeleteAllWorkerMachines waits for all nodes to be deleted
func DeleteAllWorkerMachines(ctx context.Context, cluster *v1beta1.Cluster, owner string) error {
	log.Infof("Deleting non-controlplane machines")

	client, err := GetBooststrapClient(ctx, cluster, owner)
	if err != nil {
		return errors.Wrap(err, "failed to get clusterapi client")
	}

	err = deleteMachineDeployments(client)
	if err != nil {
		log.Infof("failed to delete machine deployments: %v", err)
	}

	err = deleteMachineSets(client)
	if err != nil {
		log.Infof("failed to delete machinesetes: %v", err)
	}

	err = deleteMachines(client)
	if err != nil {
		log.Infof("failed to delete machines: %v", err)
	}

	log.Infof("successfully deleted non-controlplane machines")
	return nil
}

// deletes machinedeployments in all namespaces
func deleteMachineDeployments(client clusterclient.Client) error {
	err := client.DeleteMachineDeployments(corev1.NamespaceAll)
	if err != nil {
		return err
	}

	return wait.PollImmediate(RetryInterval, RetryTimeout, func() (done bool, err error) {
		deployList, err := client.GetMachineDeployments(corev1.NamespaceAll)
		if err != nil {
			log.Infof("failed to list machine deployments: %v", err)
			return false, nil
		}
		if len(deployList) == 0 {
			log.Infof("successfully deleted machine deployments")
			return true, nil
		}
		log.Infof("machine deployments are not deleted yet")
		return false, nil
	})
}

// deletes machinesets in all namespaces
func deleteMachineSets(client clusterclient.Client) error {
	err := client.DeleteMachineSets(corev1.NamespaceAll)
	if err != nil {
		return err
	}

	return wait.PollImmediate(RetryInterval, RetryTimeout, func() (done bool, err error) {
		machineSetList, err := client.GetMachineSets(corev1.NamespaceAll)
		if err != nil {
			log.Infof("failed to list machine sets: %v", err)
			return false, nil
		}
		if len(machineSetList) == 0 {
			log.Infof("successfully deleted machinesets")
			return true, nil
		}
		log.Infof("machinesets are not deleted yet")
		return false, nil
	})
}

// deletes machines in all namespaces
func deleteMachines(client clusterclient.Client) error {
	// delete non-controlplane machines
	machineList, err := client.GetMachines(corev1.NamespaceAll)
	for _, machine := range machineList {
		if !util.IsControlPlaneMachine(machine) && machine.DeletionTimestamp == nil {
			err = client.ForceDeleteMachine(machine.Namespace, machine.Name)
			if err != nil {
				log.Infof("failed to delete machine %s in namespace %s", machine.Namespace, machine.Name)
			}
		}
	}

	// wait for machines to be deleted
	return wait.PollImmediate(RetryInterval, RetryTimeout, func() (done bool, err error) {
		machineList, err := client.GetMachines(corev1.NamespaceAll)
		for _, machine := range machineList {
			if !util.IsControlPlaneMachine(machine) {
				log.Infof("machine %s in namespace %s is not deleted yet", machine.Name, machine.Namespace)
			}
		}

		log.Infof("successfully deleted non-controlplane machines")
		return true, nil
	})
}
