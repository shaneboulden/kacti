---
sidebar_position: 4
---

# Specifying CVEs for trials
`kacti` can use a CVE identifier instead of an explicit image reference for a trial:
```
kacti trials --deploy --cve CVE-2021-44228 -n kacti log4shell
 -> Success, Deployment creation was blocked
```

## Supported CVEs and images
When you specify a CVE `kacti` uses a signed image to perform the trial. The following table shows the currently supported CVEs and images for `kacti`:

| CVE | Image | Source | Comments |
| --- | ----- | ------ | -------- |
| CVE-2021-44228 | quay.io/kacti/log4shell | https://github.com/shaneboulden/log4shell-vulnerable-app | Log4Shell image |

