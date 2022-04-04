# Terraform Provider Version Check

A utility to check whether Terraform providers configured in a project are up-to-date.

* Reads `.terraform.lock.hcl` file to get details of providers used in a project
* Looks up `registry.terraform.io` to get the latest version for each
* If updates are available this is notified on STDOUT
  * Can optionally return code on exit if update is available

## Usage

```shell
tfpvc help
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
      --findLockFiles   Search for lockfiles in tfDir (default true)
  -h, --help            help for tfpvc
      --tfDir string    Directory with TF Files (default ".")

```

## pre-commit hook
With pre-commit, you can ensure you are notified of updates to your Terraform provider config each time you make a commit.

First install `pre-commit` and then create or update a `.pre-commit-config.yaml` in the root of your Git repo with at least the following content:

```yaml
repos:
  - repo: https://github.com/cloudreach/tf-provider-version-check
    rev: "1.0.0"
    hooks:
      - id: terraform-provder-version-check
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
