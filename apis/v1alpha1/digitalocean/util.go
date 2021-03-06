/*
Copyright The Pharmer Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package digitalocean

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

// ClusterConfigFromProviderSpec unmarshals a provider config into an DigitalOcean Cluster type
func ClusterConfigFromProviderSpec(providerConfig clusterapi.ProviderSpec) (*DigitalOceanClusterProviderSpec, error) {
	var config DigitalOceanClusterProviderSpec
	if providerConfig.Value == nil {
		return &config, nil
	}

	if err := json.Unmarshal(providerConfig.Value.Raw, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// ClusterStatusFromProviderStatus unmarshals a raw extension into an DigitalOcean Cluster type
func ClusterStatusFromProviderStatus(extension *runtime.RawExtension) (*DigitalOceanClusterProviderStatus, error) {
	if extension == nil {
		return &DigitalOceanClusterProviderStatus{}, nil
	}

	status := new(DigitalOceanClusterProviderStatus)
	if err := json.Unmarshal(extension.Raw, status); err != nil {
		return nil, err
	}

	return status, nil
}

// MachineSpecFromProviderSpec unmarshals a raw extension into an DigitalOcean machine type
func MachineConfigFromProviderSpec(providerConfig clusterapi.ProviderSpec) (*DigitalOceanMachineProviderSpec, error) {
	var config DigitalOceanMachineProviderSpec
	if providerConfig.Value == nil {
		return &config, nil
	}

	if err := json.Unmarshal(providerConfig.Value.Raw, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// MachineStatusFromProviderStatus unmarshals a raw extension into an DigitalOcean machine type
func MachineStatusFromProviderStatus(extension *runtime.RawExtension) (*DigitalOceanMachineProviderStatus, error) {
	if extension == nil {
		return &DigitalOceanMachineProviderStatus{}, nil
	}

	status := new(DigitalOceanMachineProviderStatus)
	if err := json.Unmarshal(extension.Raw, status); err != nil {
		return nil, err
	}

	return status, nil
}

// EncodeMachineStatus marshals the machine status
func EncodeMachineStatus(status *DigitalOceanMachineProviderStatus) (*runtime.RawExtension, error) {
	if status == nil {
		return &runtime.RawExtension{}, nil
	}

	var rawBytes []byte
	var err error

	if rawBytes, err = json.Marshal(status); err != nil {
		return nil, err
	}

	return &runtime.RawExtension{
		Raw: rawBytes,
	}, nil
}

// EncodeMachineSpec marshals the machine provider spec.
func EncodeMachineSpec(spec *DigitalOceanMachineProviderSpec) (*runtime.RawExtension, error) {
	if spec == nil {
		return &runtime.RawExtension{}, nil
	}

	var rawBytes []byte
	var err error

	if rawBytes, err = json.Marshal(spec); err != nil {
		return nil, err
	}

	return &runtime.RawExtension{
		Raw: rawBytes,
	}, nil
}

// EncodeClusterStatus marshals the cluster status.
func EncodeClusterStatus(status *DigitalOceanClusterProviderStatus) (*runtime.RawExtension, error) {
	if status == nil {
		return &runtime.RawExtension{}, nil
	}

	var rawBytes []byte
	var err error

	if rawBytes, err = json.Marshal(status); err != nil {
		return nil, err
	}

	return &runtime.RawExtension{
		Raw: rawBytes,
	}, nil
}

// EncodeClusterSpec marshals the cluster provider spec.
func EncodeClusterSpec(spec *DigitalOceanClusterProviderSpec) (*runtime.RawExtension, error) {
	if spec == nil {
		return &runtime.RawExtension{}, nil
	}

	var rawBytes []byte
	var err error

	if rawBytes, err = json.Marshal(spec); err != nil {
		return nil, err
	}

	return &runtime.RawExtension{
		Raw: rawBytes,
	}, nil
}

//func (c *Cluster) SetLinodeProviderConfig(cluster *clusterapi.Cluster, config *ClusterConfig) error {
func SetDigitalOceanClusterProviderConfig(cluster *clusterapi.Cluster) error {
	conf := &DigitalOceanClusterProviderSpec{
		TypeMeta: metav1.TypeMeta{
			APIVersion: DigitalOceanProviderGroupName + "/" + DigitalOceanProviderAPIVersion,
			Kind:       DigitalOceanProviderKind,
		},
	}
	bytes, err := json.Marshal(conf)
	if err != nil {
		return err

	}
	cluster.Spec.ProviderSpec = clusterapi.ProviderSpec{
		Value: &runtime.RawExtension{
			Raw: bytes,
		},
	}
	return nil
}
