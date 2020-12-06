package main

import (
	"fmt"
	"net/http"
	"practice/trie/controller"
	"practice/trie/service/impl"

	"go.uber.org/dig"
)

type ServerConfig struct {
	Host string
	Port string
}

type Router struct {
	Routes *map[string]http.HandlerFunc
}

type Server struct {
	ServerConfig *ServerConfig
	Router       *Router
}

func (svc *Server) Run() {
	for path := range *svc.Router.Routes {
		http.HandleFunc(path, (*svc.Router.Routes)[path])
	}
	_ = http.ListenAndServe(svc.ServerConfig.Host+":"+svc.ServerConfig.Port, nil)
}
func NewTrieService() *impl.Trie {
	return &impl.Trie{
		Word:  0,
		IsEnd: false,
		Child: make(map[rune]*impl.Trie),
	}
}

func NewTrieController(trieService *impl.Trie) *controller.Trie {
	return &controller.Trie{Service: trieService}
}

func NewRouter(trieController *controller.Trie) *Router {
	router := make(map[string]http.HandlerFunc)
	router["/walk"] = trieController.Walk
	return &Router{Routes: &router}
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Host: "127.0.0.1",
		Port: "9721",
	}
}

func NewServer(svcCfg *ServerConfig, router *Router) *Server {
	return &Server{
		ServerConfig: svcCfg,
		Router:       router,
	}
}

func BuildContainer() *dig.Container {
	container := dig.New()
	_ = container.Provide(NewTrieService)
	_ = container.Provide(NewTrieController)
	_ = container.Provide(NewRouter)
	_ = container.Provide(NewServerConfig)
	_ = container.Provide(NewServer)
	return container
}

func main() {
	container := BuildContainer()
	err := container.Invoke(func(trie *impl.Trie) {
		trie.Insert("你大爷", "大姨妈", "姨妈jin", "jin子", "bitch", "bitches", "妈了个吧", "狗日的", "去你吗的")
		trie.Walk()
		matched := trie.Search("你大爷的")
		fmt.Println(fmt.Sprintf("search: %s, expected: %v, actual: %v", "你大爷的", []string{"你大爷"}, matched))
		matched = trie.Search("英文单词bitches意思是母狗")
		fmt.Println(fmt.Sprintf("search: %s, expected: %v, actual: %v", "英文单词bitches意思是母狗", []string{"itches"}, matched))
		matched = trie.Search("狗日的大姨妈啊")
		fmt.Println(fmt.Sprintf("search: %s, expected: %v, actual: %v", "狗日的大姨妈啊", []string{"狗日的", "大姨妈"}, matched))
		matched = trie.Search("我去你大爷的")
		fmt.Println(fmt.Sprintf("search: %s, expected: %v, actual: %v", "我去你大爷的", []string{"你大爷"}, matched))
		matched = trie.Search("大姨妈jin子")
		fmt.Println(fmt.Sprintf("search: %s, expected: %v, actual: %v", "大姨妈jin子", []string{"大姨妈", "姨妈jin", "jin子"}, matched))
		matched = trie.Search("我很正常")
		fmt.Println(fmt.Sprintf("search: %s, expected: %v, actual: %v", "我很正常", []string{}, matched))
		trie.Delete("你大爷")
		matched = trie.Search("你大爷")
		fmt.Println(fmt.Sprintf("search: %s, expected: %v, actual: %v", "你大爷", []string{}, matched))
		result := trie.Filter("大姨妈jin子啊")
		fmt.Println(fmt.Sprintf("filter: %s, expected: %v, actual: %v", "大姨妈jin子啊", "******啊", result))
	})
	if err != nil {
		panic(err)
	}
	err = container.Invoke(func(server *Server) {
		server.Run()
	})
	if err != nil {
		panic(err)
	}
}
