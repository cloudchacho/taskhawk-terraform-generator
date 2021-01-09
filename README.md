# Taskhawk Terraform Generator

[![Build Status](https://travis-ci.org/standard-ai/taskhawk-terraform-generator.svg?branch=master)](https://travis-ci.org/standard-ai/taskhawk-terraform-generator)

[Taskhawk](https://github.com/standard-ai/taskhawk) is a replacement for celery that works on AWS SQS/SNS, while
keeping things pretty simple and straight forward. 

Taskhawk Terraform Generator is a CLI utility that makes the process of managing 
[Taskhawk Terraform modules](https://registry.terraform.io/search?q=taskhawk&verified=false) easier by abstracting 
away details about [Terraform](https://www.terraform.io/).

## Usage 

### Installation

Download the latest version of the release from [Github releases](https://github.com/standard-ai/taskhawk-terraform-generator/releases) - 
it's distributed as a zip containing a Go binary file.

### Configuration

Configuration is specified as a JSON file. Run 

```sh
./taskhawk-terraform-generator config-file-structure
```

to get the sample configuration file.

**Advanced usage**: The config *may* contain references to other terraform resources, as long as they resolve to 
an actual resource at runtime. 

### How to use

Run 

```sh
./taskhawk-terraform-generator --provider <cloud provider> apply-config <config file path>
```

to create Terraform modules. The module is named `taskhawk` by default in the current directory.

Re-run on any changes.

## Development

### Prerequisites

Install go1.11

### Getting Started

Assuming that you have go installed, set up your environment:

```sh
$ # in a location NOT in your GOPATH:
$ git checkout github.com/standard-ai/taskhawk-terraform-generator
$ cd taskhawk-terraform-generator
$ go-bindata -prefix "templates/" templates/*
$ go get -mod=readonly -v ./...
$ GO111MODULE=off go get github.com/go-bindata/go-bindata/...
$ GO111MODULE=off go get -u github.com/client9/misspell/cmd/misspell
$ GO111MODULE=off go get -u honnef.co/go/tools/cmd/staticcheck
```

### Running Tests

You can run tests in using ``make test``. By default, it will run all of the unit and functional tests, but you can 
also run specific tests directly using go test:

```sh
$ go test ./...
$ go test -run TestGenerate ./...
```

## Release Notes

[Github Releases](https://github.com/standard-ai/taskhawk-terraform-generator/releases)

## How to publish


```sh
make clean build
```

Upload to Github and attach the zip files created in above step.
