# kubechain

kubechain is a small blockchain implementation on top of kubernetes - it includes:
1. A small implementation of a blockchain (roughly following [this](https://jeiwan.cc/posts/building-blockchain-in-go-part-2/) series of awesom blog posts.)
2. A set of kubernetes CRDs to manage and interact with you're blockchain. (ala `kubectl get blocks`)
3. A custome controller for controling / updating / computing PoW for each new block presented to the blockchain.

A word of caution: kubechain is currently a work in progress.

## Usage Example:
Simply create a Block CRD in you're k8s cluster:
```
> cat examples/block.yml
apiVersion: kubechain.com/v1alpha1
kind: Block
metadata:
  name: "example-block"
spec:
  data: "Move one bitcoin from Alice to Bob."

> kubectl create -f examples/block.yml
```



