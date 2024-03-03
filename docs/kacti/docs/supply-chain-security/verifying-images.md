---
sidebar_position: 2
---

# Verifying kacti images
Kacti images are signed using Sigstore, and provenance is recorded in the public-good Rekor instance.

You can verify the provenance of Kacti images using `cosign`:
```
$ kacti images --list
CVE: CVE-2021-44228 -> Image: quay.io/kacti/log4shell@sha256:f72efa1cb3533220212bc49716a4693b448106b84ca259c20422ab387972eed9

$ cosign verify \
   --certificate-oidc-issuer https://github.com/login/oauth \
   --certificate-identity shane.boulden@gmail.com \
   quay.io/kacti/log4shell@sha256:f72efa1cb3533220212bc49716a4693b448106b84ca259c20422ab387972eed9

Verification for quay.io/kacti/log4shell@sha256:f72efa1cb3533220212bc49716a4693b448106b84ca259c20422ab387972eed9 --
The following checks were performed on each of these signatures:
  - The cosign claims were validated
  - Existence of the claims in the transparency log was verified offline
  - The code-signing certificate was verified using trusted certificate authority certificates 
```