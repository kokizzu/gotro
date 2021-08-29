package main

import (
	"os"
	"strings"
)

// purpose: convert CLI/STDIN input into STDOUT/STDERR output thru domain functions

func cliParseInput(file *os.File, in interface{}) error {
	// TODO: parse input from stdin, depends on input format: json, yaml, toml, ini/env-like
	return nil
}

type node struct {
	pattern string
	nodes   map[string]*node
}

// build a segment trie
func (n *node) ParseSegment(segments []string, pattern string) {
	n.pattern = pattern
	if len(segments) == 0 {
		return
	}
	segment := segments[0]
	if segment == `` {
		n.ParseSegment(segments[1:], pattern)
	} else {
		if n.nodes == nil {
			n.nodes = map[string]*node{}
		}
		if _, ok := n.nodes[segment]; !ok {
			n.nodes[segment] = &node{}
		}
		n.nodes[segment].ParseSegment(segments[1:], pattern)
	}
}

// match based on url
func (n *node) GetPattern(url string) string {
	segments := splitUrlToSegments(url)
	visit := n
	for _, segment := range segments {
		if segment == `` {
			continue
		}
		if child, ok := visit.nodes[segment]; ok {
			visit = child
		} else {
			break
		}
	}
	return visit.pattern
}

func splitUrlToSegments(url string) []string {
	qs := strings.Split(url, `?`) // remove ?queryString
	return strings.Split(qs[0], `/`)
}

// build segment trie and match against url, return nearest pattern
func cliUrlPattern(url string, patterns map[string]map[string]int) string {
	root := node{}

	for pattern := range patterns {
		root.ParseSegment(splitUrlToSegments(pattern), pattern)
	}
	root.ParseSegment([]string{``}, `/`) // to clear root pattern
	// TODO: continue this
	return root.GetPattern(url)
}

// return matching segment
func cliSegmentFromIdx(url string, idx int) string {
	segments := strings.Split(url, `/`)
	if idx > 0 && idx < len(segments) {
		return segments[idx]
	}
	return ``
}
