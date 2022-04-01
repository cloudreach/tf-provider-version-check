# Terraform Provider Version Check

A utility to check whether Terraform providers configured in a project are up-to-date.

* Reads `.terraform.lock.hcl` file to get details of providers used in a project
* Looks up `registry.terraform.io` to get the latest version for each
* If updates are available this is notified on STDOUT
  * Can optionally return code on exit if update is available

## Usage

```shell
tfpvc-osx-amd64 --help
A utility to check whether Terraform providers configured in a project are up-to-date.

* Reads .terraform.lock.hcl file to get details of providers used in a project
* Looks up registry.terraform.io to get the latest version for each
* If updates are available this is notified on STDOUT
  * Can optionally return code on exit if update is available

Usage:
  tfpvc [flags]
  tfpvc [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Returns version data

Flags:
      --errorOnUpdate   Exit with error code if updates are available
  -h, --help            help for tfpvc
      --tfDir string    Directory with TF Files (default ".")

```

## Requirements
* Terraform must be installed and available on the system PATH

## Building
* Golang 1.17 is required to build
* To build binary for all platform/architectures (Linux/amd64, OSX/amd64, OSX/arm64)
```shell
make
```

## Installation
* To install binary to /usr/local/bin:

```shell
# Linux/amd64
make install-linux

# OSX/amd64
make install-mac-intel

# OSX/arm64
make install-mac-applesilicon
```