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

type Blockchain struct {
	Chain []*Block
}

// AddBlock adds a new block to the blockchain.
func (bc *Blockchain) AddBlock(data string) {
	if len(bc.Chain) == 0 {
		genesis := NewBlock(data, []byte{})
		chain := []*Block{genesis}
		bc.Chain = chain
	}
	prevBlock := bc.Chain[len(bc.Chain)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.Chain = append(bc.Chain, newBlock)
}

// NewBlockchain is a constructor for our blockchain struct.
func NewBlockchain() *Blockchain {
	blockchain := new(Blockchain)
	blockchain.AddBlock("New Genesis Block")
	return blockchain
}
