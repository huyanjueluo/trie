package schema

type Trie struct {
	Word  rune           // 节点的值
	IsEnd bool           // 是否叶子节点
	Child map[rune]*Trie //子节点
}
