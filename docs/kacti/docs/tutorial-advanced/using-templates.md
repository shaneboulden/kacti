---
sidebar_position: 1
---

# Using templates

`kacti` supports templated deployments. You can reference templates in tests using the `template` keyword:
```
---
kacti-tests:
  - name: log4shell
    description: |
      Tests whether container images vulnerable to Log4Shell (CVE-2021-44228)
      are accepted by the cluster
    image: quay.io/smileyfritz/log4shell-app:v0.5
    namespace: kacti
    template: deploy-template.yaml
```
Templates should contain a complete Kubernetes deployment. `kacti` will replace the name, image and namespace for the test using the template.

You can use templates to validate deployment configuration, for example:
- testing whether deployments without limits are accepted to the cluster
- testing whether containers in a deployment can use a privileged security context
- testing whether deployments requesting host access are permitted by the cluster

## Example template usage
Create a template file:
```
cat << EOF > template.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    deployment.kubernetes.io/revision: "1"
  labels:
    app: test
  name: test
  namespace: kacti
spec:
  replicas: 1
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: test
      namespace: test
    spec:
      containers:
      - image: quay.io/podman/hello
        imagePullPolicy: IfNotPresent
        name: test
EOF
```
Reference the template if a `kacti-tests` YAML file:
```
cat << EOF > kacti.yaml
kacti-tests:
  - name: log4shell
    description: |
      Tests whether container images vulnerable to Log4Shell (CVE-2021-44228)
      are accepted by the cluster
    image: quay.io/smileyfritz/log4shell-app:v0.5
    namespace: kacti
    template: deploy-template.yaml
EOF
```
Run the test:
```
kacti test kacti.yaml
Setting up kubeconfig from: /home/user/blue-env/kubeconfig
Using tests from: kacti.yaml
Running test: pwnkit { ns: kacti / img: quay.io/the-worst-containers/pwnkit:v0.2 }
Error creating Deployment: admission webhook "policyeval.stackrox.io" denied the request:
The attempted operation violated 1 enforced policy, described below:

Policy: No resource requests or limits specified
- Description:
    ↳ Alert on deployments that have containers without resource requests and limits
- Rationale:
    ↳ If a container does not have resource requests or limits specified then the host
      may become over-provisioned.
- Remediation:
    ↳ Specify the requests and limits of CPU and Memory for your deployment.
- Violations:
    - CPU limit set to 2 cores for container 'test'
    - CPU request set to 0 cores for container 'test'
    - Memory limit set to 0 MB for container 'test'
    - Memory request set to 0 MB for container 'test'


In case of emergency, add the annotation {"admission.stackrox.io/break-glass": "ticket-1234"} to your deployment with an updated ticket number


Results:
pwnkit { ns: kacti / img:quay.io/the-worst-containers/pwnkit:v0.2 }
 -> Success, Deployment creation was blocked
```