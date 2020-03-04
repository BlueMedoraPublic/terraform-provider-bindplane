BindPlane Terraform Provider
==================

* [bindplane.bluemedora.com](https://bindplane.bluemedora.com)
* [Bindplane API Documentation](https://docs.bindplane.bluemedora.com/reference#introduction)
* [Features](#features)
* [Installation](#installation)
* [Usage](#usage)
* [Building The Provider](#building-the-provider)
* [terraform.io](https://www.terraform.io)

[![Build Status](https://travis-ci.com/BlueMedoraPublic/terraform-provider-bindplane.svg?branch=master)](https://travis-ci.com/BlueMedoraPublic/terraform-provider-bindplane)
[![Go Report Card](https://goreportcard.com/badge/github.com/BlueMedoraPublic/terraform-provider-bindplane)](https://goreportcard.com/report/github.com/BlueMedoraPublic/terraform-provider-bindplane)

Features
------------
This provider can manage the following resources:

- Metric Sources
- Metric Credentials
- Metric Collectors
- Log Sources
- Log Destinations
- Log Templates
- Log Agents

See `RESOURCES.md` for more detailed descriptions


Installation
------------

1) download the latest release for your platform
2) unzip the plugin
3) copy plugin to `~/.terraform.d/plugins` For Mac / Linux and `%APPDATA%\terraform.d\plugins` for Windows


Usage
------------

See `USAGE.md` and `examples/` for detailed examples


Building
---------------------

Install the following:
- [Make](https://www.gnu.org/software/make/)
- [Docker](https://www.docker.com/) 18.x (primary build method)
- [Go](https://golang.org/doc/install) 1.14+ (alternative build method)

Clone repository anywhere on your system (outside of your GOPATH),
this repository uses go modules, and does not need to be in the GOPATH

Build with Docker:
```sh
make test
make
```

Build artifacts can be found in the `artifacts/` directory

Build without Docker:
```sh
make quick
```

Building with Docker is ideal for production use, as your binary
will be built the same way our releases are built.

<small><i><a href='http://ecotrust-canada.github.io/markdown-toc/'>Table of contents generated with markdown-toc</a></i></small>
