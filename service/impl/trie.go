package impl

import (
	"fmt"
	"strings"
)

type Trie struct {
	Word  rune           // 节点的值
	IsEnd bool           // 是否叶子节点
	Child map[rune]*Trie //子节点
}

func (t *Trie) Insert(words ...string) {
	for _, word := range words {
		word = strings.TrimSpace(word)
		ptr := t
		for _, u := range word {
			_, ok := ptr.Child[u]
			if !ok {
				node := make(map[rune]*Trie)
				ptr.Child[u] = &Trie{
					Word:  u,
					IsEnd: false,
					Child: node,
				}
			}
			ptr = ptr.Child[u]
		}
		ptr.IsEnd = true
	}

}

func (t *Trie) Walk() {
	var walk func(string, *Trie)
	walk = func(pfx string, node *Trie) {
		if node == nil {
			return
		}
		if node.Word != 0 {
			pfx += string(node.Word)
		}
		if node.IsEnd {
			fmt.Println(pfx)
		}
		for _, v := range node.Child {
			walk(pfx, v)
		}
	}
	walk("", t)
}

func (t *Trie) Search(segment string) []string {
	segment = strings.TrimSpace(segment)
	segmentRune := []rune(segment)
	ptr := t
	var matched []string
	item := ""
	index := 0
	for i := 0; i < len(segmentRune); i++ {
		c, ok := ptr.Child[segmentRune[i]]
		if !ok {
			i = index
			index++
			item = ""
			ptr = t
			continue
		}
		item += string(c.Word)

		if c.IsEnd {
			matched = append(matched, item)
			if len(c.Child) == 0 {
				i = index
				index++
				ptr = t
				item = ""
				continue
			}
		}
		ptr = c
	}
	return matched
}

func (t *Trie) Delete(word string) {
	word = strings.TrimSpace(word)
	var branch []*Trie
	ptr := t
	for _, u := range word {
		branch = append(branch, ptr)
		c, ok := ptr.Child[u]
		if !ok {
			return
		}
		ptr = c
	}
	// 只命中字典中部分词
	if !ptr.IsEnd {
		return
	}
	// 如bitch和bitches
	// 删除bitch时，只需要将bitch最后一个节点的IsEnd改为false即可
	if len(ptr.Child) != 0 {
		ptr.IsEnd = false
		return
	}
	for len(branch) > 0 {
		p := branch[len(branch)-1]
		branch = branch[:len(branch)-1]

		delete(p.Child, ptr.Word)
		// IsEnd == true 如bitch和bitches，删除bitches时，只需要删除后面的"es"即可
		// len(Child) != 0 整个敏感词全删除
		if p.IsEnd || len(p.Child) != 0 {
			break
		}
		ptr = p
	}
}

func (t *Trie) Filter(msg string) string {
	matched := t.Search(msg)
	if len(matched) != 0 {
		var oldNew []string
		for _, v := range matched {
			oldNew = append(oldNew, v, "***")
		}
		replacer := strings.NewReplacer(oldNew...)
		return replacer.Replace(msg)
	}
	return msg
}
