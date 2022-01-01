package bazel

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)


// LabelCharacters is the list of valid character for a Bazel package
var LabelCharacters = "a-zA-Z0-9-_"

// RepoLabelRegex the regex to find labels of the current repository
var RepoLabelRegex = fmt.Sprintf("//([%s]+/)*[%s]+(:[%s]+){0,1}", LabelCharacters, LabelCharacters, LabelCharacters)

// PackageLabelRegex the regex to find labels of the current package
var PackageLabelRegex = fmt.Sprintf(":[%s]+", LabelCharacters)


// Query performs a Bazel query and returns each line of the result
func Query(folder, query string, flags []string)  ([]string, error) {
	cmd := exec.Command("bazel", "query", query)
	cmd.Args = append(cmd.Args, flags...)
	cmd.Dir = folder
	var (
		cmdOut []byte
		err    error
	)
	if cmdOut, err = cmd.Output(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "BazelShell Query - Failure to execute command: ", cmd.Args)
		_, _ = fmt.Fprintln(os.Stderr, "BazelShell Query - Error:", err)
		if cmdOut != nil && cmd.ProcessState.ExitCode() == 3 {
			// There is an error, it exited with return code three, meaning that we got a partial error, which
			// is not problematic as it's most likely du to --keep_going.
			// See https://docs.bazel.build/versions/main/guide.html#what-exit-code-will-i-get
			return strings.Split(string(cmdOut), "\n"), nil
		}
		return nil, err
	}
	return strings.Split(string(cmdOut), "\n"), nil
}
