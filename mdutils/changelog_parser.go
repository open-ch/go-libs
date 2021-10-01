package mdutils

import (
	"fmt"
	"regexp"

	"github.com/russross/blackfriday/v2"
)

const reStr = `^(\d{1,2}\.\d{1,2}\.\d{4}).*`

var dateRe = regexp.MustCompile(reStr)

// Change contains the parsed content of a single change from a CHANGELOG file
type Change struct {
	Date         string
	HeaderNode   *blackfriday.Node
	ContentNodes []*blackfriday.Node
}

// GetChanges explores the passed ast, expecting it to contain a changelog
// with section titles that start with a date in the form DD.MM.YYYY
// it returns an ordered slice of changes
func GetChanges(ast *blackfriday.Node) ([]*Change, error) {

	// We are passed the entire document AST.
	directChildren := getDirectChildren(ast)
	onlyReleaseData := dropNonReleaseInfo(directChildren)
	return groupByRelease(onlyReleaseData)
}

// getAllChildren returns a slice containing all direct children of the passed node.
func getDirectChildren(ast *blackfriday.Node) []*blackfriday.Node {

	if ast.FirstChild == nil {
		// Empty document: nothing to return
		return nil
	}

	// Headers and paragraphs are first-level children of the document,
	// only lists have nested content, meaning we can just iterate over the
	// document's children

	current := ast.FirstChild
	children := []*blackfriday.Node{current}

	for current.Next != nil {
		children = append(children, current.Next)
		current = current.Next
	}

	return children
}

// drop everything that comes before the first header in DD.MM.YYYY format
func dropNonReleaseInfo(nodes []*blackfriday.Node) []*blackfriday.Node {
	for idx, node := range nodes {
		if node.Type == blackfriday.Heading && dateRe.Match(node.FirstChild.Literal) {
			return nodes[idx:]
		}
	}
	return []*blackfriday.Node{}
}

func groupByRelease(nodes []*blackfriday.Node) ([]*Change, error) {
	var releases = make([]*Change, 0)
	if len(nodes) == 0 {
		return releases, nil
	}

	var headerNode = nodes[0]
	var releaseContent = make([]*blackfriday.Node, 0)

	// Walk through all nodes
	for _, node := range nodes[1:] {
		if node.Type == blackfriday.Heading && dateRe.Match(node.FirstChild.Literal) {
			// we reached another release: add the previous one to the
			// slice that will be returned
			dateStr, err := getDateStringFromHeaderNode(headerNode)
			if err != nil {
				return nil, err
			}
			releases = append(releases, &Change{
				Date:         dateStr,
				HeaderNode:   nodes[0],
				ContentNodes: releaseContent,
			})
			// ... and initialize another one
			headerNode = node
			releaseContent = []*blackfriday.Node{}
		} else {
			releaseContent = append(releaseContent, node)
		}
	}
	dateStr, err := getDateStringFromHeaderNode(headerNode)
	if err != nil {
		return nil, err
	}
	// Add the last release in the nodes slice
	releases = append(releases, &Change{
		Date:         dateStr,
		HeaderNode:   headerNode,
		ContentNodes: releaseContent,
	})

	return releases, nil
}

func getDateStringFromHeaderNode(node *blackfriday.Node) (string, error) {
	if node.Type != blackfriday.Heading {
		return "", fmt.Errorf("expected a Heading node. Was: %v", node.Type)
	}
	parts := dateRe.FindSubmatch(node.FirstChild.Literal)
	return string(parts[1]), nil
}
