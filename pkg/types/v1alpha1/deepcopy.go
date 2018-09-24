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
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto copies infromation from one (pointer of) block to another.
func (in *Block) DeepCopyInto(out *Block) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta

	out.Timestamp = in.Timestamp
	out.Data = in.Data
	out.PrevBlockHash = in.PrevBlockHash
	out.Hash = in.Hash
	out.Nonce = in.Nonce
}

// DeepCopyObject returns a generically typed copy of an object
func (in *Block) DeepCopyObject() runtime.Object {
	out := Block{}
	in.DeepCopyInto(&out)
	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *BlockList) DeepCopyObject() runtime.Object {
	out := BlockList{}
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Items = make([]Block, len(in.Items))
	for idx := range in.Items {
		out.Items[idx].DeepCopyInto(&out.Items[idx])
	}
	return &out
}
