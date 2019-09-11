# kube-vault

[![](https://github.com/exelban/kube-vault/workflows/Deploy%20to%20docker%20hub/badge.svg)](https://github.com/exelban/kube-vault/actions")
[![](https://images.microbadger.com/badges/version/exelban/kube-vault.svg)](https://github.com/exelban/kube-vault")

Manage your secrets in kubernetes using simple UI.

[![Preview](https://serhiy.s3.eu-central-1.amazonaws.com/Github_repo/kube-vault/v0.0.1.png)](https://serhiy.s3.eu-central-1.amazonaws.com/Github_repo/kube-vault/v0.0.1.png)

## Features
kube-vault is a simple application which allows you to manage secrets in kubernetes.

What can you do:  

- see the list of secret
- create a secret
- delete a secret
- create new value in secret
- delete a value from secret
- modify a value in secret

## Installation
### Quick version
```sh
kubectl create -f ./rbac-config.yaml
kubectl create -f ./deployment.yaml
```

To expose the application locally:

```sh
kubectl port-forward deployment/kube-vault 8080:8080
```

After that application will be available under [http://localhost:8080](http://localhost:8080).

### Long version
Application requires [rbac-account](https://github.com/exelban/kube-vault/blob/master/kubernetes/rbac-config.yaml) for receiving secrets from namespace. To create this account just create a new account in kubernetes:

```sh
kubectl create -f ./rbac-config.yaml
```

If there will be no RBAC-account, the application could not authorize and fetch a list of secrets.

To create the deployment just use:

```sh
kubectl create -f ./deployment.yaml
```

There is no service for exposing an application. For accessing a web interface better use port-forwarding:

```sh
kubectl port-forward deployment/kube-vault 8080:8080
```

If you want to create your deployment you can use [docker](https://cloud.docker.com/u/exelban/repository/docker/exelban/kube-vault) image: `exelban/kube-vault:latest`.

## Application
The application contains two parts:

- backend (golang)
- frontend (svelte)

Backend application receive next environment variables:

| Variable | Default | Description |
| --- | --- | --- |
| NAMESPACE | default | Default namespace |
| HOME | | Absolute path to the kubeconfig file (optional)  |
| PORT | 8080 | Rest server port |


Frontend application receive environment variables:

| Variable | Default | Description |
| --- | --- | --- |
| NAMESPACE | default | Default namespace |
| API_HOST | http://localhost:8080 | Host for backend application |

## License
[MIT License](https://github.com/exelban/kube-vault/blob/master/LICENSE)
