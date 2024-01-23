---
sidebar_position: 1
---
![kacti-logo](/img/kacti-logo.png)

## What is kacti?

`kacti` is a command-line tool for verification of Kubernetes admission controllers.

`kacti` is designed to functionally test whether admission control is correctly configured. It attempts to deploy known-bad containers to Kubernetes clusters, and verifies whether the containers successfully deploy.

`kacti` uses a simple, human-readable format for admission control validatation tests ([trials](/docs/kacti-trials/kacti-trials)), shown below:
```yaml
---
- name: log4shell
  description: |
    Tests whether container images vulnerable to Log4Shell (CVE-2021-44228)
    are accepted by the cluster
  image: quay.io/smileyfritz/log4shell-app:v0.5
  namespace: kacti
```
