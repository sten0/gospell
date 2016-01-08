package main

import (
	"errors"
//	"fmt"
)

// TODO: make this package protected?
type Trie struct {
	// TODO: distinguish between root and subsequent
	// nodes to avoid having to copy values?
	alphabet Alphabet
	// only handle lowercase english letters
	children   [28]*Trie // hardcoded. yup, it's big time hacky
	terminates bool
}

func (t *Trie) Contains(word []rune) bool {
	for _, l := range word {
		idx, ok := t.alphabet.Index(l)
		if !ok {
			return false
		}
		child := t.children[idx]
		if child == nil {
			return false
		}
		t = child
	}

	if !t.terminates {
		return false
	}

	return true
}

func (t *Trie) Alphabet() Alphabet {
	return t.alphabet
}

func NewTrie(words []string, alphabet Alphabet) (*Trie, error) {
	trie := &Trie{alphabet: alphabet}
	for _, w := range words {
		letters := []rune(w)
		currentTrie := trie
		for _, l := range letters {
			idx, ok := alphabet.Index(l)
			if !ok {
				return nil, errors.New("unicode char '%c' not found in alphabet.")
			}
//			fmt.Printf("currentTrie.children -> '%v'\n", currentTrie.children)
//			fmt.Printf("rune -> '%c', idx -> '%d'\n", l, idx)
			if currentTrie.children[idx] == nil {
				nTrie := &Trie{alphabet: alphabet}
				currentTrie.children[idx] = nTrie
				currentTrie = nTrie
			} else {
				currentTrie = currentTrie.children[idx]
			}
		}
		currentTrie.terminates = true
	}
	return trie, nil
}
