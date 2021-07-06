package mdutils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetChanges(t *testing.T) {
	testFile := filepath.Join(getTestDir(), "test-changelog.md-ext")
	ast, _ := ParseFileToAst(testFile)
	changes, err := GetChanges(ast)
	assert.Nil(t, err, "Expected to successfully interpret changelog.")
	assert.NotNil(t, changes, "Expected to find some content in the changelog.")
	assert.Equal(t, 11, len(changes))

	// Check ordering
	assert.Equal(t, "11.05.2020", changes[0].Date)
	assert.Equal(t, "02.06.2021", changes[10].Date)
}
