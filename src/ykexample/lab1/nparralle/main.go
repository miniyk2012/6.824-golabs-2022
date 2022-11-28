package main

import (
	"fmt"
	"log"
	"plugin"

	"x6.824/mr"
)

/*
	go build -race -buildmode=plugin ../../../mrapps/mtiming.go
	go run -race main.go
*/

func main() {
	mapf, reducef := loadPlugin("mtiming.so")
	_ = reducef
	content := `hello world
	i love you
	what do you looking for
	`
	filename := "example.txt"
	kv := mapf(filename, content)
	fmt.Println(kv)

}

//
// load the application Map and Reduce functions
// from a plugin file, e.g. ../mrapps/wc.so
//
func loadPlugin(filename string) (func(string, string) []mr.KeyValue, func(string, []string) string) {
	p, err := plugin.Open(filename)
	if err != nil {
		log.Fatalf("cannot load plugin %v, because %s", filename, err)
	}
	xmapf, err := p.Lookup("Map")
	if err != nil {
		log.Fatalf("cannot find Map in %v", filename)
	}
	mapf := xmapf.(func(string, string) []mr.KeyValue)
	xreducef, err := p.Lookup("Reduce")
	if err != nil {
		log.Fatalf("cannot find Reduce in %v", filename)
	}
	reducef := xreducef.(func(string, []string) string)

	return mapf, reducef
}
