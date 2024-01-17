---
sidebar_position: 1
---

# Verifying kacti binaries

`kacti` creates intoto provenance with every release, aligning with the Supply-chain levels for Software Artifacts (SLSA) level 3 requirements. This allows you to verify that a particular binary was created from the project and that it hasn't been tampered with post-release.

To start you'll need to install the `slsa-verifier` from the SLSA project. You can either use `go` or collect a release from the [releases page](https://github.com/slsa-framework/slsa-verifier/releases)
```
o install github.com/slsa-framework/slsa-verifier/v2/cli/slsa-verifier@v2.4.1
```
Check the version of a `kacti` binary on your system:
```
$ kacti version
Version: 0.1.4
Commit: ca9d59798b31d7eb6ac126a5737e0dab7ea1302c
Commit Date: 1705532141
Tree State: clean
```
Once you have the version, you can collect the intoto provenance from the `kacti` releases page, using the version number prefixed with 'v':
```
$ curl -LO https://github.com/shaneboulden/kacti/releases/download/v0.1.4/kacti-linux-amd64.intoto.jsonl
```
You can then use the `slsa-verifier` to verify the binary. The source tag will match the release version, prefixed with 'v'.
```
$ slsa-verifier verify-artifact --provenance-path kacti-linux-amd64.intoto.jsonl --source-tag v0.1.4 --source-uri github.com/shaneboulden/kacti /path/to/kacti
```
If everything checks out you'll see that the verification was successful:
```
Verified signature against tlog entry index 64436711 at URL: https://rekor.sigstore.dev/api/v1/log/entries/24296fb24b8ad77aa170f0953fbd003e3c832b9c0a940636aa7f8dda165af6af1ca4b5e87503cc5d
Verified build using builder "https://github.com/slsa-framework/slsa-github-generator/.github/workflows/builder_go_slsa3.yml@refs/tags/v1.9.0" at commit ca9d59798b31d7eb6ac126a5737e0dab7ea1302c
Verifying artifact /path/to/kacti: PASSED
```