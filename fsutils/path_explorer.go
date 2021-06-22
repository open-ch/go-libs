package fsutils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// SearchByFileName Given a path, returns all sub-paths to files that are named exactly like fileName
// rootPath must be absolute
// TODO some thing that supports "real" globs would be nice.
func SearchByFileName(rootPath string, baseName string) ([]string, error) {
	if !filepath.IsAbs(rootPath) {
		return nil, fmt.Errorf("rootPath is not absolute: %s", rootPath)
	}
	if len(baseName) == 0 {
		return nil, fmt.Errorf("baseName cannot be empty")
	}

	return basenameGlob(rootPath, baseName)
}

// SearchByFileNameRegex Given a path, returns all sub-paths to files that have names that match the passed regex.
// rootPath must be absolute
// TODO some thing that supports "real" globs would be nice.
func SearchByFileNameRegex(rootPath string, baseNameRegex string) ([]string, error) {
	if !filepath.IsAbs(rootPath) {
		return nil, fmt.Errorf("rootPath is not absolute: %s", rootPath)
	}
	if len(baseNameRegex) == 0 {
		return nil, fmt.Errorf("baseNameRegex cannot be empty")
	}

	return basenameRegexGlob(rootPath, baseNameRegex)
}

// SearchByExtension Given a path, returns all sub-paths to files that have the specified extension 'ext'.
// Note that 'ext' must include a dot.
func SearchByExtension(rootPath string, ext string) ([]string, error) {
	if !filepath.IsAbs(rootPath) {
		return nil, fmt.Errorf("rootPath is not absolute: %s", rootPath)
	}
	if len(ext) == 0 {
		return nil, fmt.Errorf("extension cannot be empty")
	}
	if !strings.HasPrefix(ext, ".") {
		return nil, fmt.Errorf("extension must start with a dot (.): %s", ext)
	}

	return extensionGlob(rootPath, ext)
}

// SearchClosestParentContaining returns the closest parent (which may be 'startPath' itself) that contains
// a file or directory named 'contains'. Note that the name must match exactly.
// startPath may be relative or absolute.
// On success, it will return an absolute path to a directory containing something named 'contain'
func SearchClosestParentContaining(startPath string, contains string) (string, error) {
	absPath, err := filepath.Abs(startPath)
	if err != nil {
		return "", fmt.Errorf("failed to locate repo root for path: %s", startPath)
	}
	var currentLoc = absPath
	for {
		_, err = os.Stat(filepath.Join(currentLoc, contains))
		if err == nil {
			// We found something and may return immediately
			return currentLoc, nil
		}
		// At this point, if we have an error, it should tell us that the path did not exist
		if !os.IsNotExist(err) {
			// Anything else than "NotExist" is an issue we need to report.
			return "", err
		}
		if currentLoc == "/" {
			return "", fmt.Errorf("failed to locate anything named %s in %s or any of its parents", contains, startPath)
		}
		// Go up one level and start again...
		currentLoc = filepath.Dir(currentLoc)
	}
}

// filepath.Glob does not support things like '**/file'
func basenameGlob(dir string, baseName string) ([]string, error) {

	var files []string
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if filepath.Base(path) == baseName {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

// filepath.Glob does not support things like '**/file'
func basenameRegexGlob(dir string, baseNameRegex string) ([]string, error) {
	r, err := regexp.Compile(baseNameRegex)
	if err != nil {
		return nil, err
	}

	var files []string
	err = filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if r.MatchString(filepath.Base(path)) {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func extensionGlob(dir string, ext string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == ext {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}
