package controller

import (
	"net/http"
	"practice/trie/service"
)

type Trie struct {
	Service service.Trie
}

func (t *Trie) Walk(w http.ResponseWriter, r *http.Request) {
	t.Service.Walk()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte{})
}
