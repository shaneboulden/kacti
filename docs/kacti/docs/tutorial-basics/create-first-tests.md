---
sidebar_position: 3
---

# Creating your first tests

`kacti` uses a human-readable, simple testing format to describe tests. You can see some examples below:
```
---
kacti-tests:
  - name: log4shell
    description: |
      Tests whether container images vulnerable to Log4Shell (CVE-2021-44228)
      are accepted by the cluster
    image: quay.io/smileyfritz/log4shell-app:v0.5
    namespace: kacti
```

Let's look at this file a little closer. The file starts with `kacti-tests`. This file indicates to `kacti` that it contains a number of tests that should be run.

`kacti-tests` contains an array, and each test has the following features: 
- **name** The name of the test, allows you to see which test is being run and results in the output.
- **description** This is a more detailed explanation of what is being tested, e.g. specific CVEs or 
- **image** This is a vulnerable container image that **should be blocked* by the cluster, if the test is successful. This corresponds directly with the test, i.e. if this test is for CVE-2021-44228, then the image will be vulnerable to this CVE.
- **namespace** The namespace to create this deployment in. StackRox allows you to apply admission control policies to different namespaces, and this allows you to run different tests in different namespaces.
