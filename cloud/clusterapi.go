package cloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/appscode/go/log"
	"github.com/appscode/go/wait"
	api "github.com/pharmer/pharmer/apis/v1beta1"
	"github.com/pkg/errors"
	"gomodules.xyz/cert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	kubeadmconsts "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/clusterdeployer/clusterclient"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/phases"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"
)

type ClusterApi struct {
	ctx                 context.Context
	cluster             *api.Cluster
	pharmerCertificates *api.PharmerCertificates
	namespace           string
	token               string
	kc                  kubernetes.Interface

	providerComponenet ClusterApiProviderComponent

	clusterapiClient clientset.Interface
	bootstrapClient  clusterclient.Client
}

type ApiServerTemplate struct {
	ClusterName         string
	Provider            string
	ControllerNamespace string
	ControllerImage     string
}

var MachineControllerImage = "pharmer/machine-controller:0.3.0"

const (
	BasePath = ".pharmer/config.d"
)

func NewClusterApi(cm Interface, cluster *api.Cluster, namespace string, kc kubernetes.Interface, pc ClusterApiProviderComponent) (*ClusterApi, error) {
	var token string
	var err error
	if token, err = GetExistingKubeadmToken(kc, kubeadmconsts.DefaultTokenDuration); err != nil {
		return nil, err
	}

	bc, err := GetBooststrapClient(cm, cluster)
	if err != nil {
		return nil, err
	}
	clusterClient, err := GetClusterClient(cm, cluster)
	if err != nil {
		return nil, err
	}

	pharmerCertificates := cm.GetPharmerCertificates()

	return &ClusterApi{
		cluster:             cluster,
		namespace:           namespace,
		pharmerCertificates: pharmerCertificates,
		kc:                  kc,
		clusterapiClient:    clusterClient,
		providerComponenet:  pc,
		token:               token,
		bootstrapClient:     bc}, nil
}

func GetClusterClient(cm Interface, cluster *api.Cluster) (clientset.Interface, error) {
	conf, err := NewRestConfig(cm, cluster)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get rest config")
	}
	return clientset.NewForConfig(conf)
}

func (ca *ClusterApi) Apply(controllerManager string) error {
	log.Infof("Deploying the addon apiserver and controller manager...")
	if err := ca.CreateMachineController(controllerManager); err != nil {
		return errors.Wrap(err, "can't create machine controller")
	}

	if err := phases.ApplyCluster(ca.bootstrapClient, ca.cluster.Spec.ClusterAPI); err != nil && !api.ErrAlreadyExist(err) {
		return errors.Wrap(err, "failed to add cluster")
	}
	namespace := ca.cluster.Spec.ClusterAPI.Namespace
	if namespace == "" {
		namespace = ca.bootstrapClient.GetContextNamespace()
	}

	c, err := ca.clusterapiClient.ClusterV1alpha1().Clusters(namespace).Get(ca.cluster.Name, metav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to update cluster provider status")
	}

	c.Status = ca.cluster.Spec.ClusterAPI.Status
	if _, err := ca.clusterapiClient.ClusterV1alpha1().Clusters(namespace).UpdateStatus(c); err != nil && !api.ErrObjectModified(err) {
		return errors.Wrap(err, "failed to update cluster")
	}

	if err := ca.updateProviderStatus(); err != nil {
		log.Infoln(err)
		return errors.Wrap(err, "failed to update provider status")
	}

	masterMachine, err := GetLeaderMachine(ca.cluster)
	if err != nil {
		return errors.Wrap(err, "failed to get leader machine")
	}

	masterMachine.Annotations = make(map[string]string)
	masterMachine.Annotations[InstanceStatusAnnotationKey] = ""

	log.Infof("Adding master machines...")
	err = phases.ApplyMachines(ca.bootstrapClient, namespace, []*clusterv1.Machine{masterMachine})
	if err != nil && !api.ErrAlreadyExist(err) && !api.ErrObjectModified(err) {
		return errors.Wrap(err, "failed to add master machine")
	}

	// get the machine object and update the provider status field
	err = ca.updateMachineStatus(namespace, masterMachine)
	if err != nil {
		return errors.Wrap(err, "failed to update machine status")
	}

	return nil
}

func (ca *ClusterApi) updateProviderStatus() error {
	return wait.PollImmediate(RetryInterval, RetryTimeout, func() (bool, error) {
		cluster, err := ca.clusterapiClient.ClusterV1alpha1().Clusters(ca.cluster.Spec.ClusterAPI.Namespace).Get(ca.cluster.Name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}
		if cluster.Status.ProviderStatus != nil {
			ca.cluster.Spec.ClusterAPI.Status.ProviderStatus = cluster.Status.ProviderStatus
			if _, err := Store(ca.ctx).Clusters().Update(ca.cluster); err != nil {
				return false, nil
			}
			return true, nil
		}
		return false, nil
	})
}

