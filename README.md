# s3-downloader

Helper utility to download S3 files without installing the whole AWS CLI.

## Requirements

- Golang
- [`dep`](https://github.com/golang/dep)
- `make`

## Installation

Download latest release from [GitHub releases](https://github.com/espoon-voltti/s3-downloader/releases)

or compile from source with:

```
dep ensure
make build-linux
```

## Usage

Usage: `./bin/s3downloader-linux-amd64 bucket prefix targetDir`

## Development

Uses [`dep`](https://github.com/golang/dep)
for dependency configuration.

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

Workflow based on: https://circleci.com/blog/publishing-to-github-releases-via-circleci/
and: https://circleci.com/docs/2.0/workflows/#executing-workflows-for-a-git-tag

1. Create a git tag:

        git tag -a v1.0.1 -m "- This is a change included in this release"

    - Name releases like: `vX.X.X`, following [semantic versioning](https://semver.org/)
    - Include all changes in message
1. Push git tag:

        git push --tags

1. CircleCI creates a GitHub release with artifacts from the tag
    - **NOTE:**: Requires a GitHub Personal access token configured as `GITHUB_TOKEN` for the CircleCI build
        - Exists for @ci-voltti (with scope `repo`)
