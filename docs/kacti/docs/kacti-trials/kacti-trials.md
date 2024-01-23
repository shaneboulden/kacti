---
sidebar_position: 1
---

# Kacti trials

When a ship is newly constructed or comes out of a significant refit period it will go through "sea trials", or a "shakedown". This is a series of trials to test the vessel's seaworthiness - testing its speed, maneuverability, safety equipment, etc. It's conducted prior to commissioning and acceptance.

In a similar way, `kacti` uses trials to validate Kubernetes admission control. How does the admission controller perform - does it block workloads containing critical CVEs, or trying to expose SSH? Does it permit valid workloads to be accepted by the cluster?

Each trial represents a distinct test, validating whether the container image / configuration is blocked, or accepted by the Kubernetes cluster. 

Trials consist of a Kubernetes API under test (currently only Deployments are supported), a name and description, a namespace, and an image.

`kacti` uses a simple, human-readable format to describe admission control trials. You can see an example below:
```yaml
---
- name: log4shell
  description: |
    Verifies whether container images vulnerable to Log4Shell (CVE-2021-44228)
    are accepted by the cluster
  image: quay.io/smileyfritz/log4shell-app:v0.5
  namespace: kacti
  template: deploy-template.yaml
```

Let's look at this file a little closer. Each trial is a list element and has the following features: 
- **name** The name of the test, allows you to see which test is being run and results in the output.
- **description** This is a more detailed explanation of what is being tested, e.g. specific CVEs or 
- **image** This is a vulnerable container image that **should be blocked* by the cluster, if the test is successful. This corresponds directly with the test, i.e. if this test is for CVE-2021-44228, then the image will be vulnerable to this CVE.
- **namespace** The namespace to create this deployment in. StackRox allows you to apply admission control policies to different namespaces, and this allows you to run different tests in different namespaces.
- **template** (optional) A templated Kubernetes deployment. You can use templates to specify requests and limits, set the security context on a deployment, or test whether deployments are accepted that request privilege escalation.

