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
func GitResolveRevision(inPath string, branch string) string {
	var (
		cmdOut []byte
		err    error
	)
    if cmdOut, err = exec.Command("git","-C", inPath, "rev-parse", branch).Output(); err != nil {
        fmt.Fprintln(os.Stderr, "There was an error running the git rev-parse command: ", err)
        os.Exit(1)
    }
    sha := string(cmdOut)
    return sha[0:40]
}

// GitAdd adds a change in the working directory to the staging area
// see https://git-scm.com/docs/git-add for more details
func GitAdd(inPath string, filePath string ){
	if err := exec.Command("git", "-C", inPath, "add", filePath).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Successfully added File")
}

// GitCommit saves your changes to the local repository
// see https://git-scm.com/docs/git-commit for more details
func GitCommit(inPath string, commitMsg string){
	if err := exec.Command("git","-C", inPath, "commit","-m", commitMsg).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Successfully committed with given commit message")
}

// GitCommitMessageFromHash returns the commit message from the given commit hash
// see https://git-scm.com/docs/git-log for more details
func GitCommitMessageFromHash(inPath string, hash string) string{
	var (
		cmdOut []byte
		err    error
	)
    if cmdOut, err = exec.Command("git","-C", inPath, "log", "-n", "1", "--pretty=format:%B", hash).Output(); err != nil {
        fmt.Fprintln(os.Stderr, "There was an error returning the git commit message: ", err)
        os.Exit(1)
    }
    message := string(cmdOut)
    return message

}

// GitCheckout lets you navigate between the branches
// see https://git-scm.com/docs/git-checkout for more details
func GitCheckout(inPath string, commitOrBranch string ){
	if err := exec.Command("git", "-C", inPath, "checkout", commitOrBranch).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Checkout out branch or commit")
}

// GitResolveRoot finds the root of a git repo given a path
// see https://git-scm.com/docs/git-rev-parse#Documentation/git-rev-parse.txt---show-toplevel for more details
func GitResolveRoot(inPath string) string {
	var (
		cmdOut []byte
		err    error
	)
	if cmdOut, err =  exec.Command("git", "-C", inPath, "rev-parse", "--show-toplevel").Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error reading the repository root: ", err)
		os.Exit(1)
	}
	message := string(cmdOut)
	return strings.TrimSuffix(message, "\n")
}

// GitReset hard reset to a given commit
func GitReset(inPath string, commit string) error {
	if _, err :=  exec.Command("git", "-C", inPath, "reset", "--hard", commit).Output(); err != nil {
		fmt.Fprintf(os.Stderr, "There was an error reset the repo to %s: %v\n", commit, err)
		return err
	}
	fmt.Printf("Successfully reset to %s\n", commit)
	return nil
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
	r := regexp.MustCompile("[\\s]+")
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
