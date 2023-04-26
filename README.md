## kacti
`kacti` (Kubernetes admission control testing) is a command-line tool for verification of Kubernetes admission controllers.

`kacti` is designed to functionally test whether admission control is correctly configured. It attempts to deploy known-bad containers to Kubernetes clusters, and verifies whether the containers successfully deploy.

## kacti tests
`kacti` uses a standard format for describing tests, shown below:
```
kacti-tests:
- name: pwnkit
  image: quay.io/the-worst-containers/pwnkit:v0.2
  namespace: app-deploy
```

Kacti returns results and information about why the test passed / failed.
