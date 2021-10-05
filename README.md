<!--
SPDX-FileCopyrightText: 2018-2020 City of Espoo

SPDX-License-Identifier: MIT
-->

# s3-downloader

Helper utility built in [Go](https://golang.org/) to download S3 files without installing the whole AWS CLI. Uses AWS SDK for Go, so supports all its configuration methods
(`~/.aws/credentials`, environment variables, credential servers etc).

Uses [Go modules](https://blog.golang.org/using-go-modules) for dependency configuration.

## Requirements

- Golang >=1.17.x with Go modules enabled (see [development](#development))
- *OPTIONAL*: `make`

## Installation

Download latest release from [GitHub releases](https://github.com/espoon-voltti/s3-downloader/releases)

or [compile from source](#build)

## Usage

Usage: `./bin/s3downloader-linux-amd64 <bucket> <prefix> <targetDir>`

### Logging

This tool logs in a simple JSON format (Voltti's "app-misc" type logs):

```json
{
    "@timestamp": "2020-12-29T12:02:51+02:00",
    "appBuild": "123",
    "appCommit": "22eb5a510ae608615a1487ea2538673c7f84a765",
    "appName": "my-app",
    "env": "dev",
    "exception": "Error", // Optional
    "hostIp": "10.0.0.2",
    "logLevel": "error",
    "message": "Unable to list objects in bucket: my-bucket",
    "stackTrace": "MissingRegion: could not find region configuration", // Optional
    "type": "app-misc",
    "userIdHash": "",
    "version": "1"
}
```

Configure the output with the following environment variables:

| Field | Key | Default |
|-|-|-|
| `appBuild` | `APP_BUILD` | `"local"` |
| `appCommit` | `APP_COMMIT` | `"HEAD"` |
| `appName` | `APP_NAME` | `""` |
| `env` | `VOLTTI_ENV` or `ENV` | `"local"` |
| `hostIp` | `HOST_IP` | `""` |

## Development

**NOTE:** This project uses Go modules, so you should check this repository outside `$GOPATH/src` to build it.

### Install Golang

With [homebrew](http://mxcl.github.io/homebrew/):

```sh
sudo brew install go
```

With [apt](http://packages.qa.debian.org/a/apt.html)-get:

```sh
sudo apt-get install golang
```

Or [install Golang manually](https://golang.org/doc/install)
or
[compile it yourself](https://golang.org/doc/install/source)

### Git pre-commit hook

Configure with:

```sh
ln -s ../../scripts/pre-commit.sh .git/hooks/pre-commit
```

### Build

```sh
make build
# or
make
```

## Releasing

**NOTE:**: Requires a GitHub Personal access token configured as `GITHUB_TOKEN` for the CircleCI build (exists for @ci-voltti with scope `repo`)

Workflow based on: <https://circleci.com/blog/publishing-to-github-releases-via-circleci/>
and: <https://circleci.com/docs/2.0/workflows/#executing-workflows-for-a-git-tag>

1. Ensure all changes have been reviewed & merged to `master` and you have no local changes
1. `make release`
1. CI creates a GitHub release with artifacts from the tag
