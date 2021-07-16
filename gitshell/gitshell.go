package gitshell

import (
	"fmt"
	"os"
	"os/exec"
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
