# GitShell

Minimal go library to execute git commands from the shell.

Obviously requires to have `git` on your `PATH`.

This library is implemented with 2 main side effects:
- Many logs printed (no disable option)
- Exit on errors rather than returning them to the caller

Since it's only used by kaeter and kaeter-ci we will refactor it over time.
