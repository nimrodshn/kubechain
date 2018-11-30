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
	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Blockchain is our internal blockchain implementation.
type Blockchain struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Chain []*Block `json:"chain"`
}

// AddBlock adds a new block to the blockchain.
func (bc *Blockchain) AddBlock(block *Block) {
	glog.Infof("Adding new block...")
	if len(bc.Chain) == 0 {
		chain := []*Block{block}
		bc.Chain = chain
	} else {
		prevBlock := bc.Chain[len(bc.Chain)-1]
		block.Spec.PrevBlockHash = prevBlock.Spec.Hash
		bc.Chain = append(bc.Chain, block)
	}
}
