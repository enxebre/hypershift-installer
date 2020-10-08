# CAPI / Hypershift poc

Export your management cluster kubeconfig
```
export KUBECONFIG=path-to-management-kubeconfig
```

Apply the CAPI CRDs
```
kubectl apply -f ./pkg/capi/crds/
```

Clone CAPI repo, build it and run the core controllers:
```
git clone git@github.com:kubernetes-sigs/cluster-api.git
make manager-core
./bin/manager
```

Build the CAPI custom controllers binary and run them:
```
make build-capi
./bin/capiManager --alsologtostderr -v 4
```

Apply the CRs:
```
kubectl create ns cluster-farm
kubectl apply -f ./pkg/capi/resources/cluster.yaml
```

Get the managed cluster kubeconfig
```
kubectl get secret -nagl-control-plane admin-kubeconfig --template={{.data.kubeconfig}} | base64 -D
```

To delete the hosted control plane
```
kubectl delete ns test-hosted-control-plane
```

To generate a new hostedControlPlane resource:
```bash
./bin/hypershift-installer create install-config hcp
See hostedControlPlane.yaml
```

