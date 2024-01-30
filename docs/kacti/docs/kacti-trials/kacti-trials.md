---
sidebar_position: 1
---

# Kacti trials

When a ship is newly constructed or comes out of a significant refit period it will go through "sea trials", or a "shakedown". This is a series of trials to test the vessel's seaworthiness - testing its speed, maneuverability, safety equipment, etc. It's conducted prior to commissioning and acceptance.

In a similar way, `kacti` uses trials to validate Kubernetes admission control. How does the admission controller perform - does it block workloads containing critical CVEs, or trying to expose SSH? Does it permit valid workloads to be accepted by the cluster?

Each trial represents a distinct test, validating whether the container image / configuration is blocked, or accepted by the Kubernetes cluster.

Trials consist of a Kubernetes API under test (currently only Deployments are supported), a name and description, a namespace, and an image.

## Running trials
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
 -> Success, Deployment scaled to zero replicas
```
## Using short-hand
`kacti` also accepts short-hand flags for imperative trials. You can see all of the options using `kacti trials -h`.
```
Perform functional verification trials against Kubernetes admission controllers.

Usage:
  kacti trials [flags]

Flags:
  -d, --deploy             Run a deployment trial
  -f, --file               Run a set of trials from a file
  -h, --help               help for trials
  -i, --image string       Image for the trial
  -n, --namespace string   Namespace for the trial
  -v, --verbose            Verbose output
```
You can see an example here:
```
$ kacti trials -d -n kacti -i quay.io/smileyfritz/log4shell-app:v0.5 log4shell
```

## More trial details
You can get more detail on trial activity using the `--verbose` flag:
```
$ kacti trials --deploy --namespace kacti --image quay.io/smileyfritz/log4shell-app:v0.5 log4shell --verbose
Setting up kubeconfig from: /home/user/.kube/config
Running trial: log4shell { ns: kacti / img: quay.io/smileyfritz/log4shell-app:v0.5 }
 -> Success, Deployment creation was blocked
```