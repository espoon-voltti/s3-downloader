# s3-downloader

Helper utility built in [Go](https://golang.org/) to download S3 files without installing the whole AWS CLI.

Uses [`dep`](https://github.com/golang/dep)
for dependency configuration.

## Requirements

- Golang >=1.12.x
  - and `$GOPATH` set (see [#development]())
- [`dep`](https://github.com/golang/dep)
- `make`

## Installation

Download latest release from [GitHub releases](https://github.com/espoon-voltti/s3-downloader/releases)

or compile from source with:

```sh
dep ensure
make build-linux
```

## Usage

Usage: `./bin/s3downloader-linux-amd64 bucket prefix targetDir`

## Development

**NOTE:** Golang has the concept of [`$GOPATH`](https://github.com/golang/go/wiki/GOPATH#directory-layout).

This means you should clone this repository to your `$GOPATH`
(e.g. `~/go/src/github.com/espoon-voltti/s3-downloader`) to build it.

### Install Golang

with [homebrew](http://mxcl.github.io/homebrew/):

```sh
sudo brew install go
```

with [apt](http://packages.qa.debian.org/a/apt.html)-get:

```sh
sudo apt-get install golang
```

[install Golang manually](https://golang.org/doc/install)
or
[compile it yourself](https://golang.org/doc/install/source)

### Git pre-commit hook

Configure with:

```sh
ln -s ../../scripts/pre-commit .git/hooks/pre-commit
```

### Install dependencies

```sh
dep ensure
```

### Build

```sh
make build-linux
```

## Releasing

Workflow based on: <https://circleci.com/blog/publishing-to-github-releases-via-circleci/>
and: <https://circleci.com/docs/2.0/workflows/#executing-workflows-for-a-git-tag>

1. Ensure all changes have been reviewed & merged to `master`
1. Create a git tag in the `master` branch:

    ```sh
    git tag -a v1.0.1 -m "- This is a change included in this release"
    ```

    - Name releases like: `vX.X.X`, following [semantic versioning](https://semver.org/)
    - Include all changes in message
1. Push git tag:

    ```sh
    git push --follow-tags
    ```

    - **NOTE:** The push **must only contain the tag**,
      and not commits as they will be rejected (all changes must go through PRs)
1. CircleCI creates a GitHub release with artifacts from the tag
    - **NOTE:**: Requires a GitHub Personal access token configured as `GITHUB_TOKEN` for the CircleCI build
        - Exists for @ci-voltti (with scope `repo`)
