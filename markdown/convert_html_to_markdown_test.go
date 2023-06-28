package markdown

import (
	"io"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

// TestWalk tests the walk function.
func TestWalk(t *testing.T) {
	option := &Option{
		TrimSpace: true,
		customRulesMap: map[string]WalkFunc{
			"b": func(node *html.Node, w io.Writer, nest int, option *Option) {
				w.Write([]byte("*"))
				walk(node, w, nest, option)
				w.Write([]byte("*"))
			},
		},
	}

	node := &html.Node{
		Type: html.ElementNode,
		Data: "b",
		FirstChild: &html.Node{
			Type: html.TextNode,
			Data: "test",
		},
	}

	var b strings.Builder
	walk(node, &b, 0, option)

	if result := b.String(); result != "*test*" {
		t.Errorf("Expected '*test*', got '%s'", result)
	}
}

// Test that the method `isChildOf` works correctly
func TestIsChildOf(t *testing.T) {
	divNode := &html.Node{Data: "div"}
	spanNode := &html.Node{Data: "span", Parent: divNode}
	pNode := &html.Node{Data: "p", Parent: spanNode}

	if !isChildOf(pNode, "span") {
		t.Error("pNode should be a child of spanNode")
	}

	if isChildOf(pNode, "div") {
		t.Error("pNode should not be a direct child of divNode")
	}

	if !isChildOf(spanNode, "div") {
		t.Error("spanNode should be a child of divNode")
	}
}

// Test that the method `hasClass` works correctly
func TestHasClass(t *testing.T) {
	node := &html.Node{
		Attr: []html.Attribute{
			{Key: "class", Val: "class1 class2"},
		},
	}

	if !hasClass(node, "class1") {
		t.Error("Node should have class1")
	}

	if !hasClass(node, "class2") {
		t.Error("Node should have class2")
	}

	if hasClass(node, "class3") {
		t.Error("Node should not have class3")
	}
}

// Test that the method `attr` works correctly
func TestAttr(t *testing.T) {
	node := &html.Node{
		Attr: []html.Attribute{
			{Key: "class", Val: "myclass"},
			{Key: "id", Val: "myid"},
		},
	}

	if attr(node, "class") != "myclass" {
		t.Error("Node should have class as myclass")
	}

	if attr(node, "id") != "myid" {
		t.Error("Node should have id as myid")
	}

	if attr(node, "style") != "" {
		t.Error("Node should not have style attribute")
	}
}

// Test that the method `langFromClass` works correctly
func TestLangFromClass(t *testing.T) {
	codeNode := &html.Node{Data: "code"}
	preNode := &html.Node{Data: "pre", FirstChild: codeNode}
	preNode.Attr = []html.Attribute{{Key: "class", Val: "language-golang"}}

	if lang := langFromClass(preNode); lang != "golang" {
		t.Errorf("Expected golang, but got %s", lang)
	}

	preNode.Attr = []html.Attribute{{Key: "class", Val: "other-class language-python"}}

	if lang := langFromClass(preNode); lang != "python" {
		t.Errorf("Expected python, but got %s", lang)
	}

	preNode.Attr = []html.Attribute{{Key: "class", Val: "other-class"}}

	if lang := langFromClass(preNode); lang != "" {
		t.Errorf("Expected an empty string, but got %s", lang)
	}
}
