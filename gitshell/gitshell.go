package gitshell

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// GitResolveRevision prints the SHA1 hash given a revision specifier
// see https://git-scm.com/docs/git-rev-parse for more details
func GitResolveRevision(inPath, branch string) string {
	var (
		cmdOut []byte
		err    error
	)
	// --verify gives us a more compact error output
	if cmdOut, err = exec.Command("git", "-C", inPath, "rev-parse", "--verify", branch).CombinedOutput(); err != nil {
		if notFound, _ := regexp.Match("fatal: Needed a single revision", cmdOut); notFound {
			fmt.Fprintln(os.Stderr, "Could not resolve passed commit identifier: ", branch)
		} else {
			fmt.Fprintln(os.Stderr, "There was an error running the git rev-parse command: ", err)
		}
		os.Exit(1)
	}
	sha := string(cmdOut)

	return sha[0:40]
}

// GitAdd adds a change in the working directory to the staging area
// see https://git-scm.com/docs/git-add for more details
func GitAdd(inPath, filePath string) {
	if output, err := exec.Command("git", "-C", inPath, "add", filePath).CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("gitshell failed to add %s to %s", filePath, inPath), string(output), err)
		os.Exit(1)
	}
	fmt.Println("Successfully added File")
}

// GitCommit saves your changes to the local repository
// see https://git-scm.com/docs/git-commit for more details
func GitCommit(inPath, commitMsg string) {
	if output, err := exec.Command("git", "-C", inPath, "commit", "-m", commitMsg).CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("gitshell failed to commit (%s) to %s", commitMsg, inPath), string(output), err)
		os.Exit(1)
	}
	fmt.Println("Successfully committed with given commit message")
}

// GitCommitMessageFromHash returns the commit message from the given commit hash
// see https://git-scm.com/docs/git-log for more details
func GitCommitMessageFromHash(inPath, hash string) (string, error) {
	output, err := exec.Command("git", "-C", inPath, "log", "-n", "1", "--pretty=format:%B", hash).CombinedOutput()
	return string(output), err
}

// GitCheckout lets you navigate between the branches
// returns the output from git and error object if the command failed.
// see https://git-scm.com/docs/git-checkout for more details
func GitCheckout(inPath, commitOrBranch string) (string, error) {
	output, err := exec.Command("git", "-C", inPath, "checkout", commitOrBranch).CombinedOutput()
	return string(output), err
}

// GitBranchContains returns the branches containing a given commit hash
// This will check all branches (--all) including remotes, using the pattern to filter (via --list)
func GitBranchContains(inPath, hash, pattern string) (string, error) {
	output, err := exec.Command("git", "-C", inPath, "branch", "--all", "--contains", hash, "--list", pattern).CombinedOutput()
	return string(output), err
}

// GitFetch does what you think it does
func GitFetch(inPath string) (string, error) {
	output, err := exec.Command("git", "-C", inPath, "fetch").CombinedOutput()
	return string(output), err
}

// GitResolveRoot finds the root of a git repo given a path
// see https://git-scm.com/docs/git-rev-parse#Documentation/git-rev-parse.txt---show-toplevel for more details
func GitResolveRoot(inPath string) string {
	var (
		cmdOut []byte
		err    error
	)
	if cmdOut, err = exec.Command("git", "-C", inPath, "rev-parse", "--show-toplevel").Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error reading the repository root: ", err)
		os.Exit(1)
	}
	message := string(cmdOut)

	return strings.TrimSuffix(message, "\n")
}

// GitReset hard reset to a given commit
func GitReset(inPath, commit string) (string, error) {
	output, err := exec.Command("git", "-C", inPath, "reset", "--hard", commit).CombinedOutput()
	return string(output), err
}

// GitChange is an enumeration of possible actions perform on files within a commit.
type GitChange int

const (
	// Modified signals that the file was modified
	Modified GitChange = iota
	// Added signals that the file was modified
	Added
	// Deleted signals that the file was modified
	Deleted
)

func fromString(modifier string) (GitChange, error) {
	switch modifier {
	case "M":
		return Modified, nil
	case "A":
		return Added, nil
	case "D":
		return Deleted, nil
	default:
		return Modified, fmt.Errorf("could not parse Git modifier")
	}
}

// GitFileDiff extracts the map of files and the action that was performed on them: added, modified or delete.
func GitFileDiff(inPath, previousCommit, currentCommit string) map[string]GitChange {
	m := make(map[string]GitChange)
	var (
		cmdOut []byte
		err    error
	)
	if cmdOut, err = exec.Command("git", "-C", inPath, "diff", "--no-renames", "--name-status", previousCommit, currentCommit).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error reading the changed files: ", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(strings.NewReader(string(cmdOut))) // f is the *os.File
	r := regexp.MustCompile(`\s+`)
	for scanner.Scan() {
		s := r.Split(scanner.Text(), 2)
		if len(s) == 2 {
			if mod, err := fromString(s[0]); err == nil {
				m[s[1]] = mod
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error reading the changed files: ", err)
		os.Exit(1)
	}

	return m
}
