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

package main

import (
	clientset "github.com/nimrodshn/kubechain/pkg/clientset/v1alpha1"

	"flag"
	"github.com/nimrodshn/kubechain/pkg/controllers/blockchain"
	"github.com/nimrodshn/kubechain/pkg/types/v1alpha1"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"

	"log"
)

var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to Kubernetes config file")
	flag.Parse()
}

func main() {
	var config *rest.Config
	var err error

	if kubeconfig == "" {
		log.Printf("using in-cluster configuration")
		config, err = rest.InClusterConfig()
	} else {
		log.Printf("using configuration from '%s'", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		panic(err)
	}

	v1alpha1.AddToScheme(scheme.Scheme)

	client, err := clientset.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	informer, indexer := blockchain.NewInformer(v1.NamespaceDefault, client)
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	controller := blockchain.NewController(
		queue,
		indexer,
		informer,
	)

	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(1, stop)

	select {}
}
