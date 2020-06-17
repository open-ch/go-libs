# Go Libs

This is the home for various Open Systems public Go libraries. Its goal is mostly to be able to share
some code across some open sourced tools like [kaeter](https://github.com/open-ch/kaeter) or [checkdoc](https://github.com/open-ch/checkdoc)

To keep things simple, we want to keep smaller libraries in a single repository, but nothing speaks
against hosting modules in a separate repository whenever it is required.

# TOC

 - [fsutils](./fsutils) - Quick file system exploration
 - [mdutils](./mdutils) - Quick markdown parsing and exploration
 - [untar](./untar) - expand .tar.gz files to a directory
 - [logger](./logger) - simple logger interface for injecting any logger

## Note For GitHub Readers

While the content of this module is managed in an internal repository,
you may still submit PR's.

## Licensing

Unless specified otherwise in the source file (as is the case for the single file in [untar](./untar/untar.go)),
the content of this repository is licensed under the Apache 2.0 license (see LICENSE file).
