---
sidebar_position: 2
---

# Running trials from a file

You can also run trials from a file. `kacti` uses a simple, human-readable format to describe admission control trials. You can see an example below:
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
- **description** This is a more detailed explanation of what is being tested, e.g. specific CVEs or misconfiguration
- **image** This is a vulnerable container image that *should be blocked* if the test is successful. This corresponds directly with the test, i.e. if this test is for CVE-2021-44228, then the image will be vulnerable to this CVE.
- **namespace** The namespace to create this deployment in. StackRox allows you to apply admission control policies to different namespaces, and this allows you to run different tests in different namespaces.
- **template** (optional) A templated Kubernetes deployment. You can use templates to specify requests and limits, set the security context on a deployment, or test whether deployments are accepted that request privilege escalation.

You can run trials from a file using the following command:
```
$ kacti trials --file /path/to/file.yaml
```
`kacti` will print where the trials are being loaded from, and provide a result for each trial extracted from the file. You can see a complete example here:
```
$ cat << EOF > kacti.yaml
---
- name: log4shell
  description: |
    Tests whether container images vulnerable to Log4Shell (CVE-2021-44228)
    are accepted by the cluster
  image: quay.io/smileyfritz/log4shell-app:v0.5
  namespace: kacti
EOF

$ kacti trials --file kacti.yaml
Setting up kubeconfig from: /home/user/.kube/config
Using trials from: kacti.yaml
Running trial: log4shell { ns: kacti / img: quay.io/smileyfritz/log4shell-app:v0.5 }
Results:
log4shell { ns: kacti / img:quay.io/smileyfritz/log4shell-app:v0.5 }
 -> Success, Deployment scaled to zero replicas
```
