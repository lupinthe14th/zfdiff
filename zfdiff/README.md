# zfdiff

## Introduction

Can see the difference between the information exported from the DNS records configured by Route53 and the locally managed Zone files.


## Installation

```
go get -u github.com/lupinthe14th/dns/zfdiff
```

## Getting Started

### Preparation

Export Route53 as a BIND zone file:

```
$ cli53 export example.com > r53.txt
```

SeeAlso: [cli53](https://github.com/barnybug/cli53)


Excluding comments and blank lines in local zone files:

```
$ grep -v '^;' example.com.zonefile.txt | grep -v '^$' > gglocal.txt
```

### Differential comparison of zone files

```
zfdiff r53.txt local.txt
```
