# neostat

`neostat` is a small command-line tool that reports the currently running
[Neo4j](https://neo4j.com/) Enterprise version and lists installed versions.
It is a Go port of the original `neos.sh` script.

## Features

- Detect and print the version of the Neo4j Enterprise instance currently running.
- List all installed Neo4j Enterprise versions under a base directory.
- Configurable install base directory.

## Requirements

- Go 1.27.0 or later (to build from source).
- A Unix-like environment with the `ps` command available (used to inspect
  running processes).

## Installation

Build from source with the Go toolchain:

```sh
go build -o neostat neostat.go
```

This produces a `neostat` binary in the current directory. Move it somewhere on
your `PATH` if desired:

```sh
mv neostat /usr/local/bin/
```

## Usage

```
neostat [-c|-h|-l] [-d <dir>]
```

### Options

| Option     | Description                                          |
| ---------- | ---------------------------------------------------- |
| `-h`       | Print the help message.                              |
| `-c`       | Show the version of the currently running Neo4j.     |
| `-d <dir>` | Base directory of Neo4j installs.                    |
| `-l`       | List installed versions of Neo4j.                    |

Running `neostat` with no arguments prints the currently running version
(equivalent to `-c`).

### Examples

Print the currently running version:

```sh
neostat
# or
neostat -c
```

List installed versions using the default base directory
(`$HOME/neo4j/installs/instance1`):

```sh
neostat -l
```

List installed versions from a custom base directory:

```sh
neostat -l -d /opt/neo4j/installs
```

## How it works

- **Running version:** `neos` scans the process table (`ps ux`) for a
  `neo4j-enterprise` process and extracts the version from its command line.
- **Installed versions:** `neos` reads the entries of the install base
  directory and extracts versions from directory names matching
  `neo4j-enterprise-<version>`.

## License

See the repository for license details.
