package Util

/*
脏字过滤库
*/
type Trie struct {
	Root *TrieNode
}

type TrieNode struct {
	children map[rune]*TrieNode
	end      bool
	data     interface{}
}

func NewTrie() Trie {
	return Trie{
		Root: NewTrieNode(),
	}
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
	}
}

func (t *Trie) Insert(txt string) {
	if len(txt) < 1 {
		return
	}
	node := t.Root
	key := []rune(txt)
	for i := 0; i < len(key); i++ {
		if _, exists := node.children[key[i]]; !exists {
			node.children[key[i]] = NewTrieNode()
		}
		node = node.children[key[i]]
	}

	node.data = txt
	node.end = true
}

func (t *Trie) HasDirty(txt string) bool {
	if len(txt) < 1 {
		return false
	}
	node := t.Root
	key := []rune(txt)
	sLen := len(key)
	for i := 0; i < sLen; i++ {
		word := key[i]
		if _, exists := node.children[word]; exists {
			node = node.children[word]
			for j := i + 1; j < sLen; j++ {
				ret, exists1 := node.children[key[j]]
				if !exists1 {
					node = t.Root
					break
				} else {
					if ret.end == true {
						//fmt.Println(txt, ret.data.(string))
						return true
					} else {
						node = ret
					}
				}
			}
		}
	}
	return false
}
