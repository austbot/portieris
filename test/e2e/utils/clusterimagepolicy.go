// Copyright 2018 IBM
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"testing"

	uuid "github.com/satori/go.uuid"

	"github.com/IBM/portieris/pkg/apis/securityenforcement/v1beta1"
	"github.com/IBM/portieris/test/framework"
	corev1 "k8s.io/api/core/v1"
)

func CheckIfTesting(t *testing.T, boolToCheck bool) {
	if !boolToCheck {
		t.Skip()
	}
}

// DeleteThenReturnClusterImagePolicy is used for temporary deletion of a cluster image policy for a given test
// The returned ClusterImagePolicy should be used to recreate after the test is complete using a defer
func DeleteThenReturnClusterImagePolicy(t *testing.T, fw *framework.Framework, clusterImagePolicy string) *v1beta1.ClusterImagePolicy {
	defaultClusterPolicy, err := fw.GetClusterImagePolicy(clusterImagePolicy)
	if err != nil {
		return nil
	}
	if err := fw.DeleteClusterImagePolicy(clusterImagePolicy); err != nil {
		t.Errorf("error deleting ClusterImagePolicy %q: %v", clusterImagePolicy, err)
	}
	return defaultClusterPolicy
}

func CreateClusterImagePolicyAndNamespace(t *testing.T, fw *framework.Framework, manifestPath string) (*v1beta1.ClusterImagePolicy, *corev1.Namespace) {
	ns := uuid.NewV4().String()
	clusterImagePolicy, err := fw.LoadClusterImagePolicyManifest(manifestPath)
	if err != nil {
		t.Errorf("error loading %q ClusterImagePolicy manifest: %v", clusterImagePolicy.Name, err)
	}
	namespace, err := fw.CreateNamespaceWithIPS(ns)
	if err != nil {
		t.Errorf("error creating %q namespace: %v", ns, err)
	}
	if err := fw.CreateClusterImagePolicy(clusterImagePolicy); err != nil {
		t.Errorf("error creating %q ClusterImagePolicy: %v", clusterImagePolicy.Name, err)
	}
	return clusterImagePolicy, namespace
}

func CleanUpClusterImagePolicyTest(t *testing.T, fw *framework.Framework, clusterPolicy, namespace string) {
	if err := fw.DeleteNamespace(namespace); err != nil {
		t.Logf("failed to delete namespace %q: %v", namespace, err)
	}
	if err := fw.DeleteClusterImagePolicy(clusterPolicy); err != nil {
		t.Logf("failed to delete ClusterImagePolicy %q: %v", clusterPolicy, err)
	}
}
