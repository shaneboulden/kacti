---
sidebar_position: 3
---


# Running trials

Once you've created your trials, you can simply run them using `kacti trials`.
```
$ kacti trials kacti.yaml
```
For each trial, `kacti` will attempt to create a deployment in the specified namespace referencing the vulnerable image. If the image is successfully deployed and scaled up, `kacti` will report a **failed** test.

If the image is scaled down to zero replicas, or the deployment is blocked, `kacti` will report **success**.

You can see this in the output:
```
$ kacti trials kacti.yaml
Setting up kubeconfig from: /home/user/.kube/config
Using tests from: kacti.yaml
Running test: pwnkit { ns: kacti / img: quay.io/the-worst-containers/pwnkit:v0.2 }
Running test: log4shell { ns: kacti / img: quay.io/smileyfritz/log4shell-app:v0.5 }
Results:
pwnkit { ns: kacti / img:quay.io/the-worst-containers/pwnkit:v0.2 }
 -> Failed, Deployment was created successfully and scaled up

log4shell { ns: kacti / img:quay.io/smileyfritz/log4shell-app:v0.5 }
 -> Success, Deployment scaled to zero replicas
```
Once `kacti` trials are completed, it will clean up any deployments / pods that were created.

