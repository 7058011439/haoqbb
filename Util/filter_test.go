package Util

import "testing"

func TestFilter(t *testing.T) {
	trie := NewTrie()
	trie.Insert("习近平")
	trie.Insert("胡锦涛")

	fibTests := map[string]bool{
		"习近平是人": true,
		"是习近平人": true,
		"是人习近平": true,
		"习金平是人": false,
		"胡锦涛是人": true,
		"是胡锦涛人": true,
		"是胡锦涛平": true,
		"胡锦涛":   true,
		"习近平":   true,
		"习 近 平": false,
		"胡 锦 涛": false,
	}
	for word, ret := range fibTests {
		actual := trie.HasDirty(word)
		if actual != ret {
			t.Errorf("HasDirty(%v) = %v, expected = %v", word, actual, ret)
		}
	}
}
