---
sidebar_position: 1
---

# Verifying kacti binaries

`kacti` creates intoto provenance with every release, aligning with the Supply-chain levels for Software Artifacts (SLSA) level 3 requirements. This allows you to verify that a particular binary was created from the project and that it hasn't been tampered with post-release.

To start you'll need to install the `slsa-verifier` from the SLSA project. You can either use `go` or collect a release from the [releases page](https://github.com/slsa-framework/slsa-verifier/releases)
```bash
$ go install github.com/slsa-framework/slsa-verifier/v2/cli/slsa-verifier@v2.4.1
```
Collect the intoto provenance from the `kacti` releases page
```bash
$ curl -LO https://github.com/shaneboulden/kacti/releases/download/$(kacti version)/kacti-linux-amd64.intoto.jsonl
```
You can then use the `slsa-verifier` to verify the binary
```bash
$ slsa-verifier verify-artifact --provenance-path kacti-linux-amd64.intoto.jsonl --source-tag $(kacti version) --source-uri github.com/shaneboulden/kacti $(which kacti)
```
If everything checks out, you'll see that the verification was successful:
```bash
Verified signature against tlog entry index 67558870 at URL: https://rekor.sigstore.dev/api/v1/log/entries/24296fb24b8ad77a7489ad0f9b01a8af8a046835e0e1d7015c6e8985be217c32cda158add851335e
Verified build using builder "https://github.com/slsa-framework/slsa-github-generator/.github/workflows/builder_go_slsa3.yml@refs/tags/v1.9.0" at commit 8bac905a81d747c47dc4ec9dfa6f1e28f9e432df
Verifying artifact /usr/local/bin/kacti: PASSED

PASSED: Verified SLSA provenance
```