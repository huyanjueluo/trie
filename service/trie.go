package service

type Trie interface {
	Insert(words ...string)
	Walk()
	Search(segment string) []string
	Delete(word string)
	Filter(msg string) string
}
