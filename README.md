# zfdiff
![GitHub release (latest by date)](https://img.shields.io/github/v/release/lupinthe14th/zfdiff)
![](https://github.com/lupinthe14th/zfdiff/workflows/CI/badge.svg)
![](https://github.com/lupinthe14th/zfdiff/workflows/release/badge.svg)
[![codecov](https://codecov.io/gh/lupinthe14th/zfdiff/branch/main/graph/badge.svg?token=QXYBJF3NBE)](undefined)
[![Go Report Card](https://goreportcard.com/badge/github.com/lupinthe14th/zfdiff)][goreportcard]

## Introduction

Can see the difference between the information exported from the DNS records configured by Route53 and the locally managed Zone files.


## Installation

Download the binary from [GitHub Releases][release] and drop it in your `$PATH`.

- [Darwin / Mac][release]
- [Linux][release]
- [Windows][release]


## Getting Started

### Preparation

Export Route53 as a BIND zone file:

```
$ cli53 export example.com > r53.txt
```

SeeAlso: [cli53](https://github.com/barnybug/cli53)


Excluding comments and blank lines in local zone files:

```
$ grep -v '^;' example.com.zonefile.txt | grep -v '^$' > local.txt
```

### Differential comparison of zone files

```
zfdiff r53.txt local.txt
```


<!-- links -->
[goreportcard]: https://goreportcard.com/report/github.com/lupinthe14th/zfdiff
[release]: https://github.com/lupinthe14th/zfdiff/releases/latest
