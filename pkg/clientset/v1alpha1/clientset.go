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
	"github.com/nimrodshn/kubechain/pkg/types/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

// KubechainV1Alpha1Interface is an entrypoint for our client.
type KubechainV1Alpha1Interface interface {
	Block(namespace string) BlockInterface
}

// KubechainV1Alpha1Client implements KubechainV1Alpha1Interface
// and is the entrypoint for all CRUD operations on the "Block" resource.
type KubechainV1Alpha1Client struct {
	restClient rest.Interface
}

// NewForConfig creates a new client implementing KubechainV1Alpha1Interface
func NewForConfig(c *rest.Config) (*KubechainV1Alpha1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1alpha1.GroupName, Version: v1alpha1.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &KubechainV1Alpha1Client{restClient: client}, nil
}

// Block creates a returns a client adhering to the BlockInterface. (see block.go)
func (c *KubechainV1Alpha1Client) Block(namespace string) BlockInterface {
	return &blockClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
