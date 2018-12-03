# kubechain

kubechain is a small blockchain implementation on top of kubernetes - it includes:
1. A small implementation of a blockchain - roughly following [this](https://jeiwan.cc/posts/building-blockchain-in-go-part-2/) series of awesom blog posts.
2. A set of kubernetes CRDs to manage and interact with you're blockchain. (ala `kubectl get blocks`) - roughly following [Martin Helmlich's awsome blog post](https://www.martin-helmich.de/en/blog/kubernetes-crd-client.html) as well as [kubernetes/sample-controller](https://github.com/kubernetes/sample-controller)
3. A custome controller for controling / updating / computing PoW for each new block presented to the blockchain.

A word of caution: kubechain is currently a work in progress.

## Deploying Kubchain:
There are two ways to run and buil kubechain:
1. Build kubechain from source by running `make`.
2. Deploy a kubechain container on a kubernetes cluster using the `deployment.yml` file. This file contains a k8s `Deployment` for `kubechain` which uses the in-cluster configuration to create and monitor the crds.
3. Build and run in a container using `make image` followed by `docker run nimrodshn/kubechain`.

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



