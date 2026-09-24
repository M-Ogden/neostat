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
	Syntax: neostat [-c|-h|-l|-s] [-d <dir>]
	options:
	h        Print this Help.
	c        Show the version of currently running Neo4j.
	d <dir>  Base directory of Neo4j installs.
	l        List installed versions of Neo4j.
	s        Save the base directory (from -d) to ~/.neostat.
	`

// configFile is the name of the file in the user's home directory that stores
// the saved base directory of Neo4j installs.
const configFile = ".neostat"

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
	var save, list, running bool

	// Parse all arguments first so that flag order does not matter (e.g.
	// "-l -d <dir>" behaves the same as "-d <dir> -l").
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
		case "-s":
			save = true
		case "-l":
			list = true
		case "-c":
			running = true
		default:
			fmt.Fprintf(os.Stderr, "unknown option: %s\n\n", args[i])
			fmt.Print(help)
			os.Exit(1)
		}
	}

	if save {
		saveBaseDir(baseDir)
	}

	if running {
		printRunningVersion()
		return
	}

	// Default action (list) applies for -l, or when only -d/-s were supplied.
	_ = list
	listInstalled(baseDir)
}

// configPath returns the full path to the ~/.neostat config file.
func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configFile), nil
}

// saveBaseDir writes the given base directory to ~/.neostat. It requires a
// non-empty directory (supplied via -d).
func saveBaseDir(dir string) {
	if strings.TrimSpace(dir) == "" {
		fmt.Fprintln(os.Stderr, "-s requires a base directory; supply one with -d <dir>")
		os.Exit(1)
	}
	path, err := configPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot determine home directory: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(path, []byte(dir+"\n"), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "cannot write %s: %v\n", path, err)
		os.Exit(1)
	}
	fmt.Printf("Saved base directory to %s\n", path)
}

// readSavedBaseDir returns the base directory stored in ~/.neostat, or an
// empty string if the file does not exist or cannot be read.
func readSavedBaseDir() string {
	path, err := configPath()
	if err != nil {
		return ""
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
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
		line := scanner.Text()
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
// If baseDir is empty, it falls back to the directory saved in ~/.neostat, and
// then to the default $HOME/neo4j/installs/instance1.
func listInstalled(baseDir string) {
	// Resolution order: explicit -d value > saved ~/.neostat value > default.
	dir := baseDir
	if dir == "" {
		dir = readSavedBaseDir()
	}
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot determine home directory: %v\n", err)
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
