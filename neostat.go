// neostat is a small CLI that reports the running Neo4j Enterprise version
// and lists installed versions. It is a Go port of neos.sh.
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const help = `Print the currently running version of Neo4j.
	Syntax: neostat [-c|-h|-l] [-d <dir>]
	options:
	h        Print this Help.
	c        Show the version of currently running Neo4j.
	d <dir>  Base directory of Neo4j installs.
	l        List installed versions of Neo4j.
	`

// versionRe pulls the version out of a "neo4j-enterprise-<version>" token,
// whether it appears in a process command line or an install directory name.
var versionRe = regexp.MustCompile(`neo4j-enterprise-([^/\s]+)`)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		printRunningVersion()
		return
	}

	// baseDir may be overridden with -d <dir>; empty means use the default.
	var baseDir string

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-h", "-help", "--help":
			fmt.Print(help)
			return
		case "-d":
			if i+1 >= len(args) {
				fmt.Fprintln(os.Stderr, "-d requires a directory argument")
				os.Exit(1)
			}
			baseDir = args[i+1]
			i++
		case "-l":
			listInstalled(baseDir)
			return
		case "-c":
			printRunningVersion()
			return
		default:
			fmt.Fprintf(os.Stderr, "unknown option: %s\n\n", args[i])
			fmt.Print(help)
			os.Exit(1)
		}
	}

	// Only -d was supplied (no action); list installed versions using it.
	listInstalled(baseDir)
}

// printRunningVersion scans the process table for a running
// neo4j-enterprise process and prints its version.
func printRunningVersion() {
	out, err := exec.Command("ps", "ux").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read process list: %v\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner. Text()
		if !strings.Contains(line, "neo4j-enterprise") {
			continue
		}
		if m := versionRe.FindStringSubmatch(line); m != nil {
			fmt.Printf("%s is currently running.\n", m[1])
			return
		}
	}
	fmt.Println("No instances running.")
}

// listInstalled prints the versions installed under the given base directory.
// If baseDir is empty, it defaults to $HOME/neo4j/installs/instance1.
func listInstalled(baseDir string) {
	dir := baseDir
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os. Stderr, "cannot determine home directory: %v\n", err)
			os.Exit(1)
		}
		dir = filepath.Join(home, "neo4j", "installs", "instance1")
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read %s: %v\n", dir, err)
		os.Exit(1)
	}

	var versions []string
	for _, e := range entries {
		if m := versionRe.FindStringSubmatch(e.Name()); m != nil {
			versions = append(versions, m[1])
		}
	}
	sort.Strings(versions)
	for _, v := range versions {
		fmt.Println(v)
	}
}

