# terraform-manager
Terraform version manager to allow for easy swapping between Terraform versions.

By default, terraform-manager will download the version of the Terraform for values of GOOS and GOARCH. In the future, this will be possible to override in config if needed.

# Getting Started
## Installation

Grab the newest release from the [releases page](https://github.com/dantdj/terraform-manager/releases) - make sure to download the relevant archive for your architecture and OS.

Optionally, rename the binary as desired.

## Usage

Run `terraform-manager download <terraform-version>` to get started - this will download the specified version of Terraform, and add it to the application config file.

Downloaded versions of Terraform are put in a sub-directory under the directory containing `terraform-manager` itself, and the config file is created on the same level as `terraform-manager`.

# Development

A Makefile exists to automate some tasks - currently this is only guaranteed to work on MacOS. See the comments in the Makefile itself for more detail.

## Building

* Makefile -  `make build/app`
* Go CLI - `go build -o <exe-name>`

The Makefile performs some additional steps to inject version information into the binary.