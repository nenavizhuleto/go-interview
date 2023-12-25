package utils

import (
	"fmt"
	"helpdesk/internals/core/tree"
	"strings"
)

func ToUpper(v string) string {
	return strings.ToUpper(v)
}

func ToLower(v string) string {
	return strings.ToLower(v)
}

func ToSlug(v string) string {
	return ToLower(strings.ReplaceAll(v, " ", "-"))
}

func PrintTree(node *tree.Node, indent int) string {
	if node == nil {
		return "--*\n"
	}

	tree := fmt.Sprintf("%sID: %s VALUE: %v\n", strings.Repeat("\t", indent), node.ID, node.Value)

	if len(node.Children) > 0 {
		for _, child := range node.Children {
			tree += PrintTree(child, indent+1)
		}
	}

	return tree

}
