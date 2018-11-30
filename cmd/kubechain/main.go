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
	v1alpha1 "github.com/nimrodshn/kubechain/pkg/types/v1alpha1"

	"github.com/nimrodshn/kubechain/pkg/controllers/blockchain"
	"k8s.io/apimachinery/pkg/util/wait"

	"flag"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"

	"log"
)

var kubeconfig string

// The number of threads to process events.
const threadCount = 3

// The namespace to run our controller on.
const defaultNamespace = "default"

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

	err = v1alpha1.AddToScheme(scheme.Scheme)
	if err != nil {
		panic(err)
	}

	client, err := clientset.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// Create the queue for block events.
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// Create the informer which has a cache of all the blocks representing the current system.
	// This cache is the used by the informer to react to create/update/delete block events which are then passed to the queue
	// to be processed.
	informer := blockchain.NewInformer(defaultNamespace, client, queue)

	// Construct our controller from the given queue and informers and a new blockcahin.
	controller := blockchain.NewController(queue, informer, new(v1alpha1.Blockchain))

	controller.Run(threadCount, wait.NeverStop)

}