func (ca *ClusterApi) updateMachineStatus(namespace string, masterMachine *clusterv1.Machine) error {
	return wait.PollImmediate(RetryInterval, RetryTimeout, func() (bool, error) {
		m, err := ca.clusterapiClient.ClusterV1alpha1().Machines(namespace).Get(masterMachine.Name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}
		m.Status.ProviderStatus = masterMachine.Status.ProviderStatus
		if _, err := ca.clusterapiClient.ClusterV1alpha1().Machines(namespace).UpdateStatus(m); err != nil {
			return false, nil
		}
		return true, nil
	})
}

func (ca *ClusterApi) CreateMachineController(controllerManager string) error {
	log.Infoln("creating pharmer secret")
	if err := ca.CreatePharmerSecret(); err != nil {
		return err
	}

	log.Infoln("creating apiserver and controller")
	if err := ca.CreateApiServerAndController(controllerManager); err != nil && !api.ErrObjectModified(err) {
		return err
	}
	return nil
}

func (ca *ClusterApi) CreatePharmerSecret() error {
	providerConfig := ca.cluster.ClusterConfig()

	cred, err := Store(ca.ctx).Credentials().Get(ca.cluster.ClusterConfig().CredentialName)
	if err != nil {
		return err
	}
	credData, err := json.MarshalIndent(cred, "", "  ")
	if err != nil {
		return err
	}

	if err = CreateNamespace(ca.kc, ca.namespace); err != nil {
		return err
	}

	if err = CreateSecret(ca.kc, "pharmer-cred", ca.namespace, map[string][]byte{
		fmt.Sprintf("%v.json", ca.cluster.ClusterConfig().CredentialName): credData,
	}); err != nil {
		return err
	}

	if ca.providerComponenet != nil {
		if err = ca.providerComponenet.CreateCredentialSecret(ca.kc, cred.Spec.Data); err != nil {
			return err
		}
	}

	cluster, err := json.MarshalIndent(ca.cluster, "", "  ")
	if err != nil {
		return err
	}
	if err = CreateSecret(ca.kc, "pharmer-cluster", ca.namespace, map[string][]byte{
		fmt.Sprintf("%v.json", ca.cluster.Name): cluster,
	}); err != nil {
		return err
	}

	publicKey, privateKey, err := Store(ca.ctx).SSHKeys(ca.cluster.Name).Get(ca.cluster.ClusterConfig().Cloud.SSHKeyName)
	if err != nil {
		return err
	}
	if err = CreateSecret(ca.kc, "pharmer-ssh", ca.namespace, map[string][]byte{
		fmt.Sprintf("id_%v", providerConfig.Cloud.SSHKeyName):     privateKey,
		fmt.Sprintf("id_%v.pub", providerConfig.Cloud.SSHKeyName): publicKey,
	}); err != nil {
		return err
	}

	if err = CreateSecret(ca.kc, "pharmer-certificate", ca.namespace, map[string][]byte{
		"ca.crt":             cert.EncodeCertPEM(ca.pharmerCertificates.CACert.Cert),
		"ca.key":             cert.EncodePrivateKeyPEM(ca.pharmerCertificates.CACert.Key),
		"front-proxy-ca.crt": cert.EncodeCertPEM(ca.pharmerCertificates.FrontProxyCACert.Cert),
		"front-proxy-ca.key": cert.EncodePrivateKeyPEM(ca.pharmerCertificates.FrontProxyCACert.Key),
		"sa.crt":             cert.EncodeCertPEM(ca.pharmerCertificates.ServiceAccountCert.Cert),
		"sa.key":             cert.EncodePrivateKeyPEM(ca.pharmerCertificates.ServiceAccountCert.Key),
	}); err != nil {
		return err
	}

	if err = CreateSecret(ca.kc, "pharmer-etcd", ca.namespace, map[string][]byte{
		"ca.crt": cert.EncodeCertPEM(ca.pharmerCertificates.EtcdCACert.Cert),
		"ca.key": cert.EncodePrivateKeyPEM(ca.pharmerCertificates.EtcdCACert.Key),
	}); err != nil {
		return err
	}

	return nil
}

func (ca *ClusterApi) CreateApiServerAndController(controllerManager string) error {
	tmpl, err := template.New("config").Parse(controllerManager)
	if err != nil {
		return err
	}
	var tmplBuf bytes.Buffer
	err = tmpl.Execute(&tmplBuf, ApiServerTemplate{
		ClusterName:         ca.cluster.Name,
		Provider:            ca.cluster.ClusterConfig().Cloud.CloudProvider,
		ControllerNamespace: ca.namespace,
		ControllerImage:     MachineControllerImage,
	})
	if err != nil {
		return err
	}

	return ca.bootstrapClient.Apply(tmplBuf.String())
}

func (ca *ClusterApi) GetMachines() (*[]clusterv1.Machine, error) {
	machineList, err := ca.clusterapiClient.ClusterV1alpha1().Machines("default").List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return &machineList.Items, err
}

func (ca *ClusterApi) UpdateMachine(machine *clusterv1.Machine) error {
	_, err := ca.clusterapiClient.ClusterV1alpha1().Machines("default").Update(machine)

	return err
}
