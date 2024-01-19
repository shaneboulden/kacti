---
sidebar_position: 1
---

# Verifying kacti binaries

`kacti` creates intoto provenance with every release, aligning with the Supply-chain levels for Software Artifacts (SLSA) level 3 requirements. This allows you to verify that a particular binary was created from the project and that it hasn't been tampered with post-release.

To start you'll need to install the `slsa-verifier` from the SLSA project. You can either use `go` or collect a release from the [releases page](https://github.com/slsa-framework/slsa-verifier/releases)
```
$ go install github.com/slsa-framework/slsa-verifier/v2/cli/slsa-verifier@v2.4.1
```
Collect the intoto provenance from the `kacti` releases page
```
$ curl -LO https://github.com/shaneboulden/kacti/releases/download/$(kacti version)/kacti-linux-amd64.intoto.jsonl
```
You can then use the `slsa-verifier` to verify the binary
```
$ slsa-verifier verify-artifact --provenance-path kacti-linux-amd64.intoto.jsonl --source-tag $(kacti version) --source-uri github.com/shaneboulden/kacti /path/to/kacti
```
If everything checks out, you'll see that the verification was successful:
```
Verified signature against tlog entry index 64436711 at URL: https://rekor.sigstore.dev/api/v1/log/entries/24296fb24b8ad77aa170f0953fbd003e3c832b9c0a940636aa7f8dda165af6af1ca4b5e87503cc5d
Verified build using builder "https://github.com/slsa-framework/slsa-github-generator/.github/workflows/builder_go_slsa3.yml@refs/tags/v1.9.0" at commit ca9d59798b31d7eb6ac126a5737e0dab7ea1302c
Verifying artifact /path/to/kacti: PASSED

PASSED: Verified SLSA provenance
```