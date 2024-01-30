---
sidebar_position: 3
---
# Running your first trials

`kacti` uses trials to validate Kubernetes admission control. How does the admission controller perform - does it block workloads containing critical CVEs, or trying to expose SSH? Does it permit valid workloads to be accepted by the cluster?

Each trial represents a distinct test, validating whether the container image / configuration is blocked, or accepted by the Kubernetes cluster.

You can run trials using the following command:
```
$  kacti trials --deploy --namespace kacti --image quay.io/smileyfritz/log4shell-app:v0.5 log4shell
```
In this example:
- The name of the trial is `log4shell`
- `kacti` will attempt to create a Kubernetes deployment in the namespace `kacti`, named `log4shell`
- The container image used for the deployment will be `quay.io/smileyfritz/log4shell-app:v0.5`

`kacti` will display the result of the trial. If the deployment was successfully created and scaled up, the result will be a `failure`. Otherwise, if the deployment creation was blocked, or the number of replicas was scaled to zero, the result will be `success`.
```
$ kacti trials --deploy --namespace kacti --image quay.io/smileyfritz/log4shell-app:v0.5 log4shell
Setting up kubeconfig from: /home/user/.kube/config
Running trial: log4shell { ns: kacti / img: quay.io/smileyfritz/log4shell-app:v0.5 }
 -> Success, Deployment scaled to zero replicas
```

