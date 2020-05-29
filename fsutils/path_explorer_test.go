package fsutils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestDir() string {
	workDir, _ := os.Getwd()
	return filepath.Join(workDir, "test-data")
}

func TestSearchByFileName(t *testing.T) {
	testDir := getTestDir()

	_, notAbsErr := SearchByFileName("relative/path", "README")
	assert.NotNil(t, notAbsErr, "Expected an error if provided with a relative path")

	_, emptyNameErr := SearchByFileName(testDir, "")
	assert.NotNil(t, emptyNameErr, "Expected an error if provided with an empty filename")

	singleNoExtension, _ := SearchByFileName(testDir, "README")
	assert.Equal(t, 1, len(singleNoExtension), "Expected a single match, at the root of the test directory")
	assert.True(t, strings.HasSuffix(singleNoExtension[0], "README"))

	withExtension, _ := SearchByFileName(testDir, "README.md-ext")
	assert.Equal(t, 2, len(withExtension), "Expected two matches")
}

func TestSearchByExtension(t *testing.T) {
	testDir := getTestDir()

	_, notAbsErr := SearchByExtension("relative/path", ".md-ext")
	assert.NotNil(t, notAbsErr, "Expected an error if provided with a relative path")

	_, emptyNameErr := SearchByExtension(testDir, "")
	assert.NotNil(t, emptyNameErr, "Expected an error if provided with an empty extension")

	_, noDot := SearchByExtension(testDir, "md-ext")
	assert.NotNil(t, noDot, "Expected an error if provided with an extension not starting with a dot.")

	withExtension, _ := SearchByExtension(testDir, ".md-ext")
	assert.Equal(t, 2, len(withExtension), "Expected a single match, at the root of the test directory")
}

func TestSearchClosestParentContaining(t *testing.T) {
	testDir := filepath.Join(getTestDir(), "sub-dir")

	sameLevel, err := SearchClosestParentContaining(testDir, "README.md-ext")
	assert.Nil(t, err, "Should not fail if a valid parent exists")
	assert.True(t, strings.HasSuffix(sameLevel, "/test-data/sub-dir"))

	levelAbove, err := SearchClosestParentContaining(testDir, "README")
	assert.Nil(t, err, "Should not fail if a valid parent exists")
	assert.True(t, strings.HasSuffix(levelAbove, "/test-data"))

	empty, err := SearchClosestParentContaining(testDir, "IHopeFullyDontExistAnyWhereInTheParentHierarchy")
	assert.NotNil(t, err, "Should not fail if a valid parent exists")
	assert.Empty(t, empty, "On errors an empty string should be returned")
}
