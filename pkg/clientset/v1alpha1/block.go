// Copyright 2018 Nimrod Shneor <nimrodshn@gmail.com>
// and other contributors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	scheme "github.com/nimrodshn/kubechain/pkg/clientset/scheme"
	"github.com/nimrodshn/kubechain/pkg/types/v1alpha1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

// BlockInterface is the interface for CRUD actions on blocks
type BlockInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.BlockList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.Block, error)
	Create(*v1alpha1.Block) (*v1alpha1.Block, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

// blockClient implements BlockInterface for the namespace ns.
type blockClient struct {
	restClient rest.Interface
	ns         string
}

func (c *blockClient) List(opts metav1.ListOptions) (*v1alpha1.BlockList, error) {
	result := v1alpha1.BlockList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("blocks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *blockClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.Block, error) {
	result := v1alpha1.Block{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("blocks").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *blockClient) Create(block *v1alpha1.Block) (*v1alpha1.Block, error) {
	result := v1alpha1.Block{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("blocks").
		Body(block).
		Do().
		Into(&result)

	return &result, err
}

func (c *blockClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("blocks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
