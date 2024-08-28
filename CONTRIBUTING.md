# How to contribute to k6lint

Thank you for your interest in contributing to **k6lint**!

Before you begin, make sure to familiarize yourself with the [Code of Conduct](CODE_OF_CONDUCT.md). If you've previously contributed to other open source project, you may recognize it as the classic [Contributor Covenant](https://contributor-covenant.org/).


## Prerequisites

The tools listed in the [tools] section should be installed before contributing. It is advisable to first install the [cdo] tool, which can be used to easily perform the tasks described here. The [cdo] tool can most conveniently be installed using the [eget] tool.

```bash
eget szkiba/cdo
```

The [cdo] tool can then be used to perform the tasks described in the following sections.

Help about tasks:

```
cdo
```

[cdo]: (https://github.com/szkiba/cdo)
[eget]: https://github.com/zyedidia/eget

## tools - Install the required tools

Contributing will require the use of some tools, which can be installed most easily with a well-configured [eget] tool.

```bash
eget mikefarah/yq
eget atombender/go-jsonschema
eget szkiba/mdcode
eget golangci/golangci-lint
eget goreleaser/goreleaser
```

[tools]: (#tools---install-the-required-tools)

## schema - Contribute to the JSON schema

The JSON schema of the compliance test can be found in the [compliance.schema.yaml] file, after modification of which the schema in JSON format ([compliance.schema.json]) and the golang data model ([compliance_gen.go]) must be regenerated.

```bash
yq -o=json -P docs/compliance.schema.yaml > docs/compliance.schema.json
go-jsonschema --capitalization ID -p k6lint --only-models -o compliance_gen.go docs/compliance.schema.yaml
```

[compliance.schema.json]: docs/compliance.schema.json
[compliance_gen.go]: compliance_gen.go

## readme - Update README.md

After changing the CLI tool or examples directory, the documentation must be updated in README.md.

```bash
go run ./tools/gendoc README.md
mdcode update README.md
```

### lint - Run the linter

The `golangci-lint` tool is used for static analysis of the source code.
It is advisable to run it before committing the changes.

```bash
golangci-lint run
```

### test - Run the tests

```bash
go test -count 1 -race -coverprofile=coverage.txt ./...
```

[test]: <#test---run-the-tests>

### coverage - View the test coverage report

Requires
: [test]

```bash
go tool cover -html=coverage.txt
```

### build - Build the executable binary

This is the easiest way to create an executable binary (although the release process uses the `goreleaser` tool to create release versions).

```bash
go build -ldflags="-w -s" -o k6lint ./cmd/k6lint
```

[build]: <#build---build-the-executable-binary>

### snapshot - Creating an executable binary with a snapshot version

The goreleaser command-line tool is used during the release process. During development, it is advisable to create binaries with the same tool from time to time.

```bash
goreleaser build --snapshot --clean --single-target -o k6lint
```

[snapshot]: <#snapshot---creating-an-executable-binary-with-a-snapshot-version>

### docker - Build docker image

Building a Docker image. Before building the image, it is advisable to perform a snapshot build using goreleaser. To build the image, it is advisable to use the same `Docker.goreleaser` file that `goreleaser` uses during release.

Requires
: snapshot

```bash
docker build -t k6lint -f Dockerfile.goreleaser .
```

### clean - Delete the build directory

```bash
rm -rf build
```

## all - Update everything

The most robust thing is to update everything after modifying the source.

Requires
: schema, readme
