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
package vfs

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	api "pharmer.dev/pharmer/apis/v1alpha1"
	"pharmer.dev/pharmer/store"

	"github.com/pkg/errors"
	"gomodules.xyz/stow"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

type machineFileStore struct {
	container stow.Container
	prefix    string
	cluster   string
	owner     string
}

var _ store.MachineStore = &machineFileStore{}

func (s *machineFileStore) resourceHome() string {
	return filepath.Join(s.owner, s.prefix, "clusters", s.cluster, "machine")
}

func (s *machineFileStore) resourceID(name string) string {
	return filepath.Join(s.resourceHome(), name+".json")
}

func (s *machineFileStore) List(opts metav1.ListOptions) ([]*clusterapi.Machine, error) {
	result := make([]*clusterapi.Machine, 0)
	cursor := stow.CursorStart
	for {
		page, err := s.container.Browse(s.resourceHome()+string(os.PathSeparator), string(os.PathSeparator), cursor, pageSize)
		if err != nil {
			return nil, errors.Errorf("failed to list node groups. Reason: %v", err)
		}
		for _, item := range page.Items {
			r, err := item.Open()
			if err != nil {
				return nil, errors.Errorf("failed to list node groups. Reason: %v", err)
			}
			var obj clusterapi.Machine
			err = json.NewDecoder(r).Decode(&obj)
			if err != nil {
				return nil, errors.Errorf("failed to list node groups. Reason: %v", err)
			}
			result = append(result, &obj)
			r.Close()
		}
		cursor = page.Cursor
		if stow.IsCursorEnd(cursor) {
			break
		}
	}
	return result, nil
}

func (s *machineFileStore) Get(name string) (*clusterapi.Machine, error) {
	if s.cluster == "" {
		return nil, errors.New("missing cluster name")
	}
	if name == "" {
		return nil, errors.New("missing node group name")
	}

	item, err := s.container.Item(s.resourceID(name))
	if err != nil {
		return nil, errors.Errorf("NodeGroup `%s` does not exist. Reason: %v", name, err)
	}

	r, err := item.Open()
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var existing clusterapi.Machine
	err = json.NewDecoder(r).Decode(&existing)
	if err != nil {
		return nil, err
	}
	return &existing, nil
}

func (s *machineFileStore) Create(obj *clusterapi.Machine) (*clusterapi.Machine, error) {
	if s.cluster == "" {
		return nil, errors.New("missing cluster name")
	}
	if obj == nil {
		return nil, errors.New("missing node group")
	} else if obj.Name == "" {
		return nil, errors.New("missing node group name")
	}
	err := api.AssignTypeKind(obj)
	if err != nil {
		return nil, err
	}

	id := s.resourceID(obj.Name)
	_, err = s.container.Item(id)
	if err == nil {
		return nil, errors.Errorf("NodeGroup `%s` already exists", obj.Name)
	}

	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return nil, err
	}
	_, err = s.container.Put(id, bytes.NewBuffer(data), int64(len(data)), nil)
	return obj, err
}

func (s *machineFileStore) Update(obj *clusterapi.Machine) (*clusterapi.Machine, error) {
	if s.cluster == "" {
		return nil, errors.New("missing cluster name")
	}
	if obj == nil {
		return nil, errors.New("missing node group")
	} else if obj.Name == "" {
		return nil, errors.New("missing node group name")
	}
	err := api.AssignTypeKind(obj)
	if err != nil {
		return nil, err
	}

	id := s.resourceID(obj.Name)

	_, err = s.container.Item(id)
	if err != nil {
		return nil, errors.Errorf("NodeGroup `%s` does not exist. Reason: %v", obj.Name, err)
	}

	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return nil, err
	}
	_, err = s.container.Put(id, bytes.NewBuffer(data), int64(len(data)), nil)
	return obj, err
}

func (s *machineFileStore) Delete(name string) error {
	if s.cluster == "" {
		return errors.New("missing cluster name")
	}
	if name == "" {
		return errors.New("missing node group name")
	}
	return s.container.RemoveItem(s.resourceID(name))
}

func (s *machineFileStore) UpdateStatus(obj *clusterapi.Machine) (*clusterapi.Machine, error) {
	if s.cluster == "" {
		return nil, errors.New("missing cluster name")
	}
	if obj == nil {
		return nil, errors.New("missing node group")
	} else if obj.Name == "" {
		return nil, errors.New("missing node group name")
	}
	err := api.AssignTypeKind(obj)
	if err != nil {
		return nil, err
	}

	id := s.resourceID(obj.Name)

	item, err := s.container.Item(id)
	if err != nil {
		return nil, errors.Errorf("NodeGroup `%s` does not exist. Reason: %v", obj.Name, err)
	}

	r, err := item.Open()
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var existing clusterapi.Machine
	err = json.NewDecoder(r).Decode(&existing)
	if err != nil {
		return nil, err
	}
	existing.Status = obj.Status

	data, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		return nil, err
	}
	_, err = s.container.Put(id, bytes.NewBuffer(data), int64(len(data)), nil)
	return &existing, err
}
