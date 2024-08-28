<h1 name="title">k6lint</h1>

**Linter for k6 extensions**

k6lint is a command line tool and a library for static analysis of the source of k6 extensions. The analysis is done without building a k6 executable with the extension.

The contents of the source directory are used for analysis. If the directory is a git workdir, it also analyzes the git metadata. The analysis is completely local and does not use external APIs (e.g. repository manager API) or services.

The result of the analysis is compliance expressed as a percentage (0-100). This value is created as a weighted, normalized value of the scores of each checker. A compliance grade is created from the percentage value (A-F).

The detailed result of the checks are described in a [JSON schema](https://grafana.github.io/k6lint/compliance.schema.json).

### Checkers

  - `module` - checks if there is a valid go.mod
  - `replace` - checks if there is a replace directive in go.mod
  - `readme` - checks if there is a readme file
  - `examples` - checks if there are files in the examples directory
  - `license` - checks whether there is a suitable OSS license
  - `git` - checks if the directory is git workdir
  - `versions` - checks for semantic versioning git tags

## Install

Precompiled binaries can be downloaded and installed from the [Releases](https://github.com/grafana/k6lint/releases) page.

If you have a go development environment, the installation can also be done with the following command:

```
go install github.com/grafana/k6lint/cmd/k6lint@latest
```

## Use

```bash
k6lint source-directory
```

**output**

```text file=docs/example.txt
k6 extension compliance
──────────┬─────────────
 grade  A │ level  100% 

Details
───────
✔ module              
  found `github.com/grafana/xk6-sql` as go module
✔ replace             
  no `replace` directive in the `go.mod` file
✔ readme              
  found `README.md` as README file
✔ examples            
  found `examples` as examples directory
✔ license             
  found `LICENSE` as `Apache-2.0` license
✔ git                 
  found git worktree
✔ versions            
  found `6` versions, the latest is `v0.3.0`

```

<details>
<summary>JSON output</summary>

```json file=docs/example.json
{
  "checks": [
    {
      "details": "found `github.com/grafana/xk6-sql` as go module",
      "id": "module",
      "passed": true
    },
    {
      "details": "no `replace` directive in the `go.mod` file",
      "id": "replace",
      "passed": true
    },
    {
      "details": "found `README.md` as README file",
      "id": "readme",
      "passed": true
    },
    {
      "details": "found `examples` as examples directory",
      "id": "examples",
      "passed": true
    },
    {
      "details": "found `LICENSE` as `Apache-2.0` license",
      "id": "license",
      "passed": true
    },
    {
      "details": "found git worktree",
      "id": "git",
      "passed": true
    },
    {
      "details": "found `6` versions, the latest is `v0.3.0`",
      "id": "versions",
      "passed": true
    }
  ],
  "grade": "A",
  "level": 100,
  "timestamp": 1724833956
}
```
</details>

## CLI

<!-- #region cli -->
## k6lint

Linter for k6 extensions

### Synopsis

Static analyzer for k6 extensions

k6lint analyzes the source of the k6 extension without building a k6 executable with the extension.

By default, text output is generated. The `--json` flag can be used to generate the result in JSON format.

If the grade is `C` or higher, the command is successful, otherwise it returns an exit code larger than 0.
This passing grade can be modified using the `--passing` flag.


```
k6lint [flags] [directory]
```

### Flags

```
      --passing A|B|C|D|E|F   set lowest passing grade (default C)
  -q, --quiet                 no output, only validation
  -o, --out string            write output to file instead of stdout
      --json                  generate JSON output
  -c, --compact               compact instead of pretty-printed JSON output
  -V, --version               print version
  -h, --help                  help for k6lint
```

<!-- #endregion cli -->
