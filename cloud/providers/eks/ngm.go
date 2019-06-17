package eks

import (
	"context"
	"fmt"
	"strings"

	"github.com/appscode/go/log"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/ghodss/yaml"
	. "github.com/pharmer/pharmer/cloud"
	"github.com/pharmer/pharmer/store"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	core_util "kmodules.xyz/client-go/core/v1"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

type EKSNodeGroupManager struct {
	ctx  context.Context
	conn *cloudConnector
	ng   *clusterapi.MachineSet
	kc   kubernetes.Interface

	owner string
}

func NewEKSNodeGroupManager(ctx context.Context, conn *cloudConnector, ng *clusterapi.MachineSet, kc kubernetes.Interface, owner string) *EKSNodeGroupManager {
	return &EKSNodeGroupManager{ctx: ctx, conn: conn, ng: ng, kc: kc, owner: owner}
}

func (igm *EKSNodeGroupManager) Apply() error {
	fileName := igm.ng.Name
	igm.ng.Name = strings.Replace(igm.ng.Name, ".", "-", -1)
	//var template []byte
	var found bool
	var err error

	if found, err = igm.conn.isStackExists(igm.ng.Name); err != nil {
		return err
	}

	if !found {
		params := igm.buildstackParams()
		if err = igm.conn.createStack(igm.ng.Name, NodeGroupUrl, params, true); err != nil {
			return err
		}
		var ngInfo *cloudformation.Stack
		ngInfo, err = igm.conn.getStack(igm.ng.Name)
		if err != nil {
			return err
		}
		if err = igm.newNodeAuthConfigMap(igm.conn.getOutput(ngInfo, "NodeInstanceRole")); err != nil {
			return err
		}
	} else {
		if *igm.ng.Spec.Replicas == 0 || igm.ng.DeletionTimestamp != nil {

			var ngInfo *cloudformation.Stack
			ngInfo, err = igm.conn.getStack(igm.ng.Name)
			if err != nil {
				return err
			}
			if err = igm.conn.deleteStack(igm.ng.Name); err != nil {
				return err
			}
			if err = igm.deleteNodeAuthConfigMap(igm.conn.getOutput(ngInfo, "NodeInstanceRole")); err != nil {
				return err
			}
			err = store.StoreProvider.MachineSet(igm.conn.cluster.Name).Delete(fileName)
			if err != nil {
				return err
			}
			return

		} else {

			existingStack, err := igm.conn.getStack(igm.ng.Name)
			if err != nil {
				return err
			}
			params := igm.buildstackParams()
			if err = igm.conn.updateStack(igm.ng.Name, params, true, igm.conn.getOutput(existingStack, "NodeInstanceRole")); err != nil {
				log.Infoln(err)
			}

		}
	}
	igm.ng.Status.Replicas = *igm.ng.Spec.Replicas
	_, err = store.StoreProvider.MachineSet(igm.conn.cluster.Name).UpdateStatus(igm.ng)

	return err
}

func (igm *EKSNodeGroupManager) buildstackParams() map[string]string {
	providerSpec := igm.conn.cluster.EKSProviderConfig(igm.ng.Spec.Template.Spec.ProviderSpec.Value.Raw)
	return map[string]string{
		"ClusterName":                         igm.conn.cluster.Name,
		"NodeGroupName":                       igm.ng.Name,
		"KeyName":                             igm.conn.cluster.Spec.Config.Cloud.SSHKeyName,
		"NodeImageId":                         igm.conn.cluster.Spec.Config.Cloud.InstanceImage,
		"NodeInstanceType":                    providerSpec.InstanceType,
		"NodeAutoScalingGroupDesiredCapacity": fmt.Sprintf("%d", *igm.ng.Spec.Replicas),
		"NodeAutoScalingGroupMinSize":         fmt.Sprintf("%d", *igm.ng.Spec.Replicas),
		"NodeAutoScalingGroupMaxSize":         fmt.Sprintf("%d", *igm.ng.Spec.Replicas),
		"ClusterControlPlaneSecurityGroup":    igm.conn.cluster.Status.Cloud.EKS.SecurityGroup,
		"Subnets":                             igm.conn.cluster.Status.Cloud.EKS.SubnetId,
		"VpcId":                               igm.conn.cluster.Status.Cloud.EKS.VpcId,
	}
}

func (igm *EKSNodeGroupManager) deleteNodeAuthConfigMap(arn *string) error {
	configmaps, err := igm.kc.CoreV1().ConfigMaps(metav1.NamespaceSystem).Get(EKSNodeConfigMap, metav1.GetOptions{})
	if err != nil {
		return err
	}
	mapRoles := make([]map[string]interface{}, 0)
	if configmaps != nil {
		existingRules := configmaps.Data[EKSConfigMapRoles]
		if err := yaml.Unmarshal([]byte(existingRules), &mapRoles); err != nil {
			return err
		}
	}
	newRoles := make([]map[string]interface{}, 0)
	for i, r := range mapRoles {
		if r["rolearn"] != *arn {
			newRoles = append(newRoles, mapRoles[i])
			//delete(mapRoles, r)
		}
	}
	mapRolesBytes, err := yaml.Marshal(newRoles)
	if err != nil {
		return err
	}

	_, _, err = core_util.CreateOrPatchConfigMap(igm.kc,
		metav1.ObjectMeta{Namespace: metav1.NamespaceSystem, Name: EKSNodeConfigMap},
		func(in *core.ConfigMap) *core.ConfigMap {
			if in.Data == nil {
				in.Data = make(map[string]string)
			}
			in.Data[EKSConfigMapRoles] = string(mapRolesBytes)
			return in
		})
	return err
}

func (igm *EKSNodeGroupManager) newNodeAuthConfigMap(arn *string) error {
	mapRoles := make([]map[string]interface{}, 1)
	newRole := make(map[string]interface{})

	newRole["rolearn"] = arn
	newRole["username"] = "system:node:{{EC2PrivateDNSName}}"
	newRole["groups"] = []string{
		"system:bootstrappers",
		"system:nodes",
	}

	configmaps, err := igm.kc.CoreV1().ConfigMaps(metav1.NamespaceSystem).Get(EKSNodeConfigMap, metav1.GetOptions{})
	if err == nil {
		existingRules := configmaps.Data[EKSConfigMapRoles]
		if err := yaml.Unmarshal([]byte(existingRules), &mapRoles); err != nil {
			return err
		}
	}
	mapRoles = append(mapRoles, newRole)

	mapRolesBytes, err := yaml.Marshal(mapRoles)
	if err != nil {
		return err
	}

	_, _, err = core_util.CreateOrPatchConfigMap(igm.kc,
		metav1.ObjectMeta{Namespace: metav1.NamespaceSystem, Name: EKSNodeConfigMap},
		func(in *core.ConfigMap) *core.ConfigMap {
			if in.Data == nil {
				in.Data = make(map[string]string)
			}
			in.Data[EKSConfigMapRoles] = string(mapRolesBytes)
			return in
		})
	return err
}
