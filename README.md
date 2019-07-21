# kube-vault
Manage your secrets in kubernetes using simple UI

```sh
kubectl create -f ./rbac-config.yaml
kubectl create -f ./deployment.yaml
```

`kubectl port-forward deployment/kube-vault 8080:8080`