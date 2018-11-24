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
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"fmt"
	"time"
)

// Controller is the custom controller for the blockchain CRD.
type Controller struct {
	queue    workqueue.RateLimitingInterface
	informer cache.SharedIndexInformer
}

// NewController is a constructor for the block controller.
func NewController(queue workqueue.RateLimitingInterface, informer cache.SharedIndexInformer) *Controller {
	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				var key string
				var err error
				if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
					runtime.HandleError(err)
					return
				}
				queue.Add(key)
			},
		})
	return &Controller{
		informer: informer,
		queue:    queue,
	}
}

func (c *Controller) processNextItem() bool {
	// Wait until there is a new item in the working queue
	key, quit := c.queue.Get()
	if quit {
		return false
	}

	// Tell the queue that we are done with processing this key. This unblocks the key for other workers
	// This allows safe parallel processing because two pods with the same key are never processed in
	// parallel.
	defer c.queue.Done(key)

	// Invoke the method containing the business logic
	err := addBlockEventHandler(key.(string), c.informer.GetIndexer())
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		c.queue.Forget(key)
		return true
	}
	return true
}

// NewInformer Creates a new informer for the Block crd.
func NewInformer(ns string, clientSet clientset.KubechainV1Alpha1Interface, queue workqueue.RateLimitingInterface) cache.SharedIndexInformer {
	informer := cache.NewSharedIndexInformer(
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
		cache.Indexers{},
	)
	return informer
}

func addBlockEventHandler(key string, indexer cache.Indexer) error {
	item, exists, err := indexer.GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	} else if !exists {
		fmt.Printf("Block %s does not exist anymore\n", key)
	} else {
		block := item.(*v1alpha1.Block)
		block.ProcessNewBlock()
		glog.Infof("Processing new block: %v", block)
	}
	return nil
}

// Run runs the controller
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) {
	defer runtime.HandleCrash()

	// Let the workers stop when we are done
	defer c.queue.ShutDown()

	go c.informer.Run(stopCh)

	// Wait for all involved caches to be synced, before processing items from the queue is started
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}
