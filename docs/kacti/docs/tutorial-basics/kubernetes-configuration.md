---
sidebar_position: 2
---

# Kubernetes configuration

`kacti` expects access to a kubernetes config to perform testing.
- The `KUBECONFIG` environment variable takes precedence
- If this variable isn't set, `kacti` will look in the `.kube/config` directory within the user's home directory.

You can verify that your cluster is correctly configured for `kacti` by testing whether you can create deployments.
```
$ kubectl auth can-i create deploy
```
If you get a `yes`, you're good to go! Otherwise you'll need to create roles / role bindings allowing your user to create deployments.
