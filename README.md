# Taskhawk Terraform Generator

[![Build Status](https://travis-ci.org/Automatic/taskhawk-terraform-generator.svg?branch=master)](https://travis-ci.org/Automatic/taskhawk-terraform-generator)

[Taskhawk](https://github.com/Automatic/taskhawk) is a replacement for celery that works on AWS SQS/SNS, while
keeping things pretty simple and straight forward. 

Taskhawk Terraform Generator is a CLI utility that makes the process of managing 
[Taskhawk Terraform modules](https://registry.terraform.io/search?q=taskhawk&verified=false) easier by abstracting 
away details about [Terraform](https://www.terraform.io/).

## Usage 

### Installation

Download the latest version of the release from [Github releases](https://github.com/Automatic/taskhawk-terraform-generator/releases) - 
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
./taskhawk-terraform-generator apply-config <config file path>
```

to create Terraform modules. The module is named `taskhawk` by default in the current directory.

Re-run on any changes.

## Development

### Getting Started

Assuming that you have go installed, set up your environment:

```sh
$ go get github.com/kardianos/govendor
$ go get github.com/go-bindata/go-bindata/...
$ go get github.com/Automatic/taskhawk-terraform-generator
$ cd ${GOPATH}/src/github.com/Automatic/taskhawk-terraform-generator
$ govendor sync
```

### Running Tests

You can run tests in using ``make test``. By default, it will run all of the unit and functional tests, but you can 
also run specific tests directly using go test:

```sh
$ go test ./...
$ go test -run TestGenerate ./...
```

## Release Notes

[Github Releases](https://github.com/Automatic/taskhawk-terraform-generator/releases)

## How to publish


```sh
make clean build

cd bin/linux-amd64 && zip taskhawk-terraform-generator-linux-amd64.zip taskhawk-terraform-generator; cd -
cd bin/darwin-amd64 && zip taskhawk-terraform-generator-darwin-amd64.zip taskhawk-terraform-generator; cd -
```

Upload to Github and attach the zip files created in above step.
