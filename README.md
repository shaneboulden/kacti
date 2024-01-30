![kacti logo](./docs/img/kacti-logo.png)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache2.0-brightgreen.svg)](https://opensource.org/licenses/Apache-2.0) | [![SLSA 3](https://slsa.dev/images/gh-badge-level3.svg)](https://slsa.dev)

# What is kacti?
`kacti` is a command-line tool for verification of Kubernetes admission controllers.

`kacti` is designed to functionally test whether admission control is correctly configured. It attempts to deploy known-bad containers to Kubernetes clusters, and verifies whether the containers successfully deploy.

## Quick start
Grab the latest `kacti` binary:
```
$ curl -Lo kacti https://github.com/shaneboulden/kacti/releases/latest/download/kacti-linux-amd64 && \
      sudo mv kacti /usr/local/bin/kacti && \
      sudo chmod 0755 /usr/local/bin/kacti
```
Ensure that you're logged into a Kubernetes cluster and have permissions to create deployments:
```
$ export KUBECONFIG=/path/to/kubeconfig

$ kubectl auth can-i create deploy
yes
```
Run `kacti`:
```
$ kacti trials --deploy --namespace kacti --image quay.io/smileyfritz/log4shell-app:v0.5 log4shell
 -> Success, Deployment scaled to zero replicas
```
You can find more `kacti` guides in the [docs](https://kacti.dev/docs/intro).

## kacti and SLSA
`kacti` binaries are signed with Sigstore, and provenance is available and stored in the public-good Rekor instance. 

Check out the [docs](https://kacti.dev/docs/supply-chain-security/verifying-binaries) for steps to verify `kacti` binary provenance.