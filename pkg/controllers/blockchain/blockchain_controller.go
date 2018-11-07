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

package blockchain

import (
	"github.com/golang/glog"
	clientset "github.com/nimrodshn/kubechain/pkg/clientset/v1alpha1"
	v1alpha1 "github.com/nimrodshn/kubechain/pkg/types/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"fmt"
	"time"
)

// Controller is the custom controller for the blockchain CRD.
type Controller struct {
	queue    workqueue.RateLimitingInterface
	informer cache.Controller
}

// NewController is a constructor for the block controller.
func NewController(queue workqueue.RateLimitingInterface, informer cache.Controller) *Controller {
	return &Controller{
		informer: informer,
		queue:    queue,
	}
}

// NewInformer Creates a new informer for the Block crd.
func NewInformer(ns string, clientSet clientset.KubechainV1Alpha1Interface) cache.Controller {
	_, controller := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result k8sruntime.Object, err error) {
				return clientSet.Block(ns).List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return clientSet.Block(ns).Watch(lo)
			},
		},
		&v1alpha1.Block{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{
			AddFunc: addBlockEventHandler,
		},
	)
	return controller
}

func addBlockEventHandler(obj interface{}) {
	block := obj.(*v1alpha1.Block)
	block.ProcessNewBlock()
	glog.Infof("Processing new block: %v", block)
}

// Run runs the controller
func (c *Controller) Run(stopCh <-chan struct{}) {
	defer runtime.HandleCrash()

	// Let the workers stop when we are done
	defer c.queue.ShutDown()

	go c.informer.Run(stopCh)

	// Wait for all involved caches to be synced, before processing items from the queue is started
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	<-stopCh
}
